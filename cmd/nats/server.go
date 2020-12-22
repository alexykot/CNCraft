// package main - prototype of NATS Streaming Server integration
// Tasks:
//  + basic publishing/subscribing
//  + multiple subscribers to a channel
//  - publishing from multiple publishers into single channel, interleaved
//  + manual ack of received messages
//  + channel publishing access control
//  + channel subscription access control
//  + channel creation access control
//  + subscriber restarting the flow
//  + durable subscriptions and restarting from arbitrary position
//  + partitioned publishing/subscribing
//
// Notes:
//  - STAN mTLS client connection authentication w/ statically provided CA cert.
//  - STAN durable queue group name must be same for all subscribers to work.
//  - STAN does not guarantee ordering, so every subscriber needs to do one of two things:
//      - if per-entity message order matters:
//          - every message should carry enough data to set the right state, i.e.
//            it should not rely on previous messages.
//          - latest processed message timestamp should be saved with the entity.
//          - out-of-order messages that are older than the latest recorded should be silently discarded
//      - if global order matters - rethink the architecture and avoid that, that will be a massive bottleneck.
//      - if all events for an entity matter:
//          - the order must not matter, ensure that in design
//          - the subscriber should keep own list of applied messages and deduplicate incoming
//      - if all events for an entity matter but a service does not have any persistence (i.e. no state at all):
//          - don't bother and just process everything
//
package main

import (
	"fmt"
	natsd "github.com/nats-io/nats-server/server"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

var nc *nats.Conn
var subj string
var interrupt chan os.Signal

func main() {

	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cmd := &cobra.Command{Use: "poc", Short: "nats prototype"}

	cmd.PersistentFlags().StringVar(&subj, "chan", "test-subj", "")

	registerNatsVerbs(cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("While executing command: %s\n", err)
	}
}

func startServer() {
	server := natsd.New(&natsd.Options{})
	go func() {
		defer func() {
			message := "nats stopped unexpectedly"
			if r := recover(); r != nil {
				message = fmt.Sprintf("nats panicked: %v", r)
			}
			println(message)
		}()
		server.ConfigureLogger()
		server.Start()
	}()

	if ok := server.ReadyForConnections(time.Second*3); !ok {
		panic("failed to start NATS server within the timeout")
	}
	println("started NATS server")
}

func startClient() {
	var err error
	if nc, err = nats.Connect(nats.DefaultURL); err != nil {
		panic(err)
	}
	println("started NATS client")
}

func registerNatsVerbs(cmd *cobra.Command) {
	natsCmd := &cobra.Command{Use: "nats", Short: "nats controls"}
	cmd.AddCommand(natsCmd)

	natsCmd.AddCommand(&cobra.Command{
		Use:   "all",
		Args:  cobra.NoArgs,
		Short: "start test NATS server, subscriber and publisher",
		RunE: func(cmd *cobra.Command, _ []string) error {

			startServer()
			startClient()

			sub, err := nc.Subscribe(subj, func(m *nats.Msg) {
				println(string(m.Data))
			})
			if err != nil {
				panic(err)
			}
			defer sub.Unsubscribe()
			println(fmt.Sprintf("subscribed to subject `%s`", subj))

			var messageNumber int64
			for {
				select {
					case <-interrupt:
						log.Println("all stopped")
						return nil
				default:
					<-time.After(1 * time.Second)
					message := strconv.Itoa(int(messageNumber)) + " - " + time.Now().Format(time.Stamp)
					atomic.AddInt64(&messageNumber, 1)
					if err := nc.Publish(subj, []byte(message)); err != nil {
						log.Println(fmt.Sprintf("failed to publish a message: %v", err))
					}
				}
			}
		},
	})
}
