package packet

import "fmt"

func prettyPrintBytesHex(bytes []byte) {
	const breakAt = 64
	var count int
	for _, byteVal := range bytes {
		fmt.Printf("%02X", byteVal)
		count++
		if count == breakAt {
			println()
			count = 0
		}
	}
	println()
}

func prettyPrintBytesBin(bytes []byte) {
	const breakAt = 8
	var count int
	for _, byteVal := range bytes {
		fmt.Printf("%08b", byteVal)
		count++
		if count == breakAt {
			println()
			count = 0
		}
	}
	println()
}
