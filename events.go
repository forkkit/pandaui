package slugui

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

