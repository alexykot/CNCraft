package data

type PositionI struct {
	X int64
	Y int64
	Z int64
}

type PositionF struct {
	X float64
	Y float64
	Z float64
}

type RotationF struct {
	AxisX float32 // 'yaw'
	AxisY float32 // 'pitch'
}

type Location struct {
	PositionF
	RotationF
}
