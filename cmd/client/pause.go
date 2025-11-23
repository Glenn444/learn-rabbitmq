package main

import "sync"

type UnitRank string
type Location string

type Unit struct {
	ID       int
	Rank     UnitRank
	Location Location
}

type Player struct {
	Username string
	Units    map[int]Unit
}


type GameState struct {
	Player Player
	Paused bool
	mu     *sync.RWMutex
}

func handlerPause(gameState GameState){
	
}