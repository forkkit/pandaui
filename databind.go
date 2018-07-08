package slugui

import (
	"reflect"
	"sync"

	"sync/atomic"

	"github.com/influx6/faux/reflection"
)

// constant errors.,
const (
	ErrNotExists                           = Error("compute type does not exists")
	ErrNameInUse                           = Error("computed name already used")
	ErrTypeIsNotError                      = Error("computed returned second value is not error type")
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

// State defines a map type which is used to store current
// state of existing values as key-value pairs.
type State map[string]interface{}

// Get returns the value of the giving key.
func (s State) Get(name string) interface{} {
	return s[name]
}

// Has returns true/false if giving key name exists.
func (s State) Has(name string) bool {
	_, ok := s[name]
	return ok
}

// Action defines an intent about an action to be performed
// by a Source against it's previous data to produce new data.
type Action struct {
	Type  string
	Value interface{}
}

// Feed defines a giving type of stream which captures it's
// last value and can be activated to update said value through
// a user defined computation optionally based on a Action.
type Feed interface {
	Get() interface{}
	Compute(*Action) error
}

type PureFeed struct {
	rw   sync.RWMutex
	last interface{}
}

// computation represent a type which is to be evaluated
// based on specific function type it services. It
// functions must be pure and consistently return new values
// which take the place of existing values.
// See https://gitlab.com/slugui/slugui/tree/master/specs/databind.spec.md.
type computation struct {
	changes      int64
	args         []reflect.Type
	returns      []reflect.Type
	updates      map[string]bool
	graph        map[string]struct{}
	Pure         func() (interface{}, error)
	InPure       func(interface{}) (interface{}, error)
	Action       func(interface{}, *Action) (interface{}, error)
	Materialize  func(interface{}, *Action) (interface{}, error)
	Mutation     func(interface{}, State) (interface{}, error)
	CtrlMutation func(*Action, interface{}, State) (interface{}, error)
}

// HasChanges returns true/false if the computation should
// react has updates it is interested in had occurred.
// This allows us re-computed a computation as needed and
// not create complicated on-demand computation but instead
// recompute on call.
func (c computation) HasChanges() bool {
	return atomic.LoadInt64(&c.changes) > 0
}

// Watches returns true/false if giving name is part of the
// existing types which is being watched by giving computation
// in case of mutation or controlled mutations.
func (c computation) Watches(name string) bool {
	_, ok := c.graph[name]
	return ok
}

// DataFeed defines the interface which embodies the concept of the
// data system which powers the data pipeline for all Components
// Surfaces, and attributes in slugui.
// See https://gitlab.com/slugui/slugui/tree/master/specs/databind.spec.md.
type DataFeed struct {
	rw           sync.RWMutex
	lastValues   map[string]interface{}
	computations map[string]computation
}

// Feeds returns a new instance of DataFeed which implements a redux-like
// data-management layer for the slug-ui project.
// See https://gitlab.com/slugui/slugui/tree/master/specs/databind.spec.md.
func Feeds() *DataFeed {
	return &DataFeed{
		lastValues:   make(map[string]interface{}),
		computations: make(map[string]computation),
	}
}

// All returns all values of computable keys either
// based on last know value for where changes are yet
// to occur after last call, latest values returned in
// case of Pure functions who return new values on every call,
// or pure mutations that already have values calculated or
// due to existing change will be recalculated.
// The only exceptions to this is that if computations that
// use Actions has means of merging have not being called,
// they will be skipped from list.
func (db *DataFeed) All() map[string]interface{} {
	db.rw.RLock()
	defer db.rw.RUnlock()

	items := make(map[string]interface{})
	for name, comm := range db.computations {
		if !comm.HasChanges() {
			items[name] = db.lastValues[name]
			continue
		}

		// Unlock mutex has Get will hold lock
		// effectively.
		db.rw.RUnlock()
		value, err := db.Get(name)
		if err != nil {
			// Ensure to read lock again to ensure lock safety.
			db.rw.RLock()
			continue
		}

		// Add new value into items for key 'name'.
		items[name] = value

		// Ensure to read lock again to ensure lock safety.
		db.rw.RLock()
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
func (db *DataFeed) Get(targetName string) (interface{}, error) {
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

		db.notifyUpdate(targetName)

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

		db.notifyUpdate(targetName)

		return computed, nil
	}

	if !computedBefore {
		return nil, ErrUncomputed
	}

	return last, nil
}

// StaticFeed adds giving value into DataFeed has a static function call
// that always returns the same value. It calls DataFeed.Source underneath.
func (db *DataFeed) StaticFeed(value interface{}, name string) error {
	return db.Feed(func() interface{} {
		return value
	}, name)
}

// Feed defines and adds a new feed source by it's provided key name which
// will represent the computed values from feed source. No two feed sources
// can have same key name. To learn more about DataFeed.Feeds,
// See https://gitlab.com/slugui/slugui/tree/master/specs/databind.spec.md.
func (db *DataFeed) Feed(fn interface{}, name string, actions ...string) error {
	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	args, err := reflection.GetFuncArgumentsType(fn)
	if err != nil {
		return err
	}

	rets, err := reflection.GetFuncReturnsType(fn)
	if err != nil {
		return err
	}

	if err := runRulesAgainst(sourceFunctionRules, args, rets); err != nil {
		return err
	}

	var computed computation
	computed.args = args
	computed.returns = rets

	if len(args) == 0 {
		if len(rets) == 1 {
			computed.Pure = makePureWithSingleReturn(name, fn)
		}

		if len(rets) == 2 {
			computed.Pure = makePureWithNormalReturn(name, fn)
		}
	}

	if len(args) == 1 {
		if len(rets) == 1 {
			computed.InPure = makeInpureWithSingleReturn(name, fn)
		}

		if len(rets) == 2 {
			computed.InPure = makeInpureWithNormalReturn(name, fn)
		}
	}

	if len(args) == 2 {
		if len(rets) == 1 {
			computed.Action = makeActionWithSingleReturn(name, fn)
		}

		if len(rets) == 2 {
			computed.Action = makeActionWithNormalReturn(name, fn)
		}
	}

	db.computations[name] = computed
	return nil
}

// DerivedFeed are feeds which are the result of a merge computation of other
// pure feeds or other DerivedFeed which then produces a unique result to be
// set to the value of the key name.
// See https://gitlab.com/slugui/slugui/tree/master/specs/databind.spec.md.
func (db *DataFeed) DerivedFeed(fn interface{}, name string, merging ...string) error {
	if len(merging) == 0 {
		return ErrMustProvideMergeKeys
	}

	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	args, err := reflection.GetFuncArgumentsType(fn)
	if err != nil {
		return err
	}

	rets, err := reflection.GetFuncReturnsType(fn)
	if err != nil {
		return err
	}

	if err := runRulesAgainst(mutationFunctionRules, args, rets); err != nil {
		return err
	}

	var computed computation
	computed.args = args
	computed.returns = rets
	computed.updates = make(map[string]bool, len(merging))

	graph := make(map[string]struct{}, len(merging))
	for _, depends := range merging {
		graph[depends] = struct{}{}
	}
	computed.graph = graph

	if len(rets) == 1 {
		computed.Mutation = makeMutationWithSingleReturn(name, fn)
		db.computations[name] = computed
		return nil
	}

	computed.Mutation = makeMutationWithNormalReturn(name, fn)
	db.computations[name] = computed
	return nil
}

func (db *DataFeed) DerivedActionFeed(fn interface{}, name string, merging []string, actions []string) error {
	if len(merging) == 0 {
		return ErrMustProvideMergeKeys
	}

	db.rw.Lock()
	defer db.rw.Unlock()

	if _, ok := db.computations[name]; ok {
		return ErrNameInUse
	}

	args, err := reflection.GetFuncArgumentsType(fn)
	if err != nil {
		return err
	}

	rets, err := reflection.GetFuncReturnsType(fn)
	if err != nil {
		return err
	}

	if err := runRulesAgainst(ctrlMutationFunctionRules, args, rets); err != nil {
		return err
	}

	var computed computation
	computed.args = args
	computed.returns = rets
	computed.updates = make(map[string]bool, len(merging))

	graph := make(map[string]struct{}, len(merging))
	for _, depends := range merging {
		graph[depends] = struct{}{}
	}
	computed.graph = graph

	if len(rets) == 1 {
		computed.CtrlMutation = makeCtrlMutationWithSingleReturn(name, fn)
		db.computations[name] = computed
		return nil
	}

	computed.CtrlMutation = makeCtrlMutationWithNormalReturn(name, fn)
	db.computations[name] = computed
	return nil
}

func (db *DataFeed) Compute(actionName string, actionValue interface{}) error {

	return nil
}

func (db *DataFeed) notifyUpdate(actionName string) {
	db.rw.RLock()
	defer db.rw.RUnlock()
	for name, comm := range db.computations {
		if !comm.Watches(actionName) {
			continue
		}

		atomic.AddInt64(&comm.changes, 1)
		comm.updates[actionName] = true
		db.notifyOthers(name)
	}
}

func (db *DataFeed) notifyOthers(actionName string) {
	for name, comm := range db.computations {
		if !comm.Watches(actionName) {
			continue
		}

		if comm.updates[actionName] {
			continue
		}

		comm.updates[actionName] = true
		atomic.AddInt64(&comm.changes, 1)
		db.notifyOthers(name)
	}
}

func runRulesAgainst(rules []reflection.AreaValidation, args []reflect.Type, returns []reflect.Type) error {
	for _, rule := range rules {
		if err := rule(args, returns); err != nil {
			return err
		}
	}
	return nil
}

func makeMutationWithSingleReturn(name string, fn interface{}) func(interface{}, State) (interface{}, error) {
	return func(value interface{}, others State) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value, others)
		if err != nil {
			return nil, err
		}

		return res[0], nil
	}
}

