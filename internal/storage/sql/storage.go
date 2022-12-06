package sqlstorage

import (
	"context"
	"github.com/a-klimenko/go-otus-final-project/internal/storage"
	"github.com/a-klimenko/go-otus-final-project/internal/ucb"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	// import postgres lib.
	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sqlx.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) AddBanner(ctx context.Context, bannerId uuid.UUID, slotId uuid.UUID) error {
	updateQuery := `
				UPDATE rotations 
					SET deleted_at=NULL
				WHERE id IN (
				    SELECT id FROM rotations WHERE banner_id=$1 AND slot_id=$2 
				)
	`
	_, err := s.Db.ExecContext(
		ctx,
		updateQuery,
		bannerId,
		slotId,
	)
	if err != nil {
		return err
	}

	groupQuery := `
				SELECT id, description
				FROM groups
				WHERE id NOT IN (
				    SELECT group_id FROM rotations WHERE banner_id=$1 AND slot_id=$2
				)
	`
	rows, err := s.Db.QueryxContext(
		ctx, groupQuery,
		bannerId,
		slotId,
	)
	if err != nil {
		return err
	}

	insertQuery := `
				INSERT INTO rotations 
					(id, banner_id, slot_id, group_id)
				VALUES
					($1, $2, $3, $4)
	`
	for rows.Next() {
		var group storage.Group
		err := rows.StructScan(&group)
		if err != nil {
			return err
		}
		_, err = s.Db.ExecContext(
			ctx,
			insertQuery,
			uuid.New(),
			bannerId,
			slotId,
			group.ID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) RemoveBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error {
	query := `
				UPDATE rotations 
					SET deleted_at=now()
				WHERE id IN (
				    SELECT id FROM rotations WHERE banner_id=$1 AND slot_id=$2 
				)
	`
	_, err := s.Db.ExecContext(
		ctx,
		query,
		bannerId,
		slotId,
	)

	return err
}

func (s *Storage) ClickBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID, groupId uuid.UUID) error {
	query := `
				UPDATE rotations 
					SET clicks=clicks+1
				WHERE banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	_, err := s.Db.ExecContext(
		ctx,
		query,
		bannerId,
		slotId,
		groupId,
	)

	return err
}

func (s *Storage) ChooseBanner(ctx context.Context, slotId uuid.UUID, groupId uuid.UUID) (*uuid.UUID, error) {
	rotationsQuery := `
				SELECT id, banner_id, slot_id, group_id, clicks, shows
				FROM rotations 
				WHERE deleted_at IS NULL AND slot_id=$1 AND group_id=$2
	`
	rows, err := s.Db.QueryxContext(ctx, rotationsQuery, slotId, groupId)
	if err != nil {
		return nil, err
	}

	totalShows := 0
	rotations := make(map[uuid.UUID]storage.Rotation, 0)
	for rows.Next() {
		var rotation storage.Rotation
		err := rows.StructScan(&rotation)
		if err != nil {
			return nil, err
		}
		rotations[rotation.BannerId] = rotation
		totalShows += rotation.Shows
	}

	targetBannerId := ucb.MakeDecision(rotations, totalShows)

	query := `
				UPDATE rotations 
					SET shows=shows+1
				WHERE banner_id=$1 AND slot_id=$2 AND group_id=$3
	`
	_, err = s.Db.ExecContext(
		ctx,
		query,
		targetBannerId,
		slotId,
		groupId,
	)
	if err != nil {
		return nil, err
	}

	return &targetBannerId, nil
}

func (s *Storage) Connect() error {
	db, err := sqlx.Open("postgres",
		"postgres://admin:admin@rotator-Db:5432/rotator?sslmode=disable",
	)
	if err != nil {
		return err
	}

	s.Db = db

	return nil
}

func (s *Storage) Close() error {
	if err := s.Db.Close(); err != nil {
		return err
	}

	return nil
}
