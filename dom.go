package slugui

import "time"



// Event Body represents basic data provided information
// related to an event occuring within dom.
type EventBody struct {
	Type string
	Target Node
	Current Node
	FromAgent bool
	Time time.Time
	Phase int
	Detail map[string]interface{}
	Attached interface{}
}

// NodeEvent defines the type of Event that would occur
// for a giving node after the execution of associated event type.
type NodeEvent interface{
	PreventDefault();
	StopPropagation();
	Body() EventBody
}

// Node defines a rendering node which represent the tag and type
// with associated attributes and methods.
type Node interface{
	Tag() string
	Attr(string) string
	SetAttr(string,string)
	Attrs() map[string]string
	SetAttrs(map[string]string)

	Append(Node) error
	Create(string) (Node, error)
	Listen(eventName string)
	Unlisten(eventName string)

	One(string) (Node, error)
	All(string) ([]Node, error)
}
