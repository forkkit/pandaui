package elementals

import "gitlab.com/slugui/slugui"

func Div(elems ...slugui.Elemental) slugui.Elemental {
	return func(surface slugui.Surface) slugui.Element {
		var elem slugui.Element
		elem.TagName = "div"
		
	}
}