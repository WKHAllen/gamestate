package gamestate

import (
	"fmt"
)

// GameStateID is an unsigned integer that uniquely identifies a game state.
type GameStateID uint32

// GameState represents a single state a game can exist in. The execution of a
// game state involves calling a user-defined function.
type GameState struct {
	id GameStateID
	action func() any
	transitions map[any]*GameState
	defaultTransition *GameState
}

// GameMap represents a mapping of game states and values they may access.
type GameMap struct {
	nextStateID GameStateID
	states map[GameStateID]*GameState
	values map[string]any
}

// gameMap is a global instance of a GameMap in which all game states and
// values are stored.
var gameMap *GameMap = &GameMap{
	nextStateID: 0,
	states: make(map[GameStateID]*GameState),
	values: make(map[string]any),
}

// NewState creates a new state in the game map. `action` is the function that
// will be called when the state becomes active.
func NewGameState(action func() any) *GameState {
	newState := &GameState{
		id: gameMap.nextStateID,
		action: action,
		transitions: make(map[any]*GameState),
		defaultTransition: nil,
	}

	gameMap.states[gameMap.nextStateID] = newState
	gameMap.nextStateID++

	return newState
}

// Transition defines a transition from one state to another. This is used as
// the fallback transition if no transitions based on return values match.
func (gameState *GameState) Transition(toState *GameState) {
	gameState.defaultTransition = toState
}

// TransitionOn defines a transition from one state to another on the
// condition that the return value of the current state's action function
// matches the `on` parameter. This enables the possibility of conditionally
// branching between multiple distinct states.
func (gameState *GameState) TransitionOn(toState *GameState, on any) {
	gameState.transitions[on] = toState
}

// nextState attempts to find the next state to transition to based on the
// return value of the current state's action function.
func (gameState *GameState) nextState(on any) (*GameState, error) {
	if nextState, ok := gameState.transitions[on]; ok {
		return nextState, nil
	} else if gameState.defaultTransition != nil {
		return gameState.defaultTransition, nil
	} else {
		return nil, fmt.Errorf("no relevant state or default state found")
	}
}

// GetValue retrieves a value from the game map's value mapping. The first
// returned value is the value itself, and the second returned value is a
// boolean representing whether the value exists in the mapping.
func GetGameValue(key string) (any, bool) {
	value, ok := gameMap.values[key]
	return value, ok
}

// SetValue sets a value in the game map's value mapping. If the key already
// existed, the existing value will be overwritten.
func SetGameValue(key string, value any) {
	gameMap.values[key] = value
}

// DeleteValue removes a value from the map. If the key did not exist, this
// will be a no-op.
func DeleteGameValue(key string) {
	delete(gameMap.values, key)
}

// RunGame starts the game, given a starting state and ending state. This will
// return an error if the game gets to a point where there are no valid state
// transitions. Usually such an error indicates that the configuration of
// transitions between states is incorrect.
func RunGame(startState *GameState, endState *GameState) (any, error) {
	currentState := startState

	for currentState != endState {
		ret := currentState.action()
		nextState, err := currentState.nextState(ret)
		if err != nil {
			return nil, err
		}

		currentState = nextState
	}

	ret := currentState.action()
	return ret, nil
}

// RunGameToEnd starts the game, given a starting state. This will continue to
// run the game until no transition is possible, at which point it will
// return. This makes it possible to define multiple end states. It also makes
// it easier to configure invalid transitions between states, since such
// configurations will cause the game to end prematurely.
func RunGameToEnd(startState *GameState) any {
	currentState := startState

	for {
		ret := currentState.action()
		nextState, err := currentState.nextState(ret)
		if err != nil {
			return ret
		}

		currentState = nextState
	}
}
