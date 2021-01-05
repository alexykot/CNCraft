package level

import (
	buff "github.com/alexykot/cncraft/pkg/buffer"
)

// 16*16*16 blocks cubic section, part of the chunk
type Section interface {
	buff.BPush

	// position in the chunk, 0 to 15
	Index() int

	// supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int) Block
}

type section struct {
}
