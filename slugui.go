package slugui

import "gitlab.com/slugui/slugui/dom"

// HTML defines a type to represent a rendered html markup.
type HTML string

type DataBind interface{}

type AttrBind interface{}

// Renderable defines a type which exposes a method that
// receives a series of elementals to be rendered.
type Renderable interface{
	Render(Surface) error
}


// Component defines a custom type which embodies
// the content to be rendered and how such a
// content would be interacted with.
type Component interface{
	Initialize()
	Destroy()
	Mounted()
	Unmounted()
	Render(Surface) error
}

// Surface defines the interaction layer which handles the update
// and rendering sequences for handling defined component data and
// event management.
type Surface interface {
	Node() Node
	Data() DataBind
	Attr() AttrBind
	
	Use(Component)
	Render() HTML
}


// SurfaceTransformer defines a function type which takes in a giving surface
// and returns another surface with either limited data scope
type SurfaceTransformer func(Surface) Surface

// Elemental defines a associated type which giving
// a provided surface will return an Element type.
type Elemental func(Surface) Element

// Element defines the description of a specific type
// representing a DOM node. It holds associated list of
// dom events for giving element and tag name represent
// the specific type. It's up to the renderer to define how
// it wishes to render said elements.
type Element struct {
	TagName string
	Events []string
}

// Render returns giving html with provided children content
// as a rendered string type of HTML.
func (e Element) Render(children []HTML) HTML {
	return ""
}