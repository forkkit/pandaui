package slugui

import "time"

// Event Body represents basic data provided information
// related to an event occurring within dom.
type EventBody struct {
	Phase int
	Type string
	Target Node
	Current Node
	Bubbles bool
	FromAgent bool
	Time time.Time
	Attached interface{}
	Detail map[string]interface{}
}

// NodeEvent defines the type of Event that would occur
// for a giving node after the execution of associated event type.
type NodeEvent interface{
	Body() EventBody
	PreventDefault()
	StopPropagation()
}

// EventSubscriber defines a type which exposes a single method
// to end a giving subscription between itself and a object.
type EventSubscriber interface{
	Unsubscribe()
}

// EventPublisher defines an interface which represents a pub-sub
// system for node structures.
type EventPublisher interface{
	Subscribe(string) EventSubscriber
	Trigger(string, interface{}, map[string]interface{})
}

// KV defines an interface of a type for storing KV pairs.
type KV interface{
	Get(string) string
	Set(string,string)
	Map() map[string]string
}

// Queryable defines a interface to represent a implement type
// which exposes methods that allows retrieval of node/nodes
// by an associated query.
type Queryable interface{
	One(string) (Node, error)
	All(string) ([]Node, error)
}

// Document defines a root structure which expresses the DOM document
// API element for retrieving initial nodes for the head or body
// element of a DOM document.
type Document interface{
	EventPublisher
	
	Head() Node
	Body() Node
}

// Node defines an interface which exposes a type of DOM element created
// by the Document implementing type. It is the proxy through which DOM
// behaviour is interacted with by outside element and is implemented by
// those providing DOM or DOM-like systems.
type Node interface{
	KV
	Tag
	EventPublisher
	
	// CreateAndAppend returns a new Node instance which exists as
	// a child to the called Node. It defers from Node.Create in that
	// the node is instantly appended.
	CreateAndAppend(string) (Node, error)
	
	// Create returns a new Node instance which is not attached to
	// this node. It provides a method through which all node types
	// can create Node instances from. The node is still separate
	// from the Document till attached.
	Create(string) (Node, error)
	
	// Append adds provided Node as child of called Node. It must return
	// an error to signal success or failure.
	Append(Node) error
	
	// Render must overwrite existing content using provided
	// HTML as new content of itself.
	Render(HTML) error
}