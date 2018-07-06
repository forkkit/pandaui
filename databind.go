package slugui

import (
	"reflect"
	"sync"

	"github.com/influx6/faux/reflection"
)

// constant errors.,
const (
	ErrNotExists                           = Error("compute type does not exists")
	ErrNameInUse                           = Error("computed name already used")
	ErrUncomputed                          = Error("data uncomputed yet, use appropriate action")
	ErrMustProvideMergeKeys                = Error("Merge/Mutations require atleast one existing source name")
	ErrMustExpectLastValueAsArgument       = Error("Function must expect previous value as argument")
	ErrMustReturnNewStateValue             = Error("Function must return new state value as return atleast")
	ErrMustReturn                          = Error("Function must return value")
	ErrMustAcceptAction                    = Error("Function must accept *Action type")
	ErrMustAcceptLastValueAndMergeValues   = Error("Function must accept last value and more than 1 merge values")
	ErrMutationIsMergingMultipleValues     = Error("Mutations is merging of multiple values, arguments must be > 1")
	ErrInTwoReturnsLastMustBeError         = Error("Function returning two values must return an error as second")
	ErrReturnValueMustBeSettableAsArgument = Error("Function must return a new value of same type as first argument")
)

var (
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
	actionType = reflect.TypeOf((*Action)(nil)).Elem()

	baseReturnPolicy = func(arguments []reflect.Type, returns []reflect.Type) error {
		switch len(returns) {
		case 1:
			// if it's a single value then it can not be
			// an error.
			if !reflection.IsSettableType(errorType, returns[0]) {
				return nil
			}
		case 2:
			// if we have a function returning two values, then
			// the last one must be an error type.
			if reflection.IsSettableType(errorType, returns[1]) {
				return nil
			}
			return ErrInTwoReturnsLastMustBeError
		}
		return ErrMustReturn
	}

	sourceFunctionRules = []reflection.AreaValidation{
		baseReturnPolicy,
		func(arguments []reflect.Type, returns []reflect.Type) error {
			if len(returns) == 0 {
				return ErrMustReturnNewStateValue
			}
			if len(arguments) == 0 && len(returns) != 0 {
				if reflection.IsSettableType(returns[0], errorType) {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if len(arguments) == 1 && reflection.IsSettableType(arguments[0], actionType) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if len(arguments) == 1 && !reflection.IsSettableType(arguments[0], actionType) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}

				if !reflection.IsSettableType(returns[0], arguments[0]) {
					return ErrReturnValueMustBeSettableAsArgument
				}
				return nil
			}
			if len(arguments) == 2 && reflection.IsSettableType(arguments[1], actionType) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}

				if !reflection.IsSettableType(returns[0], arguments[0]) {
					return ErrReturnValueMustBeSettableAsArgument
				}
				return nil
			}
			return nil
		},
	}

	mutationFunctionRules = []reflection.AreaValidation{
		baseReturnPolicy,
		func(arguments []reflect.Type, returns []reflect.Type) error {
			if len(arguments) == 0 {
				return ErrMustExpectLastValueAsArgument
			}
			if len(returns) == 0 {
				return ErrMustReturnNewStateValue
			}
			if len(arguments) == 1 {
				return ErrMutationIsMergingMultipleValues
			}
			if len(returns) >= 1 {
				if reflection.IsSettableType(returns[0], errorType) {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			return nil
		},
	}

	ctrlMutationFunctionRules = []reflection.AreaValidation{
		baseReturnPolicy,
		func(arguments []reflect.Type, returns []reflect.Type) error {
			if len(arguments) == 0 {
				return ErrMustExpectLastValueAsArgument
			}
			if len(returns) == 0 {
				return ErrMustReturnNewStateValue
			}
			if len(arguments) == 1 {
				return ErrMutationIsMergingMultipleValues
			}
			if len(returns) >= 1 {
				if reflection.IsSettableType(returns[0], errorType) {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if len(arguments) > 1 {
				if !reflection.IsSettableType(arguments[0], actionType) {
					return ErrMustAcceptAction
				}

				if len(arguments) == 2 {
					return ErrMustAcceptLastValueAndMergeValues
				}
			}
			return nil
		},
	}
)

// Action defines an intent about an action to be performed
// by a Source against it's previous data to produce new data.
type Action struct {
	Type  string
	Value interface{}
}

type computation struct {
	linked       []string
	Pure         func() (interface{}, error)
	InPure       func(interface{}) (interface{}, error)
	Action       func(interface{}, *Action) (interface{}, error)
	Mutation     func(interface{}, ...interface{}) (interface{}, error)
	CtrlMutation func(*Action, interface{}, ...interface{}) (interface{}, error)
}

// DataBind defines the interface which embodies the concept of the
// data system which powers the data pipeline for all Components
// Surfaces, and attributes in slugui.
type DataBind struct {
	rw           sync.RWMutex
	lastValues   map[string]interface{}
	computations map[string]computation
}

func Redux() *DataBind {
	return &DataBind{
		lastValues:   make(map[string]interface{}),
		computations: make(map[string]computation),
	}
}

// All returns all last computed values and associated keys
// in a map. Values returned are values of Sources, Mutations
// and ActionMutations that have being called previously be user.
// All others without previously value computed, will be absent.
func (db *DataBind) All() map[string]interface{} {
	db.rw.RLock()
	defer db.rw.RUnlock()

	items := make(map[string]interface{})
	if db.lastValues == nil {
		return items
	}

	for key, value := range db.lastValues {
		items[key] = value
	}
	return items
}

// Get returns the computed value for a giving computation name.
//
// It follows the following rules:
//
// 1. If a computation is "pure", where it expects no action and no previous value,
//    then we are required to call it every single time, it's value
//    is requested.
//
// 2. If a computation is "in-pure", where it expects no action and but a previous value,
//    then we are required to call it every single time, it's value
//    is requested.
//
// 3. If a computation is action based, then we must return last
//    known value if previously computed or an ErrUncomputed error.
//
// 4. If a computation is merged, then we must return last
//    known value if previously computed. If no previous value then
//    we must retrieve value of all concerned interests and pass to
//    mutation for resolving, returning error if any or returned value.
//
// 5. If a computation is merged by action, then we must return last
//    known value if previously computed or an ErrUncomputed error.
//
func (db *DataBind) Get(targetName string) (interface{}, error) {
	db.rw.RLock()
	computed, ok := db.computations[targetName]
	if !ok {
		db.rw.RUnlock()
		return nil, ErrNotExists
	}

	last, computedBefore := db.lastValues[targetName]
	db.rw.RUnlock()

	// if it's a pure computation, then execute computation
	// for new value.
	if computed.Pure != nil {
		computed, err := computed.Pure()
		if err != nil {
			return nil, err
		}

		// Add computed value into last known values.
		db.rw.Lock()
		db.lastValues[targetName] = computed
		db.rw.Unlock()

		return computed, nil
	}

	if computed.InPure != nil {
		computed, err := computed.InPure(last)
		if err != nil {
			return nil, err
		}

		// Add computed value into last known values.
		db.rw.Lock()
		db.lastValues[targetName] = computed
		db.rw.Unlock()

		return computed, nil
	}

	if !computedBefore {
		return nil, ErrUncomputed
	}

	return last, nil
}

// Static adds giving value into DataBind has a static function call
// that always returns the same value. It calls DataBind.Source underneath.
func (db *DataBind) Static(value interface{}, name string) error {
	return db.Action(func() interface{} {
		return value
	}, name)
}

// Action adds giving function which based on produce provided value tagged
// with provided name.
func (db *DataBind) Action(fn interface{}, name string) error {
	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	if err := reflection.ValidateFuncArea(fn, sourceFunctionRules...); err != nil {
		return err
	}

	var computed computation
	db.computations[name] = computed
	return nil
}

// Merged adds giving function associated with name and provided merging as factors
// for key update.
func (db *DataBind) Merged(fn interface{}, name string, merging ...string) error {
	if len(merging) == 0 {
		return ErrMustProvideMergeKeys
	}

	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	if err := reflection.ValidateFuncArea(fn, mutationFunctionRules...); err != nil {
		return err
	}

	var computed computation
	computed.linked = merging

	db.computations[name] = computed
	return nil
}

func (db *DataBind) MergedBy(fn interface{}, name string, merging ...string) error {
	if len(merging) == 0 {
		return ErrMustProvideMergeKeys
	}

	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	if err := reflection.ValidateFuncArea(fn, ctrlMutationFunctionRules...); err != nil {
		return err
	}

	var computed computation
	computed.linked = merging

	db.computations[name] = computed
	return nil
}

func (db *DataBind) Compute(actionName string, actionValue interface{}) error {

	return nil
}

func (db *DataBind) ComputeFor(actionName string, actionValue string, specificSources ...string) error {
	return nil
}
