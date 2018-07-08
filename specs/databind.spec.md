SlugUI DataBind
-----------------

SlugUI DataBind is the means through which data can be passed through from any root surface to lower components written by user or by library. DataBind will provide different approaches for data provision which includes Redux, Static and more that may be added later.

## Redux
The idea is to take the approach of Redux, simplify and keep essentials which then can be used like in React to provide a single flow of existing data from higher components to lower components without loosing swift efficient updating and change propagation.

Redux in javascript land has 3 very simple concepts has explained below:

- Reducers which I like to call Data Source Reducers (I feel this is better, but because Reducers can also do some work or filtering just like in MapReduce you can understand why Redux uses Reducers), the idea of reducer is that they can the means through which properties declared in the React props will get their value from, such that the reducers can return either a static data or ever changing data values which then becomes the data for properties that React components will use.
- Actions are reducers way of just stating like facts or instructions sent to an existing or create reducer which will return some value in regards to the action description, so we can for example create a reducer that with respect to a action will return some specific data, like say book title or selected book. I want you to think of actions like DOM events on element but this event provide some context that target reducers can use to present specialised data for. When an action gets a result from a data source it targets, the result get’s set as the value for that action and not the data source.
- Connectors, this are really not complex, just like their name, they are the means redux use to bind a Reducer to a component, a component to a action and a action to a reducer. It may seem complicated but just imagine this, that when a action is called to execute, Redux will take the value received from the action and call the reducer with that value and the last value the reducer returned, this way a reducer can either return same value or modify and return another entirely new value based on action data.

Taking in this ideas which are really simple, which unfortunately get communicated very badly with too much big terms, we can reason of a means to do the same things in simpler but effective ways in Slug.

## SlugUI:DataBind

Databind approach takes inspiration from Redux and creates it's own concepts to power the means by which data is interacted,
retrieved and computed through. By simplifying to these, we hope to create a robust and flexible means to easily update 
data values with ease and in sharing to different parts of the application with easy and simply flow control. 


### Feeds

Feeds are pure sources of data which can be affected by their previous value or a user provided Action containing some data that provides some context to what its data should be like or form. Feeds can be purely static or even only ever return data without caring about Actions from users, they can also be purely action based so that the data they return is only ever supplied based on user Action context. The basic idea is they are sources, they retrieve that data from some unknown true source and only ever are interacted with through Actions from users.

Feeds are perfect for data that are retrieved from a http request based on some context or static data read from a file system. They will only ever produce data and not derive data from some other data value or key. How they do what they do and where they get their data from is entirely left to their implementation.

```go
var state = databind.Redux();

state.Source("random_value", func() int {
    return rand.Intn()
});

state.Source("visibility" func(lastState bool) bool {
    switch action.Type{
     case "GET”:
        return true
      case "hide”:
        return false
      case  "show”:
        return true  
    }
    return true
}, "hide");
```

### DerivedFeeds

DerivedFeeds are the very opposite of Feeds, which are pure and externally retrieve. Derived feeds are values derived from existing feeds or other derived feeds, they are the merged or computed values from final values of others. They exists to be some derived 
form that exists because of changes on other values. They are like streams map or reduce operation that take a series of values deriving a new unique value from their sources. They are perfect for values that must shift their based on others, like visiblity or self updating multiplications.

```go
var state = databind.Redux();

state.Mutation("visual-style" func(lastState interface, visiblity bool) (interface{}, error) {
    var visibility := db.Get("visibility”).(bool);
    return Style{display: (visibility ? "inline-block” : "hidden”)}
}, "visibility”);

```

### Derived ActionFeeds

Derived ActionFeeds are very similar to DerivedFeeds except that they are values that are computed based on a giving Action which provides context for how the values they are interested in should be merged/reduced or computed by. They exists to be a sort of derived data based on user defined rules. Derived Feeds are a merging based on change of existing values, whilst Derived ActionFeeds are a merging based on user's action that has a direct effect on their returned values.


```go
state.MutatedAction("rico" func(ac Action, lastState interface, visiblity bool, visual-styles interface{}) (interface{}, error) {
    var visibility := db.Get("visibility”).(bool);
    return Style{display: (visibility ? "inline-block” : "hidden”)}
}, "visibility", "visual-styles");
```


In databind, values are not proactively calculated has this may cause heavy peak of resource usage, instead derived feeds are lazily computed based on-demand computation where user request's value. This way if a computation's key is never called for value retrieval then no effort will be executed to have it's latest value computed.

### Values

Retreiving values for computed values should be easy and we simply use a `Getter` method that returns latest values either previously computed or computed on demand when retrieval is called.
 
```go
// Get last known state of source if available or new data if no last state which sends a GET.
var visibility = state.Get("visibility”)
```


### Actions

Actions are the user provide context which are used by DerivedActionFeeds and Feeds, these allows users to provide some information that can be used in desirable context by the target feed to provide a personalized response/value to the user. This ensures we can allow user interaction with computed values but in a safe and easy means that allows us both control and cascade update notification easily to other data parts that may have interest in computed value.


```go
value, err := state.Compute("AddClass", "component-class") 
```


```go
newValue, err := state.SendSource("class","AddClass" "component-class”);
```

### Derived Values

Databind must provide a means to convert computed values from it's internal store into a map that contains unique, copy of existing data that maintains the immutability of it's stored data. This allows users gain access to computed state or values easily as a holistic view of all keys.

```go
state.All()
```
