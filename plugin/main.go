package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/plugin"
	"log"
)

type myPlugin struct {
}

func (_ *myPlugin) OnStart(ctx actor.Context) {
	log.Printf("Middleware: starting %s\n", ctx.Self().Id)
}

func (_ *myPlugin) OnOtherMessage(ctx actor.Context, msg interface{}) {
	log.Printf("Middleware: received %#v\n", msg)
}

func main() {
	// Construct plugin implementation
	myPlugin := &myPlugin{}

	// Wrap plugin implementation in a form of InboundMiddleware
	middleware := plugin.Use(myPlugin)

	// 2019/08/25 12:17:18 Middleware: starting $1
	// 2019/08/25 12:17:18 Actor: received &actor.Started{}
	// 2019/08/25 12:17:18 Middleware: received &actor.Stopping{}
	// 2019/08/25 12:17:18 Actor: received &actor.Stopping{}
	// 2019/08/25 12:17:18 Middleware: received &actor.Stopped{}
	// 2019/08/25 12:17:18 Actor: received &actor.Stopped{}
	props := actor.
		FromFunc(func(ctx actor.Context) {
			log.Printf("Actor: received %#v\n", ctx.Message())
		}).
		WithMiddleware(middleware) // Set as a middleware

	pid := actor.Spawn(props)
	pid.Tell("dummy message")
	pid.GracefulStop() // Waits till the target actor actually stops
}
