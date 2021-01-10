package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
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
	registerGenerationTools(cmd)

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

func registerGenerationTools(cmd *cobra.Command) {
	codegenCmd := &cobra.Command{Use: "gen {cmd}", Short: "codegen tools"}

	codegenCmd.AddCommand(&cobra.Command{
		Use:   "blocks {input_file.json}",
		Short: "block ids code generation",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			type blockState struct {
				ID        int               `json:"id"`
				IsDefault bool              `json:"default"`
				Props     map[string]string `json:"properties"`
			}

			type block struct {
				States []blockState        `json:"states"`
				Props  map[string][]string `json:"properties"`
			}

			type blockList map[string]block
			var outFile *os.File
			var err error

			inputFileName := args[0]
			outputFileName := args[1]

			if outputFileName == "-" {
				outFile = os.Stdout
			} else if outFile, err = os.Open(outputFileName); err != nil {
				return fmt.Errorf("failed to open output file %s: %w", outputFileName, err)
			}
			defer outFile.Close()

			input, err := ioutil.ReadFile(inputFileName)
			if err != nil {
				return fmt.Errorf("failed to read source file %s: %w", inputFileName, err)
			}

			theList := make(blockList)
			if err := json.Unmarshal(input, &theList); err != nil {
				return fmt.Errorf("failed to open source file %s: %w", inputFileName, err)
			}

			var constBlob string
			var mapBlob string
			for blockName, blockData := range theList {
				for _, state := range blockData.States {
					constBlob = constBlob + fmt.Sprintf("%s BlockID = %d\n", getConstName(blockName, state.Props), state.ID)
					mapBlob = mapBlob + fmt.Sprintf("namesMap[%d] = \"%s\"\n", state.ID, blockName)
				}
			}

			goResult := fmt.Sprintf(`package blocks

type BlockID uint32
func(b BlockID) String() string{return namesMap[b]}
func(b BlockID) ID() uint32{return uint32(b)}

const (
%s)

var namesMap map[BlockID]string
func init(){
	namesMap = make(map[BlockID]string)
%s
}
`, constBlob, mapBlob)

			result, err := format.Source([]byte(goResult))
			if err != nil {
				return fmt.Errorf("failed to format the output: %w", err)
			}

			if _, err := outFile.Write(result); err != nil {
				return fmt.Errorf("failed to write output: %w", err)
			}

			return nil
		},
	})

	cmd.AddCommand(codegenCmd)
}

func getConstName(blockName string, props map[string]string) string {
	blockName = strings.Replace(blockName, "minecraft:", "", 1)
	blockName = strings.Replace(blockName, "_", " ", -1)
	blockName = strings.Title(blockName)
	blockName = strings.Replace(blockName, " ", "", -1)

	if len(props) > 0 {
		blockName = blockName + "_"
	}

	for prop, value := range props {
		blockName = blockName + strings.Title(prop) + strings.Title(value)
	}

	return blockName
}
