package main

import (
	"context"
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"go.uber.org/zap"

	"github.com/alexykot/cncraft/core/db/orm"
	"github.com/alexykot/cncraft/pkg/game/items"
	pItems "github.com/alexykot/cncraft/pkg/protocol/items"

	"github.com/alexykot/cncraft/cmd/tools/packet"
	coreDB "github.com/alexykot/cncraft/core/db"
)

func main() {
	ctx := context.Background()
	signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)

	cmd := &cobra.Command{Use: "tools", Short: "misc tools"}
	packet.RegisterPacketTools(ctx, cmd)
	registerGenerationTools(ctx, cmd)
	registerMiscTools(ctx, cmd)

	if err := cmd.Execute(); err != nil {
		log.Fatalf("While executing command: %s\n", err)
	}
}

func registerGenerationTools(ctx context.Context, cmd *cobra.Command) {
	codegenCmd := &cobra.Command{Use: "gen {cmd}", Short: "codegen tools"}

	codegenCmd.AddCommand(&cobra.Command{
		Use:   "blocks {input_file.json}",
		Short: "block ids code generator",
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

	codegenCmd.AddCommand(&cobra.Command{
		Use:   "items {input_file.json}",
		Short: "items ids code generator",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			type item struct {
				ID int `json:"protocol_id"`
			}
			type itemList map[string]item

			registries := struct {
				Items struct {
					Entries itemList `json:"entries"`
				} `json:"minecraft:item"`
			}{}

			var outFile *os.File
			var err error

			inputFileName := args[0]
			outputFileName := args[1]

			if outputFileName == "-" {
				outFile = os.Stdout
			} else if outFile, err = os.Create(outputFileName); err != nil {
				return fmt.Errorf("failed to open output file %s: %w", outputFileName, err)
			}
			defer outFile.Close()

			input, err := ioutil.ReadFile(inputFileName)
			if err != nil {
				return fmt.Errorf("failed to read source file %s: %w", inputFileName, err)
			}

			if err := json.Unmarshal(input, &registries); err != nil {
				return fmt.Errorf("failed to open source file %s: %w", inputFileName, err)
			}

			itemsMap := registries.Items.Entries

			itemsList := make([]string, len(itemsMap), len(itemsMap))

			var constBlob, mapBlob string
			for itemName, itemData := range itemsMap {
				itemsList[itemData.ID] = itemName
			}
			for ID, itemName := range itemsList {
				constBlob = constBlob + fmt.Sprintf("%s ItemID = %d\n", getConstName(itemName, nil), ID)
				mapBlob = mapBlob + fmt.Sprintf("namesMap[%d] = \"%s\"\n", ID, itemName)
			}

			goResult := fmt.Sprintf(`package items

type ItemID uint32
func(b ItemID) String() string{return namesMap[b]}
func(b ItemID) ID() uint32{return uint32(b)}

const (
%s)

var namesMap map[ItemID]string
func init(){
	namesMap = make(map[ItemID]string)
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

func getConstName(registryName string, props map[string]string) string {
	registryName = strings.Replace(registryName, "minecraft:", "", 1)
	registryName = strings.Replace(registryName, "_", " ", -1)
	registryName = strings.Title(registryName)
	registryName = strings.Replace(registryName, " ", "", -1)

	if len(props) > 0 {
		registryName = registryName + "_"
	}

	for prop, value := range props {
		registryName = registryName + strings.Title(prop) + strings.Title(value)
	}

	return registryName
}

func registerMiscTools(ctx context.Context, cmd *cobra.Command) {
	miscCmd := &cobra.Command{Use: "misc {cmd}", Short: "misc tools"}
	miscCmd.AddCommand(&cobra.Command{
		Use:   "readfile {file}",
		Short: "output binary contents of the fix in hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bytes, err := ioutil.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}

			println(fmt.Sprintf("contents of: %s, total byte length: %d", args[0], len(bytes)))

			var i int
			for i = 64; i < len(bytes); i = i + 64 {
				println(fmt.Sprintf("%X", bytes[i-64:i]))
			}

			if i-64 < len(bytes) {
				println(fmt.Sprintf("%X", bytes[i-64:len(bytes)-1]))
			}

			return nil
		},
	})

	miscCmd.AddCommand(&cobra.Command{
		Use:   "idkfa {player_name}",
		Short: "give player all weapons",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbURL, ok := os.LookupEnv("CNCRAFT_TEST_DB_URL")
			if !ok {
				return fmt.Errorf("CNCRAFT_TEST_DB_URL envar must be set to a valid DB URL")
			}

			db, err := coreDB.New(dbURL, false, zap.L())
			if err != nil {
				return fmt.Errorf("failed to open DB URL %s: %w", dbURL, err)
			}

			dbPlayer, err := orm.Players(orm.PlayerWhere.Username.EQ(args[0])).One(ctx, db)
			if err != nil {
				return fmt.Errorf("failed to query player: %w", err)
			}

			inventory := items.NewInventory()
			inventory.RowHotbar = [9]items.Slot{
				{
					IsPresent: true,
					ItemID:    int16(pItems.DiamondPickaxe),
					ItemCount: 1,
				},
				{
					IsPresent: true,
					ItemID:    int16(pItems.Bedrock),
					ItemCount: 40,
				},
			}

			for slotNum, slot := range inventory.ToArray() {
				if slot.IsPresent {
					dbItem := orm.Inventory{
						PlayerID:   dbPlayer.ID,
						SlotNumber: int16(slotNum),
						ItemID:     slot.ItemID,
						ItemCount:  slot.ItemCount,
					}
					if err := dbItem.Upsert(ctx, db, true,
						[]string{orm.InventoryColumns.PlayerID, orm.InventoryColumns.SlotNumber},
						boil.Whitelist(orm.InventoryColumns.ItemID, orm.InventoryColumns.ItemCount),
						boil.Infer()); err != nil {
						return fmt.Errorf("failed to update player inventory slot: %w", err)
					}
				}
			}

			return nil
		},
	})

	cmd.AddCommand(miscCmd)
}
