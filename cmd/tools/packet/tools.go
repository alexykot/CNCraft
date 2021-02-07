package packet

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alexykot/cncraft/pkg/buffer"
)

func RegisterPacketTools(cmd *cobra.Command) {
	packetCmd := &cobra.Command{Use: "packet {cmd}", Short: "packet tools"}

	decodeCmd := &cobra.Command{Use: "decode {cmd}", Short: "byte value decoding tools"}
	decodeCmd.AddCommand(&cobra.Command{
		Use: "varint {hex value}", Args: cobra.ExactArgs(1),
		Short: "decode hexed bytes that should represent a varint into a decimal value",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			println(fmt.Sprintf("integer: %d", buffer.NewFrom(hexBytes).PullVarInt()))
			return nil
		},
	})

	decodeCmd.AddCommand(&cobra.Command{
		Use: "hexdec {hex value}", Args: cobra.ExactArgs(1),
		Short: "decode hexed bytes into uint64 decimal value, 8 bytes max",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(fmt.Sprintf("%016s", args[0]))
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			println(fmt.Sprintf("integer: %d", buffer.NewFrom(hexBytes).PullUint64()))
			return nil
		},
	})

	decodeCmd.AddCommand(&cobra.Command{
		Use: "hexformat {hex value}", Args: cobra.ExactArgs(1),
		Short: "format single hex string into 64 bytes strings",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			prettyPrintBytesHex(hexBytes)
			return nil
		},
	})

	decodeCmd.AddCommand(&cobra.Command{
		Use: "hexbin {hex value}", Args: cobra.ExactArgs(1),
		Short: "decode hexed bytes into it's binary representation",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			prettyPrintBytesBin(hexBytes)
			return nil
		},
	})

	decodeCmd.AddCommand(&cobra.Command{
		Use: "hexutf {hex value}", Args: cobra.ExactArgs(1),
		Short: "decode hexed bytes that should represent an UTF8 string into the string",
		RunE: func(cmd *cobra.Command, args []string) error {
			hexBytes, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("failed to decode hex string: %w", err)
			}
			println(fmt.Sprintf("string: %s", string(hexBytes)))
			return nil
		},
	})

	cmd.AddCommand(packetCmd)
	packetCmd.AddCommand(decodeCmd)
}
