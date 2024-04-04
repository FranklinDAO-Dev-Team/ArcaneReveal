package client

type GameStateStore map[string]*GameState

func NewGameStateStore() GameStateStore {
	return make(GameStateStore)
}

func (store GameStateStore) ReplaceGameState(id string, gameState *GameState) {
	store[id] = gameState
}

func (store GameStateStore) GetGameState(id string) (*GameState, bool) {
	gameState, exists := store[id]
	return gameState, exists
}

func (store GameStateStore) DeleteGameState(id string) bool {
	if _, exists := store[id]; !exists {
		return false
	}
	delete(store, id)
	return true
}
