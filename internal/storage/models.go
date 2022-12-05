package storage

import (
	"database/sql"
	"github.com/google/uuid"
)

type Banner struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type Group struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type Slot struct {
	ID          uuid.UUID `db:"id"`
	Description string    `db:"description"`
}

type Rotation struct {
	ID        uuid.UUID    `db:"id"`
	BannerId  uuid.UUID    `db:"banner_id"`
	SlotId    uuid.UUID    `db:"slot_id"`
	GroupId   uuid.UUID    `db:"group_id"`
	Shows     int          `db:"shows"`
	Clicks    int          `db:"clicks"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
