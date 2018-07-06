SlugUI is a framework of libraries that are to be created to provide easy and simple development of full stack web applications using web assembly, service workers and webworkers with DOM manipulation to provide a simple and easy to understand, integrated solution for building web applications able to work on both server setup and in browsers. They must embody the idea of isomorphic ability and similarity to development style used on the server with Go, but suited for the client side. We need to make the learning barrier has smooth and easy as possible.

SlugUI is the cumulative work done on GU (https://github.com/gu-io) and other experiments that mimics the beauty of passed down parameters or props like React, with simple easy intelligent rendering and efficient diffing, with integrated server request-response processing for routing that allows easy transition both on page routing and resource provision, basically build a similar idea to golang http routers but for the client without much difference in style.

More so, ensure we intercept all requests and integrate page loading like is done with a library like Turbolinks. 

SlugUI will have different libraries providing specific layers for use:

* SlugUI:Surface 
    Will be the implementation of all things rendering and interfacing with the DOM. In Slug there will exists two implementation of the DOM where one will be used in browser and another within go where a browser is not available by taking advantage of Golang’s html parsing to generate said DOM, will will take a minimalistic approach to this and borrow ideas from jsdom but only implement (at least initially) rendering and event handling without going into other DOM APIs. This makes it easy to do BDD tests and more so server-side rendering without effort.

* SlugUI:DataBind
DataBind will be the data layer where data and context is passed into developed slug components, it will borrow ideas from Redux and Stimulus (which prefers data existing as attributes on DOM elements, which is a idea I like very much especially in regards to hydration and dehydration of state for a dom element, of course other approaches will be envisioned as a companion to this).

* SlugUI:HTTP
A client side http like library especially made for routing and request handling. Users should not have much mental model shift in the way the write request-response, of course what they can use and how they use it may differ but the experience should be as close and similar as possible, more so, we must include the ability for parts of rendered DOM should be able to take response from http by SlugUI:Surface and seamlessly swap them into rendering. See Turbolinks for this (I really like this).

I will be heavily inspired from the following libraries and ideas:

* https://github.com/stimulusjs/stimulus
* https://github.com/turbolinks/turbolinks
* https://gitlab.com/gu-io/gu
* https://go.isomorphicgo.org/go/
* https://github.com/jsdom/jsdom

Talked with Geofrey Ernest(https://github.com/gernest) from Gophers Slack pointed me to a nice UI language specification that seems nice, we should take inspiration on regards ready made components and design system from them, see:

* https://ant.design/
* https://github.com/ant-design?page=2

I have taking a heavily like to actors as a good concurrency model and also in general as a good means of expressing logic between different parts due to it’s native message parsing construct, I fell by chance on two projects from Google’s ChromeLabs(https://github.com/GoogleChromeLabs) that use web workers for interesting work:

* https://github.com/GoogleChromeLabs/clooney.git
* https://github.com/GoogleChromeLabs/comlink
* https://github.com/GoogleChromeLabs/comlink-loader
* https://github.com/GoogleChromeLabs/tasklets
* https://github.com/GoogleChromeLabs/application-shell

Go has Actors implemented with a interesting approach see:
    
* https://github.com/AsynkronIT/protoactor-go

Have also have ongoing work in Actors:

* https://gitlab.com/gokit/actorkit


Although actors provide interesting, they will for now take a back sit, they are included here for future indicators of ideas to chase later, but for now I will chase a pristine request-response approach with SlugUI:HTTP, this will allow later providing some other means e.g SlugUI:Actors that can be used by developers.

Also as for web assembly, see following for more inspiration:

* https://github.com/aspnet/Blazor
* https://github.com/gowasm/gopherwasm
* https://github.com/johanbrandhorst/wasm-experiments/tree/master/grpc

Note: I will be moving the work on this to Gitlab (not for any special reason but for the free private repos they offer without charge).

Repo: https://gitlab.com/slugui.