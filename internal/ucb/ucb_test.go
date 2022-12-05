package ucb

import (
	"github.com/a-klimenko/go-otus-final-project/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMakeDecision(t *testing.T) {
	t.Run("choose any", func(t *testing.T) {
		rotations := make(map[uuid.UUID]storage.Rotation, 3)

		for i := 0; i < 3; i++ {
			rotations[uuid.New()] = storage.Rotation{
				Clicks: 10,
				Shows:  30,
			}
		}
		shows := make(map[uuid.UUID]int, 3)
		totalShows := 90
		for i := 0; i < 1000; i++ {
			selectedBanner := MakeDecision(rotations, totalShows)
			totalShows += 1
			shows[selectedBanner] += 1
		}
		require.True(t, len(shows) == 3)
	})

	t.Run("choose popular", func(t *testing.T) {
		rotations := make(map[uuid.UUID]storage.Rotation, 3)

		for i := 0; i < 3; i++ {
			rotations[uuid.New()] = storage.Rotation{
				Clicks: 10,
				Shows:  30,
			}
		}
		popularBannerId := uuid.New()
		rotations[popularBannerId] = storage.Rotation{
			Clicks: 1000,
			Shows:  30,
		}
		shows := make(map[uuid.UUID]int, 4)
		totalShows := 120
		for i := 0; i < 1000; i++ {
			selectedBanner := MakeDecision(rotations, totalShows)
			totalShows += 1
			shows[selectedBanner] += 1
		}
		res := shows[popularBannerId]
		success := true
		for _, v := range shows {
			if v > res {
				success = false
				break
			}
		}
		require.True(t, success)
	})
}
