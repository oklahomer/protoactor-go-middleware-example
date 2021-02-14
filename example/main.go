package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
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
	// Setup actor system
	system := actor.NewActorSystem()

	// Run a pong actor that receives ping payload and send back pong payload.
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
	// Because ping actor receives both signal of struct{}{} and a pong message of &pong{},
	// the printed number of execution is doubled comparing to that of pong actor.
	//
	// 2019/08/25 13:10:21 ReceiverMiddleware: start handling incoming message #1: &actor.Started{}
	// 2019/08/25 13:10:21 Nested ReceiverMiddleware: start handling incoming message: &actor.Started{}
	// 2019/08/25 13:10:21 Actor: received &actor.Started{}
	// 2019/08/25 13:10:21 Nested ReceiverMiddleware: end handling incoming message: &actor.Started{}
	// 2019/08/25 13:10:21 ReceiverMiddleware: end handling incoming message #1: &actor.Started{}
	// 2019/08/25 13:10:22 ReceiverMiddleware: start handling incoming message #2: struct {}{}
	// 2019/08/25 13:10:22 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2019/08/25 13:10:22 Actor: received signal
	// 2019/08/25 13:10:22 SenderMiddleware: start sending message #1 to pong
	// 2019/08/25 13:10:22 SenderMiddleware: end sending message #1 to pong
	// 2019/08/25 13:10:22 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2019/08/25 13:10:22 ReceiverMiddleware: end handling incoming message #2: struct {}{}
	// 2019/08/25 13:10:22 ReceiverMiddleware: start handling incoming message #3: &main.pong{}
	// 2019/08/25 13:10:22 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2019/08/25 13:10:22 Actor: received pong
	// 2019/08/25 13:10:22 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2019/08/25 13:10:22 ReceiverMiddleware: end handling incoming message #3: &main.pong{}
	// 2019/08/25 13:10:23 ReceiverMiddleware: start handling incoming message #4: struct {}{}
	// 2019/08/25 13:10:23 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2019/08/25 13:10:23 Actor: received signal
	// 2019/08/25 13:10:23 SenderMiddleware: start sending message #2 to pong
	// 2019/08/25 13:10:23 SenderMiddleware: end sending message #2 to pong
	// 2019/08/25 13:10:23 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2019/08/25 13:10:23 ReceiverMiddleware: end handling incoming message #4: struct {}{}
	// 2019/08/25 13:10:23 ReceiverMiddleware: start handling incoming message #5: &main.pong{}
	// 2019/08/25 13:10:23 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2019/08/25 13:10:23 Actor: received pong
	// 2019/08/25 13:10:23 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2019/08/25 13:10:23 ReceiverMiddleware: end handling incoming message #5: &main.pong{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: start handling incoming message #6: struct {}{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: start handling incoming message: struct {}{}
	// 2019/08/25 13:10:24 Actor: received signal
	// 2019/08/25 13:10:24 SenderMiddleware: start sending message #3 to pong
	// 2019/08/25 13:10:24 SenderMiddleware: end sending message #3 to pong
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: end handling incoming message: struct {}{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: end handling incoming message #6: struct {}{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: start handling incoming message #7: &main.pong{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: start handling incoming message: &main.pong{}
	// 2019/08/25 13:10:24 Actor: received pong
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: end handling incoming message: &main.pong{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: end handling incoming message #7: &main.pong{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: start handling incoming message #8: &actor.Stopping{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopping{}
	// 2019/08/25 13:10:24 Actor: received &actor.Stopping{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopping{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: end handling incoming message #8: &actor.Stopping{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: start handling incoming message #9: &actor.Stopped{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: start handling incoming message: &actor.Stopped{}
	// 2019/08/25 13:10:24 Actor: received &actor.Stopped{}
	// 2019/08/25 13:10:24 Nested ReceiverMiddleware: end handling incoming message: &actor.Stopped{}
	// 2019/08/25 13:10:24 ReceiverMiddleware: end handling incoming message #9: &actor.Stopped{}
	// 2019/08/25 13:10:24 Finish
	pingProps := actor.
		PropsFromFunc(func(ctx actor.Context) {
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
		WithReceiverMiddleware(newReceiverMiddleware()).
		WithReceiverMiddleware(func(next actor.ReceiverFunc) actor.ReceiverFunc {
			return func(c actor.ReceiverContext, env *actor.MessageEnvelope) {
				_, msg, _ := actor.UnwrapEnvelope(env)
				log.Printf("Nested ReceiverMiddleware: start handling incoming message: %#v", msg)
				next(c, env)
				log.Printf("Nested ReceiverMiddleware: end handling incoming message: %#v", msg)
			}
		}).
		WithSenderMiddleware(newSenderMiddleware())

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
