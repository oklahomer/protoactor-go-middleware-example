package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/plugin"
	"log"
)

type myPlugin struct {
}

func (_ *myPlugin) OnStart(ctx actor.ReceiverContext) {
	log.Printf("Middleware: starting %s\n", ctx.Self().Id)
}

func (_ *myPlugin) OnOtherMessage(_ actor.ReceiverContext, env *actor.MessageEnvelope) {
	_, msg, _ := actor.UnwrapEnvelope(env)
	log.Printf("Middleware: received %#v\n", msg)
}

func main() {
	// Set up the actor system
	system := actor.NewActorSystem()

	// Construct a plugin implementation
	myPlugin := &myPlugin{}

	// Wrap the plugin implementation in a form of InboundMiddleware
	middleware := plugin.Use(myPlugin)

	// 2022/07/31 12:44:39 Middleware: starting $1
	// 2022/07/31 12:44:39 Actor: received &actor.Started{}
	// 2022/07/31 12:44:39 Middleware: received &actor.Stopping{}
	// 2022/07/31 12:44:39 Actor: received &actor.Stopping{}
	// 2022/07/31 12:44:39 Middleware: received &actor.Stopped{}
	// 2022/07/31 12:44:39 Actor: received &actor.Stopped{}
	props := actor.PropsFromFunc(
		func(ctx actor.Context) {
			log.Printf("Actor: received %#v\n", ctx.Message())
		},
		actor.WithReceiverMiddleware(middleware), // Set as a middleware
	)

	pid := system.Root.Spawn(props)
	system.Root.Send(pid, "dummy message")
	system.Root.StopFuture(pid).Wait() // Waits till the target actor actually stops
}
