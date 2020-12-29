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

	registerTools(cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("While executing command: %s\n", err)
	}
}

func registerTools(cmd *cobra.Command) {
	cmd.AddCommand(&cobra.Command{
		Use:   "varint",
		Args:  cobra.ExactArgs(1),
		Short: "varint to decimal",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}

			bufTest := buffer.NewFrom(hexBytes)
			integer := bufTest.PullVrI()
			println(fmt.Sprintf("integer: %d", integer))

			return nil
		},
	})
}
