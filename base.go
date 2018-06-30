package slugui


// HTML defines a type to represent a rendered html markup.
type HTML string

// EventBind defines a interface to be implemented for
// event watching, interaction and removal.
type EventBind interface{
	// Remove removes event name from watched list.
	Remove(string)
	
	// Has returns true/false if giving event is being watched.
	Has(string) bool
	
	// Add adds giving event into watching list with associated
	// states.
	Add(name string, preventDefault bool, stopPropagation bool)
}

type DataBind interface{}

type AttrBind interface{}

// Tag defines an interface which defines a single method
// that returns desired tag name.
type Tag interface{
	Tag() string
}

// Renderable defines a type which exposes a method that
// receives a Surface and applies operation to link giving
// surface with new content and children.
type Renderable interface{
	Render(Surface) error
}

// Mountable defines a type which exposes method to be called
// based on whether the state is mounted or unmounted.
type Mountable interface{
	Mounted()
	Unmounted()
}

// Inited defines existing states to define the
// stage of a implement type whether be it initialized or
// destroyed.
type Inited interface{
	Init()
	Destroy()
}

// Component defines a custom type which embodies
// the content to be rendered and how such a
// content would be interacted with.
type Component interface{
	Inited
	Mountable
	Renderable
}

// Elemental defines a associated type which giving
// a provided surface will return an Element type.
type Elemental func(Surface, ...Elemental) Surface

// NSurface defines the interaction layer which handles the update
// and rendering sequences for handling defined component data and
// event management.
//type NSurface interface {
//	Node() Node
//
//	Data() DataBind
//	Attr() AttrBind
//	Events() EventBind
//
//	Upgrade(Surface) Surface
//	Use(Component) error
//	Render(...Elemental) error
//}

// Surface represents the means by which a giving component interacts
// with it's rendering layer and associated DOM API.
type Surface struct{
	Root *Surface
	Next *Surface
	LastChild *Surface
	FirstChild *Surface
	
	RootNode Node
	Data DataBind
	Attrs AttrBind
}

// NewSurface returns a new Surface instance using the provided node as
// it's root source.
func NewSurface(root Node) *Surface {
	return &Surface{RootNode: root}
}



