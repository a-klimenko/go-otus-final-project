package ucb

import (
	"github.com/a-klimenko/go-otus-final-project/internal/storage"
	"github.com/google/uuid"
	"math"
)

func MakeDecision(rotations map[uuid.UUID]storage.Rotation, totalShows int) uuid.UUID {
	max := 0.0
	var targetBannerId uuid.UUID
	for bannerId, rotation := range rotations {
		avReward := float64(rotation.Clicks) / float64(rotation.Shows)
		decision := avReward + math.Sqrt(2*math.Log(float64(totalShows))/float64(rotation.Shows))
		if decision >= max {
			targetBannerId = bannerId
		}
	}

	return targetBannerId
}
