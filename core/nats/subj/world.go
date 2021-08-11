package subj

// *** World related subjects ***

// MkShardEvent creates a subject name string for shard events, i.e. all events that are processed
// by the event loop of the given shard.
func MkShardEvent(shardID string) Subj {
	return Subj("world.event." + shardID)
}
