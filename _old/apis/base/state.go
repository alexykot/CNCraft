package base

type Loads interface {
	Load()
}

type Kills interface {
	Stop()
}

type State interface {
	Loads
	Kills
}
