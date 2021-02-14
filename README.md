# Goal
This repository demonstrates how `InboundMiddleware` and `OutboundMiddleware` work in [protoacgor-go](https://github.com/asynkronIT/protoactor-go)'s actor model.

# InboundMiddleware and OutboundMiddleware
See `example/main.go` for regular middleware usage.
This shows how middleware can be defined and how nested middleware chain is called.

Execution of `go run example/main.go` will output log messages somewhat like below:
```
» go run example/main.go                                                                                                                                                                              1 ↵
2019/08/25 12:04:19 InboundMiddleware: start handling incoming message #1: &actor.Started{}
2019/08/25 12:04:19 Nested InboundMiddleware: start handling incoming message: &actor.Started{}
2019/08/25 12:04:19 Actor: received &actor.Started{}
2019/08/25 12:04:19 Nested InboundMiddleware: end handling incoming message: &actor.Started{}
2019/08/25 12:04:19 InboundMiddleware: end handling incoming message #1: &actor.Started{}
2019/08/25 12:04:20 InboundMiddleware: start handling incoming message #2: struct {}{}
2019/08/25 12:04:20 Nested InboundMiddleware: start handling incoming message: struct {}{}
2019/08/25 12:04:20 Actor: received signal
2019/08/25 12:04:20 OutboundMiddleware: start sending message #1 to pong
2019/08/25 12:04:20 OutboundMiddleware: end sending message #1 to pong
2019/08/25 12:04:20 Nested InboundMiddleware: end handling incoming message: struct {}{}
2019/08/25 12:04:20 InboundMiddleware: end handling incoming message #2: struct {}{}
2019/08/25 12:04:20 InboundMiddleware: start handling incoming message #3: &main.pong{}
2019/08/25 12:04:20 Nested InboundMiddleware: start handling incoming message: &main.pong{}
2019/08/25 12:04:20 Actor: received pong
2019/08/25 12:04:20 Nested InboundMiddleware: end handling incoming message: &main.pong{}
2019/08/25 12:04:20 InboundMiddleware: end handling incoming message #3: &main.pong{}
2019/08/25 12:04:21 InboundMiddleware: start handling incoming message #4: struct {}{}
2019/08/25 12:04:21 Nested InboundMiddleware: start handling incoming message: struct {}{}
2019/08/25 12:04:21 Actor: received signal
2019/08/25 12:04:21 OutboundMiddleware: start sending message #2 to pong
2019/08/25 12:04:21 OutboundMiddleware: end sending message #2 to pong
2019/08/25 12:04:21 Nested InboundMiddleware: end handling incoming message: struct {}{}
2019/08/25 12:04:21 InboundMiddleware: end handling incoming message #4: struct {}{}
2019/08/25 12:04:21 InboundMiddleware: start handling incoming message #5: &main.pong{}
2019/08/25 12:04:21 Nested InboundMiddleware: start handling incoming message: &main.pong{}
2019/08/25 12:04:21 Actor: received pong
2019/08/25 12:04:21 Nested InboundMiddleware: end handling incoming message: &main.pong{}
2019/08/25 12:04:21 InboundMiddleware: end handling incoming message #5: &main.pong{}
2019/08/25 12:04:22 InboundMiddleware: start handling incoming message #6: &actor.Stopping{}
2019/08/25 12:04:22 Nested InboundMiddleware: start handling incoming message: &actor.Stopping{}
2019/08/25 12:04:22 Actor: received &actor.Stopping{}
2019/08/25 12:04:22 Nested InboundMiddleware: end handling incoming message: &actor.Stopping{}
2019/08/25 12:04:22 InboundMiddleware: end handling incoming message #6: &actor.Stopping{}
2019/08/25 12:04:22 InboundMiddleware: start handling incoming message #7: &actor.Stopped{}
2019/08/25 12:04:22 Nested InboundMiddleware: start handling incoming message: &actor.Stopped{}
2019/08/25 12:04:22 Actor: received &actor.Stopped{}
2019/08/25 12:04:22 Nested InboundMiddleware: end handling incoming message: &actor.Stopped{}
2019/08/25 12:04:22 InboundMiddleware: end handling incoming message #7: &actor.Stopped{}
2019/08/25 12:04:22 Finish
```

# Plugin
While `example/main.go` shows how a regular middleware works, `plugin/main.go` describes how a middleware construction with [`plugin`](https://github.com/AsynkronIT/protoactor-go/blob/bab29b9c335d46ad2f0ffaafdec40300eea16115/plugin/plugin.go) mechanism runs.
This kind of middleware construction is especially used in `protoactor-go` to build setup a [`PassivationPlugin`](https://github.com/AsynkronIT/protoactor-go/blob/bab29b9c335d46ad2f0ffaafdec40300eea16115/plugin/passivation.go).
`PassivationPlugin` is a designated middleware that runs a initialization code on actor start and runs another logic on other message receptions.
This is mainly used to let a cluster grain start a timer on actor initialization, reset the timer on every message reception and stop the actor when no message comes before the timer ticks.
Remember that, in protoactor's cluster grain architecture, Cluster Grains always "exist."
When no Cluster Grain actor exists on message reception, an actor is created; when no message follows during a pre-defined interval, the actor stops to save server resources.

Execution of `go run plugin/main.go` will output log messages somewhat like below:
```
» go run plugin/main.go                                                                                                                                                                               1 ↵
2019/08/25 12:17:18 Middleware: starting $1
2019/08/25 12:17:18 Actor: received &actor.Started{}
2019/08/25 12:17:18 Middleware: received &actor.Stopping{}
2019/08/25 12:17:18 Actor: received &actor.Stopping{}
2019/08/25 12:17:18 Middleware: received &actor.Stopped{}
2019/08/25 12:17:18 Actor: received &actor.Stopped{}
```

# References
- [[Golang] Protoactor-go 101: Introduction to golang's actor model implementation](https://blog.oklahome.net/2018/07/protoactor-go-introduction.html)
- [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html)
- [[Golang] protoactor-go 101: How actor.Future works to synchronize concurrent task execution](https://blog.oklahome.net/2018/11/protoactor-go-how-future-works.html)
- [[Golang] protoactor-go 201: How middleware works to intercept incoming and outgoing messages](https://blog.oklahome.net/2018/11/protoactor-go-middleware.html)
- [[Golang] protoactor-go 201: Use plugins to add behaviors to an actor](https://blog.oklahome.net/2018/12/protoactor-go-use-plugin-to-add-behavior.html)

# Other Example Codes
- [oklahomer/protoactor-go-sender-example](https://github.com/oklahomer/protoactor-go-sender-example)
  - Some example codes to illustrate how protoactor-go refers to sender process
- [oklahomer/protoactor-go-future-example](https://github.com/oklahomer/protoactor-go-future-example)
  - Some example codes to illustrate how protoactor-go handles Future process
