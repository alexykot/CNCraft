package packet

import "fmt"

func prettyPrintBytesHex(bytes []byte) {
	const breakAt = 64
	for i, byteVal := range bytes {
		fmt.Printf("%X", byteVal)
		if i == breakAt {
			println()
		}
	}
	println()
}

func prettyPrintBytesBin(bytes []byte) {
	const breakAt = 8
	for i, byteVal := range bytes {
		fmt.Printf("%08b", byteVal)
		if i == breakAt {
			println()
		}
	}
	println()
}
