package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/alexykot/cncraft/pkg/buffer"
)

var interrupt chan os.Signal

func main() {

	interrupt = make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	cmd := &cobra.Command{Use: "tools", Short: "misc tools"}

	registerPacketTools(cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("While executing command: %s\n", err)
	}
}

func registerPacketTools(cmd *cobra.Command) {
	packetCmd := &cobra.Command{Use: "packet {cmd}", Short: "packet tools"}

	decodeCmd := &cobra.Command{Use: "decode {cmd}", Short: "byte value decoding tools"}
	decodeCmd.AddCommand(&cobra.Command{
		Use:   "varint {hex value}",
		Args:  cobra.ExactArgs(1),
		Short: "decode hexed bytes that should represent a varint into a decimal value",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}

			bufTest := buffer.NewFrom(hexBytes)
			integer := bufTest.PullVarInt()
			println(fmt.Sprintf("integer: %d", integer))

			return nil
		},
	})

	cmd.AddCommand(packetCmd)
	packetCmd.AddCommand(decodeCmd)
}
