package action_type

const (
	PlayerSpawnId = iota
	PlayerDespawnId
	PlayerMoveId
	OrientationChangeId
	AskChunkActionId
	SetChunkActionId
)

type PlayerInitData struct {
	Name string
	X, Y, Orientation int
	Skin int
}

type InitWorldState struct {
	// for now, all peer need to know on connection is 
	// player position
	Players []PlayerInitData
}