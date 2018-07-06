SlugUI DataBind is the means through which data can be passed through from any root surface to lower components written by user or by library. DataBind will provide different approaches for data provision which includes Redux, Static and more that may be added later.

Redux
The idea is to take the approach of Redux, simplify and keep essentials which then can be used like in React to provide a single flow of existing data from higher components to lower components without loosing swift efficient updating and change propagation.

Redux in javascript land has 3 very simple concepts has explained below:

* Reducers which I like to call Data Source Reducers (I feel this is better, but because Reducers can also do some work or filtering just like in MapReduce you can understand why Redux uses Reducers), the idea of reducer is that they can the means through which properties declared in the React props will get their value from, such that the reducers can return either a static data or ever changing data values which then becomes the data for properties that React components will use.
* Actions are reducers way of just stating like facts or instructions sent to an existing or create reducer which will return some value in regards to the action description, so we can for example create a reducer that with respect to a action will return some specific data, like say book title or selected book. I want you to think of actions like DOM events on element but this event provide some context that target reducers can use to present specialised data for. When an action gets a result from a data source it targets, the result get’s set as the value for that action and not the data source.
* Connectors, this are really not complex, just like their name, they are the means redux use to bind a Reducer to a component, a component to a action and a action to a reducer. It may seem complicated but just imagine this, that when a action is called to execute, Redux will take the value received from the action and call the reducer with that value and the last value the reducer returned, this way a reducer can either return same value or modify and return another entirely new value based on action data.

Taking in this ideas which are really simple, which unfortunately get communicated very badly with too much big terms, we can reason of a means to do the same things in simpler but effective ways in Slug.

SlugUI:DataBind
We will take the idea of redux simplified to suite a Go idiomatic style of writing code. The slug DataBind has a two simple ideas: 

* Sources
* Mutations

Sources
Source are actions that produce data originally through a http request, static content. They are able to mutate existing state by the call to actions indicating type of behaviour/action to be done with existing data if so. They do not interact with other Sources and are only able to respond to actions and return latest computed value which can either be static and constant or computed based on last action called against. They can be considered partially pure in that they will only mutate due to some action.

Mutations
Mutations are the distant cousins of Sources, their existence is based on derived and computed data from either a single or multiple Sources. They are the Reduce part of the MapReduce paradigm, where they express interest over specific Sources, and when they are called return a computed or computed a new value based off results from said sources or other mutations. Mutations do not require user interactions and will update their latest value when any of their sources or other mutations changes. Even if one item changed, mutations will receive all last known values of existing Sources by which they can behave by.

The access pattern for “Sources” and “Mutations” are exactly alike, they will be treated as keys whoes value will be retrieved once “GET” is called on them. But their definition pattern will defer:

1. Sources only ever receive a “LastState” and optionally with an “Action” object if they state interests in any.
2. Mutations receive a “Bind” which contains all latest states of “Sources” they combine and their last value which they can use for their computation if need be.

To further clarify this, See sample API below:

```
var state = databind.Redux();

state.Source(“visibility”, func(lastState bool) bool {
    switch action.Type{
     case “GET”:
        return true
      case “hide”:
        return false
      case  “show”:
        return true  
    }
    return true
});

state.Source(“class”, func(lastState []string, action Action) ([]string, error) {
    switch action.Type {
     case “AddClass”:
        return append(lastState, action.Data)
     case “GET”:
        return lastState
    }
}).Actions(“AddClass”, “RemoveClass”);

state.Mutation(“visual-style”, func(lastState interface, visiblity bool) (interface{}, error {
    var visibility := db.Get(“visibility”).(bool);
    return Style{display: (visibility ? “inline-block” : “hidden”)}
}).Combines(“visibility”);

state.MutatedAction(“rico”, func(ac Action, lastState interface, visiblity bool, visual-styles interface{}) (interface{}, error {
    var visibility := db.Get(“visibility”).(bool);
    return Style{display: (visibility ? “inline-block” : “hidden”)}
}).Combines(“visibility”, "visual-styles");

// Get last known state of source if available or new data if no last state which sends a GET.
var visibility = state.Get(“visibility”)

// Call Refresh to force an update of a Source that does not expect any action.
// Useful for Sources that may only ever
state.Refresh(“visibility”).Get(“visibility”)

// Have all data sources interested in action update their state.
// Note: Be aware any other source with interest to this source will react also.
state.Send(“hide”, null).Get(“visibility”);

state.Send(“AddClass”, “component-class”).Get(“class");

// Directly send action to only one source.
state.SendSource(“class",“AddClass”, “component-class”);

// Send action to a select few of sources, where:
// First Argument: Action
// Second Argument: Value
// Third Argument: Variadic list of source names. (…string)
state.SendMany(“AddClass”, “component”, “class”, “otherclass”);

// return all current values of both sources and mutations as a map: map[string]interface{}
state.All()
```
