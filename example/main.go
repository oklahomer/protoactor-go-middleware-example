package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"log"
	"os"
	"os/signal"
	"time"
)

func newReceiverMiddleware() actor.ReceiverMiddleware {
	cnt := 0
	return func(next actor.ReceiverFunc) actor.ReceiverFunc {
		return func(ctx actor.ReceiverContext, env *actor.MessageEnvelope) {
			_, msg, _ := actor.UnwrapEnvelope(env)
			cnt++
			log.Printf("ReceiverMiddleware: start handling incoming message #%d: %#v", cnt, msg)
			next(ctx, env)
			log.Printf("ReceiverMiddleware: end handling incoming message #%d: %#v", cnt, msg)
		}
	}
}

func newSenderMiddleware() actor.SenderMiddleware {
	cnt := 0
	return func(next actor.SenderFunc) actor.SenderFunc {
		return func(ctx actor.SenderContext, target *actor.PID, env *actor.MessageEnvelope) {
			cnt++
			log.Printf("SenderMiddleware: start sending message #%d to %s", cnt, target.Id)
			next(ctx, target, env)
			log.Printf("SenderMiddleware: end sending message #%d to %s", cnt, target.Id)
		}
	}
}

type ping struct{}

type pong struct{}

func main() {
	// Set up the actor system
	system := actor.NewActorSystem()

	// Run a pong actor that receives a ping payload and send back a pong payload.
	pongProps := actor.
		PropsFromFunc(func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *ping:
				ctx.Respond(&pong{})
			}
		})
	pongPid, _ := system.Root.SpawnNamed(pongProps, "pong")

	// Run a ping actor with nested receiver middlewares and one sender middleware.
	//
	// Output should be somewhat like below.
	// Because the ping actor receives both signal of struct{}{} and a pong message of &pong{},
	// the printed number of executions is doubled comparing to that of a pong actor.
	//
	// 2022/07/31 12:54:55 ReceiverMiddleware: start handling incoming message #1: &actor.Started{}
	// 2022/07/31 12:54:55 Nested ReceiverMiddleware: start handling incoming message: &actor.Started{}
	// 2022/07/31 12:54:55 Actor: received &actor.Started{}
	// 2022/07/31 12:54:55 Nested ReceiverMiddleware: end handling incoming message: &actor.Started{}
	// 2022/07/31 12:54:55 ReceiverMiddleware: end handling incoming message #1: &actor.Started{}
	// 2022/07/31 12:54:56 ReceiverMiddleware: start handling incoming message #2: struct {}{}
	// 2022/07/31 12:54:56 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2022/07/31 12:54:56 Actor: received signal
	// 2022/07/31 12:54:56 SenderMiddleware: start sending message #1 to pong
	// 2022/07/31 12:54:56 SenderMiddleware: end sending message #1 to pong
	// 2022/07/31 12:54:56 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2022/07/31 12:54:56 ReceiverMiddleware: end handling incoming message #2: struct {}{}
	// 2022/07/31 12:54:56 ReceiverMiddleware: start handling incoming message #3: &main.pong{}
	// 2022/07/31 12:54:56 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2022/07/31 12:54:56 Actor: received pong
	// 2022/07/31 12:54:56 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2022/07/31 12:54:56 ReceiverMiddleware: end handling incoming message #3: &main.pong{}
	// 2022/07/31 12:54:57 ReceiverMiddleware: start handling incoming message #4: struct {}{}
	// 2022/07/31 12:54:57 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2022/07/31 12:54:57 Actor: received signal
	// 2022/07/31 12:54:57 SenderMiddleware: start sending message #2 to pong
	// 2022/07/31 12:54:57 SenderMiddleware: end sending message #2 to pong
	// 2022/07/31 12:54:57 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2022/07/31 12:54:57 ReceiverMiddleware: end handling incoming message #4: struct {}{}
	// 2022/07/31 12:54:57 ReceiverMiddleware: start handling incoming message #5: &main.pong{}
	// 2022/07/31 12:54:57 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2022/07/31 12:54:57 Actor: received pong
	// 2022/07/31 12:54:57 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2022/07/31 12:54:57 ReceiverMiddleware: end handling incoming message #5: &main.pong{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #6: struct {}{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2022/07/31 12:54:58 Actor: received signal
	// 2022/07/31 12:54:58 SenderMiddleware: start sending message #3 to pong
	// 2022/07/31 12:54:58 SenderMiddleware: end sending message #3 to pong
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #6: struct {}{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #7: &main.pong{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2022/07/31 12:54:58 Actor: received pong
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #7: &main.pong{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #8: &actor.Stopping{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopping{}
	// 2022/07/31 12:54:58 Actor: received &actor.Stopping{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopping{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #8: &actor.Stopping{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: start handling incoming message #9: &actor.Stopped{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopped{}
	// 2022/07/31 12:54:58 Actor: received &actor.Stopped{}
	// 2022/07/31 12:54:58 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopped{}
	// 2022/07/31 12:54:58 ReceiverMiddleware: end handling incoming message #9: &actor.Stopped{}
	// 2022/07/31 12:54:58 Finish
	pingProps := actor.PropsFromFunc(
		func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case struct{}:
				log.Print("Actor: received signal")
				ctx.Request(pongPid, &ping{})

			case *pong:
				log.Print("Actor: received pong")

			default:
				log.Printf("Actor: received %#v\n", ctx.Message())

			}
		},
		actor.WithReceiverMiddleware(newReceiverMiddleware()),
		actor.WithReceiverMiddleware(func(next actor.ReceiverFunc) actor.ReceiverFunc {
			return func(c actor.ReceiverContext, env *actor.MessageEnvelope) {
				_, msg, _ := actor.UnwrapEnvelope(env)
				log.Printf("Nested ReceiverMiddleware: start handling incoming message: %#v", msg)
				next(c, env)
				log.Printf("Nested ReceiverMiddleware: end handling incoming message: %#v", msg)
			}
		}),
		actor.WithSenderMiddleware(newSenderMiddleware()),
	)

	pingPid, _ := system.Root.SpawnNamed(pingProps, "ping")

	// Subscribe to signal to finish interaction
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	// Periodically send ping payload till signal comes
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})

		case <-finish:
			system.Root.StopFuture(pingPid).Wait()
			system.Root.StopFuture(pongPid).Wait()
			log.Print("Finish")
			return

		}
	}
}
