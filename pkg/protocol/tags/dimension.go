package tags

type DimensionCodec struct {
	Dimensions DimensionRegistry `nbt:"minecraft:dimension_type"` // TAG_Compound 	The dimension type registry (see below).
	Biomes     BiomeRegistry     `nbt:"minecraft:worldgen/biome"` // TAG_Compound 	The biome registry (see below).
}

type DimensionRegistry struct {
	Type            string                   `nbt:"type"`  // TAG_String 	The name of the registry. Always "minecraft:dimension_type".
	RegistryEntries []DimensionRegistryEntry `nbt:"value"` // TAG_List 	List of dimension types registry entries (see below).
}

type DimensionRegistryEntry struct {
	Name    string    `nbt:"name"`    // TAG_String 	The name of the dimension type (for example, "minecraft:overworld").
	ID      int32     `nbt:"id"`      // TAG_Int 	The protocol ID of the dimension (matches the index of the element in the registry list).
	Element Dimension `nbt:"element"` // TAG_Compound 	The dimension type (see below).
}

type Dimension struct {
	PiglinSafe         uint8   `nbt:"piglin_safe"`          // TAG_Byte      Whether piglins shake and transform to zombified piglins. 	1: true, 0: false.
	Natural            uint8   `nbt:"natural"`              // TAG_Byte      When false, compasses spin randomly. When true, nether portals can spawn zombified piglins. 	1: true, 0: false.
	AmbientLight       float32 `nbt:"ambient_light"`        // TAG_Float     How much light the dimension has. 	0.0 to 1.0.
	Infiniburn         string  `nbt:"infiniburn"`           // TAG_String    A resource location defining what block tag to use for infiniburn. 	"" or minecraft resource "minecraft:...".
	RespawnAnchorWorks uint8   `nbt:"respawn_anchor_works"` // TAG_Byte      Whether players can charge and use respawn anchors. 	1: true, 0: false.
	HasSkylight        uint8   `nbt:"has_skylight"`         // TAG_Byte      Whether the dimension has skylight access or not. 	1: true, 0: false.
	BedWorks           uint8   `nbt:"bed_works"`            // TAG_Byte      Whether players can use a bed to sleep. 	1: true, 0: false.
	Effects            string  `nbt:"effects"`              // TAG_String    ? 	"minecraft:overworld", "minecraft:the_nether", "minecraft:the_end" or something else.
	HasRaids           uint8   `nbt:"has_raids"`            // TAG_Byte      Whether players with the Bad Omen effect can cause a raid. 	1: true, 0: false.
	LogicalHeight      int32   `nbt:"logical_height"`       // TAG_Int       The maximum height to which chorus fruits and nether portals can bring players within this dimension. 	0-256.
	CoordinateScale    float64 `nbt:"coordinate_scale"`     // TAG_Float     The multiplier applied to coordinates when traveling to the dimension. 	1: true, 0: false.
	Ultrawarm          uint8   `nbt:"ultrawarm"`            // TAG_Byte      Whether the dimensions behaves like the nether (water evaporates and sponges dry) or not. Also causes lava to spread thinner. 	1: true, 0: false.
	HasCeiling         uint8   `nbt:"has_ceiling"`          // TAG_Byte      Whether the dimension has a bedrock ceiling or not. When true, causes lava to spread faster. 	1: true, 0: false.
	FixedTime          int64   `nbt:"fixed_time"`           // TAG_Long      Optional, if set, the time of the day is the specified value. 	If set, 0 to 24000.
}
