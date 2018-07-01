package slugui

// BaseSurface returns a new Surface instance using whoes root is
// the provided node.
func BaseSurface(n Node) *Surface {
	return &Surface{Node: n, base: true}
}

// HeadSurface returns a new Surface instance using provided Document is
// Body node as root.
func BodySurface(doc Document) *Surface {
	return &Surface{Document: doc, Node: doc.Body()}
}

// HeadSurface returns a new Surface instance using provided Document is
// Head node as root.
func HeadSurface(doc Document) *Surface {
	return &Surface{Document: doc, Node: doc.Head()}
}

// Surface represents the means by which a giving component interacts
// with it's rendering layer and associated DOM API.
type Surface struct{
	Root *Surface
	Next *Surface
	Prev *Surface
	FirstChild *Surface
	LastChild *Surface
	
	Component Component
	Document Document
	Node Node
	Data DataBind
	Attrs DataBind
	Events EventBind
	base bool
}

// Render updates the internal rendered state of giving surface
// with provided input. Usually this is called by the
// Component to update it's current rendered state.
// Hence the logic implemented within is heavily geared
// towards that assumption, ensuring that the new state
// provided by the call to Render() applies in that
// context:
// Render has the following rules:
// 1. When initially called without any links (i.e a first and last child)
//  render assumes that it is handling an initial state and hence
func (s *Surface) Render(elems ...Element)  error {
	if s.FirstChild == nil && s.LastChild == nil {
		return s.renderAsNew(elems...)
	}
	return s.renderAsUpdate(elems...)
}

func (s *Surface) renderAsNew(elems ...Element)  error {
	for _, elem := range elems {
		es := BaseSurface(s.Node.Root())
		if err := elem(es); err != nil {
			return err
		}
		
		if s.FirstChild == nil && s.LastChild == nil {
			s.FirstChild = es
			s.LastChild = es
			continue
		}
		
		es.Prev = s.LastChild
		s.LastChild.Next = es
		
	}
	return nil
}

func (s *Surface) renderAsUpdate(elems ...Element)  error {
	
	return nil
}

// Use applies provided Component as the core of
// surface operation and interaction.
func (s *Surface) Use(c Component) error {
	if !s.base {
		return ErrSurfaceIsRoot
	}
	
	s.Component = c
	return nil
}