func makeMutationWithNormalReturn(name string, fn interface{}) func(interface{}, State) (interface{}, error) {
	return func(value interface{}, others State) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value, others)
		if err != nil {
			return nil, err
		}

		if res[1] == nil {
			return res[0], nil
		}

		callErr, ok := res[1].(error)
		if !ok {
			return res[0], ErrTypeIsNotError
		}

		return res[0], callErr
	}
}

func makeCtrlMutationWithSingleReturn(name string, fn interface{}) func(*Action, interface{}, State) (interface{}, error) {
	return func(req *Action, value interface{}, others State) (interface{}, error) {
		res, err := reflection.CallFunc(fn, req, value, others)
		if err != nil {
			return nil, err
		}

		return res[0], nil
	}
}

func makeCtrlMutationWithNormalReturn(name string, fn interface{}) func(*Action, interface{}, State) (interface{}, error) {
	return func(req *Action, value interface{}, others State) (interface{}, error) {
		res, err := reflection.CallFunc(fn, req, value, others)
		if err != nil {
			return nil, err
		}

		if res[1] == nil {
			return res[0], nil
		}

		callErr, ok := res[1].(error)
		if !ok {
			return res[0], ErrTypeIsNotError
		}

		return res[0], callErr
	}
}

func makeActionWithSingleReturn(name string, fn interface{}) func(interface{}, *Action) (interface{}, error) {
	return func(value interface{}, req *Action) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value, req)
		if err != nil {
			return nil, err
		}

		return res[0], nil
	}
}

