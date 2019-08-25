package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newInboundMiddleware() actor.InboundMiddleware {
	cnt := 0
	return func(next actor.ActorFunc) actor.ActorFunc {
		return func(ctx actor.Context) {
			cnt++
			log.Printf("InboundMiddleware: start handling incoming message #%d: %#v", cnt, ctx.Message())
			next(ctx)
			log.Printf("InboundMiddleware: end handling incoming message #%d: %#v", cnt, ctx.Message())
		}
	}
}

func newOutboundMiddleware() actor.OutboundMiddleware {
	cnt := 0
	return func(next actor.SenderFunc) actor.SenderFunc {
		return func(ctx actor.Context, target *actor.PID, envelope *actor.MessageEnvelope) {
			cnt++
			log.Printf("OutboundMiddleware: start sending message #%d to %s", cnt, target.Id)
			next(ctx, target, envelope)
			log.Printf("OutboundMiddleware: end sending message #%d to %s", cnt, target.Id)
		}
	}
}

type ping struct{}

type pong struct{}

func main() {
	pongProps := actor.
		FromFunc(func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *ping:
				ctx.Respond(&pong{})
			}
		})
	pongPid, _ := actor.SpawnNamed(pongProps, "pong")

	// Output should be somewhat like below.
	// Because ping actor receives both signal of struct{}{} and a pong message of &pong{},
	// the printed number of execution is doubled comparing to that of pong actor.
	//
	// 2019/08/25 12:04:19 InboundMiddleware: start handling incoming message #1: &actor.Started{}
	// 2019/08/25 12:04:19 Nested InboundMiddleware: start handling incoming message: &actor.Started{}
	// 2019/08/25 12:04:19 Actor: received &actor.Started{}
	// 2019/08/25 12:04:19 Nested InboundMiddleware: end handling incoming message: &actor.Started{}
	// 2019/08/25 12:04:19 InboundMiddleware: end handling incoming message #1: &actor.Started{}
	// 2019/08/25 12:04:20 InboundMiddleware: start handling incoming message # 2: struct {}{}
	// 2019/08/25 12:04:20 Nested InboundMiddleware: start handling incoming message: struct {}{}
	// 2019/08/25 12:04:20 Actor: received signal
	// 2019/08/25 12:04:20 OutboundMiddleware: start sending message #1 to pong
	// 2019/08/25 12:04:20 OutboundMiddleware: end sending message #1 to pong
	// 2019/08/25 12:04:20 Nested InboundMiddleware: end handling incoming message: struct {}{}
	// 2019/08/25 12:04:20 InboundMiddleware: end handling incoming message # 2: struct {}{}
	// 2019/08/25 12:04:20 InboundMiddleware: start handling incoming message #3: &main.pong{}
	// 2019/08/25 12:04:20 Nested InboundMiddleware: start handling incoming message: &main.pong{}
	// 2019/08/25 12:04:20 Actor: received pong
	// 2019/08/25 12:04:20 Nested InboundMiddleware: end handling incoming message: &main.pong{}
	// 2019/08/25 12:04:20 InboundMiddleware: end handling incoming message #3: &main.pong{}
	// 2019/08/25 12:04:21 InboundMiddleware: start handling incoming message #4: struct {}{}
	// 2019/08/25 12:04:21 Nested InboundMiddleware: start handling incoming message: struct {}{}
	// 2019/08/25 12:04:21 Actor: received signal
	// 2019/08/25 12:04:21 OutboundMiddleware: start sending message # 2 to pong
	// 2019/08/25 12:04:21 OutboundMiddleware: end sending message # 2 to pong
	// 2019/08/25 12:04:21 Nested InboundMiddleware: end handling incoming message: struct {}{}
	// 2019/08/25 12:04:21 InboundMiddleware: end handling incoming message #4: struct {}{}
	// 2019/08/25 12:04:21 InboundMiddleware: start handling incoming message #5: &main.pong{}
	// 2019/08/25 12:04:21 Nested InboundMiddleware: start handling incoming message: &main.pong{}
	// 2019/08/25 12:04:21 Actor: received pong
	// 2019/08/25 12:04:21 Nested InboundMiddleware: end handling incoming message: &main.pong{}
	// 2019/08/25 12:04:21 InboundMiddleware: end handling incoming message #5: &main.pong{}
	// 2019/08/25 12:04:22 InboundMiddleware: start handling incoming message #6: &actor.Stopping{}
	// 2019/08/25 12:04:22 Nested InboundMiddleware: start handling incoming message: &actor.Stopping{}
	// 2019/08/25 12:04:22 Actor: received &actor.Stopping{}
	// 2019/08/25 12:04:22 Nested InboundMiddleware: end handling incoming message: &actor.Stopping{}
	// 2019/08/25 12:04:22 InboundMiddleware: end handling incoming message #6: &actor.Stopping{}
	// 2019/08/25 12:04:22 InboundMiddleware: start handling incoming message #7: &actor.Stopped{}
	// 2019/08/25 12:04:22 Nested InboundMiddleware: start handling incoming message: &actor.Stopped{}
	// 2019/08/25 12:04:22 Actor: received &actor.Stopped{}
	// 2019/08/25 12:04:22 Nested InboundMiddleware: end handling incoming message: &actor.Stopped{}
	// 2019/08/25 12:04:22 InboundMiddleware: end handling incoming message #7: &actor.Stopped{}
	// 2019/08/25 12:04:22 Finish
	pingProps := actor.
		FromFunc(func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case struct{}:
				log.Print("Actor: received signal")
				ctx.Request(pongPid, &ping{})

			case *pong:
				log.Print("Actor: received pong")

			default:
				log.Printf("Actor: received %#v\n", ctx.Message())

			}
		}).
		WithMiddleware(newInboundMiddleware()).
		WithMiddleware(func(next actor.ActorFunc) actor.ActorFunc {
			return func(c actor.Context) {
				log.Printf("Nested InboundMiddleware: start handling incoming message: %#v", c.Message())
				next(c)
				log.Printf("Nested InboundMiddleware: end handling incoming message: %#v", c.Message())
			}
		}).
		WithOutboundMiddleware(newOutboundMiddleware())
	pingPid, _ := actor.SpawnNamed(pingProps, "ping")

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pingPid.Tell(struct{}{})

		case <-finish:
			pingPid.GracefulStop()
			pongPid.GracefulStop()
			log.Print("Finish")
			return

		}
	}
}
