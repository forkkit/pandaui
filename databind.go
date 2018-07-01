package slugui

import (
	"reflect"

	"github.com/influx6/faux/reflection"
)

// constant errors.,
const (
	ErrMustExpectLastValueAsArgument = Error("Function must expect previous value as argument")
	ErrMustReturnNewStateValue       = Error("Function must return new state value as return atleast")
	ErrMustBeActionType              = Error("Function taking second argument must have type as slugui.Action")
	ErrMustBeMoreThanOne             = Error("Function must accept more than one argument")
	ErrMustReturn                    = Error("Function must return value")
	ErrInTwoReturnsLastMustBeError   = Error("Function returning two values, " +
		"must return an error as second")
	ErrTooManyArguments                    = Error("Function must atmost expect 2 arguments")
	ErrReturnValueMustBeSettableAsArgument = Error(
		"Function must return a new value of same type as first argument or of a type settable as first argument")
)

var (
	errorType  = reflect.TypeOf((*error)(nil)).Elem()
	actionType = reflect.TypeOf((*Action)(nil)).Elem()

	areaFunctionRules = []reflection.AreaValidation{
		func(arguments []reflect.Type, returns []reflect.Type) error {
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
		},
		func(arguments []reflect.Type, returns []reflect.Type) error {
			if len(arguments) == 0 && len(returns) == 0 {
				return ErrMustReturnNewStateValue
			}
			if len(arguments) == 1 && reflection.IsSettableType(arguments[0], actionType) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if len(arguments) == 0 && reflection.IsSettableType(returns[0]) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if len(arguments) == 1 && !reflection.IsSettableType(arguments[0], actionType) {
				if len(returns) == 0 {
					return ErrMustReturnNewStateValue
				}
				return nil
			}
			if !reflection.IsSettableType(returns[0], arguments[0]) {
				return ErrReturnValueMustBeSettableAsArgument
			}
			return nil
		},
	}
	sourceFunctionRules = []reflection.TypeValidation{
		func(types []reflect.Type) error {
			switch len(types) {
			case 0:
				return ErrMustExpectLastValueAsArgument
			case 2:
				if !reflection.IsSettableType(actionType, types[1]) {
					return ErrMustBeActionType
				}
			case 3:
				return ErrTooManyArguments
			}
			return nil
		},
	}

	// mutationFunctionRules are simple, has we only expect
	// that a function at least expects one argument which is
	// it's previous state and a interested value of another
	// source or mutation.
	mutationFunctionRules = []reflection.TypeValidation{
		func(types []reflect.Type) error {
			switch len(types) {
			case 0:
				return ErrMustReturnNewStateValue
			case 1:
				return ErrMustBeMoreThanOne
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

// ReducerError defines a type which provides
// a secondary action to be called when a action update
// or a mutation that returns an error, occurs to return
// an error as a signal of failed state, which then
// allows the user to receive error for counter action.
type ReducerError func(err error)

type namespace map[string]string

// DataBind defines the interface which embodies the concept of the
// data system which powers the data pipeline for all Components
// Surfaces, and attributes in slugui.
type DataBind struct {
	namespace  map[string]namespace
	lastValues map[string]interface{}
	lastErrors map[string]interface{}
}

func (db *DataBind) All() map[string]interface{} {
	items := make(map[string]interface{})
	return items
}

func (db *DataBind) Err(targetName string) error {
	return nil
}

func (db *DataBind) Get(targetName string) interface{} {
	return nil
}

func (db *DataBind) Source(sourceName string, fn interface{}, re ...ReducerError) *DataBind     {}
func (db *DataBind) Mutation(mutationName string, fn interface{}, re ...ReducerError) *DataBind {}

func (db *DataBind) Send(actionName string, actionValue interface{}) *DataBind {}
func (db *DataBind) SendTo(actionName string, actionValue string, specificSources ...string) *DataBind {
}
