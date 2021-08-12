package subj

// *** Player related subjects ***

// MkPlayerLoading creates a subject name string for announcing new users joining server.
//  This is send after successful login and triggers client world loading and player spawn.
func MkPlayerLoading() Subj { return "players.loading" }

// MkPlayerJoined creates a subject name string for announcing players joining server.
//  This is sent after the player has successfully spawned in the world.
func MkPlayerJoined() Subj { return "players.joined.first_time" }

// MkPlayerLeft creates a subject name string for announcing players leaving server.
func MkPlayerLeft() Subj { return "players.left" }

// MkPlayerSpatialUpdate creates a subject name string for announcing player position updates.
//  This is sent every time player position changes.
func MkPlayerSpatialUpdate() Subj { return "players.update.spatial" }

// MkPlayerInventoryUpdate creates a subject name string for announcing player inventory updates.
//  This is sent every time player inventory changes (including hotbar).
func MkPlayerInventoryUpdate() Subj { return "players.update.inventory" }
