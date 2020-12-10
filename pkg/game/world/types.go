package world

type WorldType int

const (
	DEFAULT WorldType = iota
	FLAT
	LARGEBIOMES
	AMPLIFIED
	CUSTOMIZED
	BUFFET
	DEFAULT11
)

var typeToName = map[WorldType]string{
	DEFAULT:     "default",
	FLAT:        "flat",
	LARGEBIOMES: "largeBiomes",
	AMPLIFIED:   "amplified",
	CUSTOMIZED:  "customized",
	BUFFET:      "buffet",
	DEFAULT11:   "default_1_1",
}

func (l WorldType) String() string {
	return typeToName[l]
}
