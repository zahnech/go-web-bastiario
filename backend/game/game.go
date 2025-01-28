package game

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID string `json:"id"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

func InitPlayer() *Player {
	localRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Player{
		uuid.NewString(),
		int(localRand.Float64() * 800),
		int(localRand.Float64() * 600),
	}
}
