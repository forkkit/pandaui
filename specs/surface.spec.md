SlugUI:Surface
-----------------
Surface is an abstraction that sits above user component and interests with any rendering layer through RenderingProvider abstraction or interface. It is the container that has the incoming data from a parent or from the root and also contains has fields other needed instantiated fields for interaction.


 ## Usage of Surface

```
var container := Surface{
 Data: DataBind     // Provided by user to supply data to elements.
 State: StateBind   // Unique to a element and surface.
 DOM: DOMProvider.Node // Means of interaction by surface with representation. Also provides []Children.
} 

Surface.Render(DOM);
```

A surface exposes the contents to be provided to a UI component and a Render method to receive latest updates for the view. 
I think it is more wise to follow this approach than with respect to Gu's render interface approach which had certain issues:

```
Render() *tree.Markup
```

More so, we must also change the pattern on how elements markup are created. From now will just express markup has functions return constructors. 

```
function Div() func(Properties) {...}
```

As this provides a great level of control on how properties are passed down into create elements or other components. If you reason deeply the issue with most libraries is the fact that it's hard to mutate an existing dom element safely because ypu need to account for consistency which then bring about use of mutexes. But if we restrict and provide such functional interface that create element after provision of properties we can easily overcome said issue. 

### Component with Surface: 

```
var _ = ui.Register(Card); // should register type to tag ‘Card’
Var _ = ui.Register(Card, “ICard”) // should register type to tag ‘ICard'

type Card struct{
  Surface Surface
}

func (c *Card) Mounted(){
  c.Surface.StateBind.SetState();
}

func (c *Card) Unmounted(){}

func (c *Card) Clicked(){
  c.Surface.SetState();
}

func CardComponent(Surface) {
    c := new(Card)
    c.Surface = s;
    
    // Supply to a surface it’s handler or scope.
    // The surface must be aware of underline 
    // source which it is to represent and interact with.
    // The Surface must also handle registration of
    // 1. Render
    // 2. Mount
    // 3. Unmount
    // If handler meets interface requirement.
    //
    s.use(c);
}

func (c *Card) Render(s Surface) error {
  return s.Render(ui.Div(
        ui.Text(“do something”),
        ui.Event(“click”, c.clicked).preventDefault(),
        c.Surface.Children(),
  ), ui.Div(
      ui.Text(“Run away”),
      surface.Transform(func (root Surface) ui.Element {
        newSurface := Surface.WithData(root, root.Data.Get(…));
        return ui.Anchor(newSurface);
      }),
  ));
}
```

Where `ui.Div` returns:

```
func(s Surface) ui.Element
```

The `Surface` passed is a new Surface instance specifically linked to parent/root Surface but caters only for Div.

A surface exists to be that direct link for a defined component in the sense that it's an extension from a root which creates a existing tree of links, and each uniquely caters to the need of that component that is registered to it e.g. input tag getting it's value property. 

Using this approach raise an issue you get throwable types e.g. Surfaces and DOM elements representation. 

#### Why
Well Slug adheres to the notion that state is embedded into representation and is decoupled from the code handling it. See States and Events.

More so, when a component wants to update it’s rendered expression/representation, the new structures returned will simple be swap out in the surface using the `Use` directive if  of similar tag name else a ore complex repositioning would need be done, of course for basic element types eg Div and not user defined.

In Gu though w3 return a representation of what is needed but these representations carry states that's difficult to get rid of and make them reusable but if we truly build stateless code like above where the state is stored as part of the representation and our code only manages behaviour based on available state then there is a win win situation because handlers can be recycled, more so surfaces in the sense can also be recycled.

#### Reconciliation of Surface DOM
Gu and other libraries there is a form of recoiliation that happens at the root when all elements created get a diffed or update effienciently without touching rendering before an update to reduce works during rendering like buffered images in java where we render first to a buffer then replace said screen with buffer contents.

I suggest with move reconciliation from the root down to the localized levels, where each surface is responsible for appropriately it's cummincating a chqnge to the representation layer in a common format which then can be used to efficiently. 

In Gu we  had to use a 2 step strategy due to ensuring that elements in markup that are added by user do not get removed and can be maintained during updates. This I must confess is problematic, if we restrict dom creation to only ever come from within user code then we can skip a 2-step strategy and use a one-step strategy which would be that when a update has arrive, we can efficiently swap out just that area alone with the new copy without running through diffing or a line by line, one-by-one removal of tagged elements and addition of  new elements from representation, this of course is predicate on the assumption that the parent is very much able to let lower parts ideally relate through some means with representation without tight coupling.

#### States and Events
More so, surfaces must be able to take in existing representation and their states as if they were previous state as this allows slugs to immediately without replacing dom to work with existing contents.

Even more slugs UI elements must do their best to only have their states as part of representation and not as internal code level that makes it hard for transition although this should not be impossible if desired.

More so slugs surfaces should be able to infer capabilities and events from existing representation of elements through attributes. For example with events we can do the following:

```
<div slug-events=“click{nopg, dopd} touch mouse-over”></div>
```

By doing this we can essentially ensure that dom elements themselves are the very descriptors of the type of events they wish to be alerted about and by using the format 
EventName{...}

The contents with brackets will contain capabilities like no propagation using `nopg` and prevent default using `dopd` to affect behaviour of caller , we gain easy swift but reconcilable means of handling events without evening binding directly to elements using the old jquery’s idea of live events, where a constant parent element such as the body or the element which get’s mounted upon will be the means through which events that bubble up are caught and can be inspected in regards to a target, where said target has our events `slug-events` attribute to indicate if it wishes notification on such an event.

More so, if we shift our minds to states to follow this way, that may be widely important to a element especially one that may jump in between server render and client render, we can use a similar approach to the explanation of events to let the dom elements attributes contain state values which allows us quickly swap in new instances of code to take over elements lifecycle without much problems  or  loss of last state that can occur during the transition from a content rendered on the server to that handled on the client.

```
<div slug-state-last-id=1 slug-state-drug-count=100></div>
```

In this example we can have the original states of `last-id` and `drug-count` inferred back into an instantiated elements from the rendered representation because these have being persisted as part of the element properties as well. Of course, this can not be used for sensitive states that may have security concerns which then would require more sophisticated techniques like pulling from a secure remote resource or using indexedDB local storage provided by browser. But by taking this types of capabilities into account we can fairly quickly hydrate and dehydrate states with little effort has a Surface must handle this for users.

#### Document and Nodes
Node exists from a Document and a Document is the only concept that connects a surface to an underline rendering system. Without a Document. There is no node. In the interaction of a Surface and a Document, only the root document will ever have access to a Doucment, no lower Surface that lies beneath, is a child of another should directly convert with the document. The root surface must handle the responsibility of serializing the whole structure into a form it can give to the document for immediate update and must be the means by which it and it's children retrieve such their respective node representation.

Browser DOM wise, a slug UI concept of document does not necessary match a one to one relation with the browser document object. For all purposes browser wise the element that is used as a Document with a surface can be a node within the document, the body or head element. In the interaction children nodes retrieve their respective nodes through query and are done lazily based on initial call. Data or values of nodes must be repeated and left to the data layer to retrieve such from.

Outside the document used for a surface must come from a RootDocument which has a one to one relationship to the browser or headless DOM document.