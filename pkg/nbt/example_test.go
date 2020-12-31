package nbt

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/alexykot/cncraft/pkg/protocol/tags"
)

func ExampleUnmarshal() {
	var data = []byte{
		0x08, 0x00, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x00, 0x09,
		0x42, 0x61, 0x6e, 0x61, 0x6e, 0x72, 0x61, 0x6d, 0x61,
	}

	var Name string

	if err := Unmarshal(data, &Name); err != nil {
		panic(err)
	}

	fmt.Println(Name)

	// Output: Bananrama
}

func ExampleMarshal() {
	var value = struct {
		Name string `nbt:"name"`
	}{"Tnze"}

	var buf bytes.Buffer
	if err := Marshal(&buf, value); err != nil {
		panic(err)
	}

	fmt.Printf("% 02x ", buf.Bytes())

	// Output:
	//	0a 00 00 08 00 04 6e 61 6d 65 00 04 54 6e 7a 65 00
}

func TestUnmarshal(t *testing.T) {
	gzipBytes, err := ioutil.ReadFile("testdata/nbt/vanilla_world.nbt")
	if err != nil {
		panic(fmt.Errorf("failed to marshal: %w", err))
	}
	nbtBytes := inflateGZip(gzipBytes)

	codec := tags.DimensionCodec{}
	if err := Unmarshal(nbtBytes, &codec); err != nil {
		panic(err)
	}

	println(fmt.Sprintf("%#v\n", codec))
}

func TestMarshal(t *testing.T) {

	codec := tags.DimensionCodec{
		Dimensions: tags.DimensionRegistry{
			Type: "minecraft:dimension_type",
			RegistryEntries: []tags.DimensionRegistryEntry{
				{
					Name: "minecraft:overworld",
					ID:   0,
					Element: tags.Dimension{
						PiglinSafe:         0,
						Natural:            1,
						AmbientLight:       0.0,
						Infiniburn:         "minecraft:infiniburn_overworld",
						RespawnAnchorWorks: 0,
						HasSkylight:        1,
						BedWorks:           1,
						Effects:            "minecraft:overworld",
						HasRaids:           1,
						LogicalHeight:      256,
						CoordinateScale:    1.0,
						Ultrawarm:          0,
						HasCeiling:         0,
					},
				},
			},
		},
		Biomes: tags.BiomeRegistry{},
	}

	var nbtBuf bytes.Buffer
	if err := Marshal(&nbtBuf, codec); err != nil {
		panic(fmt.Errorf("failed to marshal: %w", err))
	}

	if err := ioutil.WriteFile("testdata/nbt/test_codec.nbt", compressGZip(nbtBuf.Bytes()), os.ModePerm); err != nil {
		panic(fmt.Errorf("failed to write file: %w", err))
	}

	fmt.Printf("%X\n", nbtBuf.Bytes())
}

func inflateGZip(gzipBuf []byte) []byte {
	reader, _ := gzip.NewReader(bytes.NewReader(gzipBuf))
	var out bytes.Buffer
	_, _ = io.Copy(&out, reader)
	_ = reader.Close()
	return out.Bytes()
}

func compressGZip(nbtBuf []byte) []byte {
	var buf bytes.Buffer
	writer, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	_, _ = writer.Write(nbtBuf)
	_ = writer.Close()
	return buf.Bytes()
}
