package slugui

// error constants.
const (
	ErrSurfaceIsRoot = Error("Surface is a root node, can not be used for Components")
)

// Error implements a custom error type for package to
// provide consistency and constancy of type and value.
type Error string
func (e Error) Error() string {
	return string(e)
}

// HTML defines a type to represent a rendered html markup.
type HTML string

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

// Element defines a function which takes a Surface for
// internal operation and application.
type Element func(*Surface) error

// Elemental defines a associated type which giving
// a provided surface will return an Element type.
type Elemental func(...Element) Element