func makeActionWithNormalReturn(name string, fn interface{}) func(interface{}, *Action) (interface{}, error) {
	return func(value interface{}, req *Action) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value, req)
		if err != nil {
			return nil, err
		}

		if res[1] == nil {
			return res[0], nil
		}

		callErr, ok := res[1].(error)
		if !ok {
			return res[0], ErrTypeIsNotError
		}

		return res[0], callErr
	}
}

func makeInpureWithSingleReturn(name string, fn interface{}) func(interface{}) (interface{}, error) {
	return func(value interface{}) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value)
		if err != nil {
			return nil, err
		}

		return res[0], nil
	}
}

func makeInpureWithNormalReturn(name string, fn interface{}) func(interface{}) (interface{}, error) {
	return func(value interface{}) (interface{}, error) {
		res, err := reflection.CallFunc(fn, value)
		if err != nil {
			return nil, err
		}

		if res[1] == nil {
			return res[0], nil
		}

		callErr, ok := res[1].(error)
		if !ok {
			return res[0], ErrTypeIsNotError
		}

		return res[0], callErr
	}
}

func makePureWithSingleReturn(name string, fn interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		res, err := reflection.CallFunc(fn)
		if err != nil {
			return nil, err
		}

		return res[0], nil
	}
}

func makePureWithNormalReturn(name string, fn interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		res, err := reflection.CallFunc(fn)
		if err != nil {
			return nil, err
		}

		if res[1] == nil {
			return res[0], nil
		}

		callErr, ok := res[1].(error)
		if !ok {
			return res[0], ErrTypeIsNotError
		}

		return res[0], callErr
	}
}
