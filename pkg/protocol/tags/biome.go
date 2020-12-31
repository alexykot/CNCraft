package tags

type BiomeRegistry struct {
	Type            string               `nbt:"type"`  // TAG_String 	The name of the registry. Always "minecraft:worldgen/biome".
	RegistryEntries []BiomeRegistryEntry `nbt:"value"` // TAG_List 	List of biome registry entries (see below).
}

type BiomeRegistryEntry struct {
	Name    string `nbt:"name"`    // 	TAG_String 	The name of the biome (for example, "minecraft:ocean").
	ID      int32  `nbt:"id"`      // 	TAG_Int 	The protocol ID of the biome (matches the index of the element in the registry list).
	Element Biome  `nbt:"element"` // 	TAG_Compound 	The biome properties (see below).
}

type Biome struct {
	Precipitation string       `nbt:"precipitation"` // TAG_String 	    The type of precipitation in the biome. 	"rain", "snow", or "none".
	Depth         float32      `nbt:"depth"`         // TAG_Float 	    The depth factor of the biome. 	The default values vary between 1.5 and -1.8.
	Temperature   float32      `nbt:"temperature"`   // TAG_Float 	    The temperature factor of the biome. 	The default values vary between 2.0 and -0.5.
	Scale         float32      `nbt:"scale"`         // TAG_Float 	    ? 	The default values vary between 1.225 and 0.0.
	Downfall      float32      `nbt:"downfall"`      // TAG_Float 	    ? 	The default values vary between 1.0 and 0.0.
	Category      string       `nbt:"category"`      // TAG_String 	    The category of the biome. 	Known values are "ocean", "plains", "desert", "forest", "extreme_hills", "taiga", "swamp", "river", "nether", "the_end", "icy", "mushroom", "beach", "jungle", "mesa", "savanna", and "none".
	Effects       BiomeEffects `nbt:"effects"`       // TAG_Compound 	Biome effects, see below.
}

type BiomeEffects struct {
	SkyColor           int32                      `nbt:"sky_color"`            // TAG_Int 	The color of the sky. 	Example: 8364543, which is #7FA1FF in RGB.
	WaterFogColor      int32                      `nbt:"water_fog_color"`      // TAG_Int 	Possibly the tint color when swimming. 	Example: 8364543, which is #7FA1FF in RGB.
	FogColor           int32                      `nbt:"fog_color"`            // TAG_Int 	Possibly the color of the fog effect when looking past the view distance. 	Example: 8364543, which is #7FA1FF in RGB.
	WaterColor         int32                      `nbt:"water_color"`          // TAG_Int 	The tint color of the water blocks. 	Example: 8364543, which is #7FA1FF in RGB.
	FoliageColor       int32                      `nbt:"foliage_color"`        // TAG_Int    Optional. The tint color of the grass. 	Example: 8364543, which is #7FA1FF in RGB.
	GrassColorModifier string                     `nbt:"grass_color_modifier"` // TAG_String    Optional. Unknown, likely affects foliage color. 	If set, known values are "swamp" and "dark_forest".
	AmbientSound       string                     `nbt:"ambient_sound"`        // TAG_String    Optional. Ambient soundtrack. 	If present, the ID of a soundtrack. Example: "minecraft:ambient.basalt_deltas.loop".
	Music              BiomeEffectsMusic          `nbt:"music"`                // TAG_Compound    Optional. Music properties for the biome. 	If present, contains the fields: replace_current_music (TAG_Byte), sound (TAG_String), max_delay (TAG_Int), min_delay (TAG_Int).
	AdditionsSound     BiomeEffectsAdditionsSound `nbt:"additions_sound"`      // TAG_Compound    Optional. Additional ambient sound that plays randomly. 	If present, contains the fields: sound (TAG_String), tick_chance (TAG_Double).
	MoodSound          BiomeEffectsMoodSound      `nbt:"mood_sound"`           // TAG_Compound    Optional. Additional ambient sound that plays at an interval. 	If present, contains the fields: sound (TAG_String), tick_delay (TAG_Int), offset (TAG_Double), block_search_extend (TAG_Int).
	Particle           BiomeEffectsParticle       `nbt:"particle"`             // TAG_Compound    Optional. Particles that appear randomly in the biome. 	If present, contains the fields: probability (TAG_Float), options (TAG_Compound). The "options" compound contains the field "type" (TAG_String), which identifies the particle type.
}
type BiomeEffectsMusic struct {
	ReplaceCurrentMusic uint8  `nbt:"replace_current_music"` // TAG_Byte
	Sound               string `nbt:"sound"`                 // TAG_String
	MaxDelay            int32  `nbt:"max_delay"`             // TAG_Int
	MinDelay            int32  `nbt:"min_delay"`             // TAG_Int
}
type BiomeEffectsAdditionsSound struct {
	Sound      string  `nbt:"sound"`       // TAG_String
	TickChance float64 `nbt:"tick_chance"` // TAG_Double
}

type BiomeEffectsMoodSound struct {
	Sound             string  `nbt:"sound"`               //  TAG_String
	TickDelay         int32   `nbt:"tick_delay"`          //  TAG_Int
	Offset            float64 `nbt:"offset"`              //  TAG_Double
	BlockSearchExtend int32   `nbt:"block_search_extend"` //  TAG_Int
}
type BiomeEffectsParticle struct {
	Probability float32                     `nbt:"probability"` // TAG_Float
	Options     BiomeEffectsParticleOptions `nbt:"options"`     // TAG_Compound
}

type BiomeEffectsParticleOptions struct {
	ParticleType string `nbt:"type"` // TAG_String
}
