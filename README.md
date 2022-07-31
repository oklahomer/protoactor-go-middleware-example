# Goal
This repository demonstrates how `InboundMiddleware` and `OutboundMiddleware` work in [protoacgor-go](https://github.com/asynkron/protoactor-go)'s actor model.

# InboundMiddleware and OutboundMiddleware
See `example/main.go` for the regular middleware usage.
This shows how middleware can be defined and how nested middleware chain is called.

The execution of `go run example/main.go` will output log messages somewhat like the below:
```
» go run example/main.go
2022/07/31 12:54:55 ReceiverMiddleware: start handling incoming message #1: &actor.Started{}
2022/07/31 12:54:55 Nested ReceiverMiddleware: start handling incoming message: &actor.Started{}
2022/07/31 12:54:55 Actor: received &actor.Started{}
2022/07/31 12:54:55 Nested ReceiverMiddleware: end handling incoming message: &actor.Started{}
2022/07/31 12:54:55 ReceiverMiddleware: end handling incoming message #1: &actor.Started{}
2022/07/31 12:54:56 ReceiverMiddleware: start handling incoming message #2: struct {}{}
2022/07/31 12:54:56 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
2022/07/31 12:54:56 Actor: received signal
2022/07/31 12:54:56 SenderMiddleware: start sending message #1 to pong
2022/07/31 12:54:56 SenderMiddleware: end sending message #1 to pong
2022/07/31 12:54:56 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
2022/07/31 12:54:56 ReceiverMiddleware: end handling incoming message #2: struct {}{}
2022/07/31 12:54:56 ReceiverMiddleware: start handling incoming message #3: &main.pong{}
2022/07/31 12:54:56 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
2022/07/31 12:54:56 Actor: received pong
2022/07/31 12:54:56 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
2022/07/31 12:54:56 ReceiverMiddleware: end handling incoming message #3: &main.pong{}
2022/07/31 12:54:57 ReceiverMiddleware: start handling incoming message #4: struct {}{}
2022/07/31 12:54:57 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
2022/07/31 12:54:57 Actor: received signal
2022/07/31 12:54:57 SenderMiddleware: start sending message #2 to pong
2022/07/31 12:54:57 SenderMiddleware: end sending message #2 to pong
2022/07/31 12:54:57 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
2022/07/31 12:54:57 ReceiverMiddleware: end handling incoming message #4: struct {}{}
2022/07/31 12:54:57 ReceiverMiddleware: start handling incoming message #5: &main.pong{}
2022/07/31 12:54:57 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
2022/07/31 12:54:57 Actor: received pong
2022/07/31 12:54:57 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
2022/07/31 12:54:57 ReceiverMiddleware: end handling incoming message #5: &main.pong{}
2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #6: struct {}{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
2022/07/31 12:54:58 Actor: received signal
2022/07/31 12:54:58 SenderMiddleware: start sending message #3 to pong
2022/07/31 12:54:58 SenderMiddleware: end sending message #3 to pong
2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #6: struct {}{}
2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #7: &main.pong{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
2022/07/31 12:54:58 Actor: received pong
2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #7: &main.pong{}
2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #8: &actor.Stopping{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopping{}
2022/07/31 12:54:58 Actor: received &actor.Stopping{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopping{}
2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #8: &actor.Stopping{}
2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #9: &actor.Stopped{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopped{}
2022/07/31 12:54:58 Actor: received &actor.Stopped{}
2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopped{}
2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #9: &actor.Stopped{}
2022/07/31 12:54:58 Finish
```

# Plugin
While `example/main.go` shows how a regular middleware works, `plugin/main.go` describes how a middleware construction with [`plugin`](https://github.com/asynkron/protoactor-go/blob/afd2d973a1d1/plugin/plugin.go) mechanism runs.
This kind of middleware construction is especially used in `protoactor-go` to build setup a [`PassivationPlugin`](https://github.com/asynkron/protoactor-go/blob/afd2d973a1d1/plugin/passivation.go).
`PassivationPlugin` is a designated middleware that runs a initialization code on actor start and runs another logic on other message receptions.
This is mainly used to let a cluster grain start a timer on actor initialization, reset the timer on every message reception and stop the actor when no message comes before the timer ticks.
Remember that, in protoactor's cluster grain architecture, Cluster Grains always "exist."
When no Cluster Grain actor exists on message reception, an actor is created; when no message follows during a pre-defined interval, the actor stops to save server resources.

Execution of `go run plugin/main.go` will output log messages somewhat like below:
```
» go run plugin/main.go
2022/07/31 12:44:39 Middleware: starting $1
2022/07/31 12:44:39 Actor: received &actor.Started{}
2022/07/31 12:44:39 Middleware: received &actor.Stopping{}
2022/07/31 12:44:39 Actor: received &actor.Stopping{}
2022/07/31 12:44:39 Middleware: received &actor.Stopped{}
2022/07/31 12:44:39 Actor: received &actor.Stopped{}
```

# References
- [[Golang] Protoactor-go 101: Introduction to golang's actor model implementation](https://blog.oklahome.net/2018/07/protoactor-go-introduction.html)
- [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html)
- [[Golang] protoactor-go 101: How actor.Future works to synchronize concurrent task execution](https://blog.oklahome.net/2018/11/protoactor-go-how-future-works.html)
- [[Golang] protoactor-go 201: How middleware works to intercept incoming and outgoing messages](https://blog.oklahome.net/2018/11/protoactor-go-middleware.html)
- [[Golang] protoactor-go 201: Use plugins to add behaviors to an actor](https://blog.oklahome.net/2018/12/protoactor-go-use-plugin-to-add-behavior.html)
- [[Golang] protoactor-go 301: How proto.actor's clustering works to achieve higher availability](https://blog.oklahome.net/2021/05/protoactor-clustering.html)

# Other Example Codes
- [oklahomer/protoactor-go-sender-example](https://github.com/oklahomer/protoactor-go-sender-example)
  - Some example codes to illustrate how protoactor-go refers to sender process
- [oklahomer/protoactor-go-future-example](https://github.com/oklahomer/protoactor-go-future-example)
  - Some example codes to illustrate how protoactor-go handles Future process
