package level

import (
	buff "github.com/alexykot/cncraft/pkg/buffer"
)

type Slice interface {
	buff.BPush

	Index() int

	Chunk() Chunk
	Level() Level

	// supports values x:[0:15] y:[0:15] z: [0:15]
	GetBlock(x, y, z int) Block
}
