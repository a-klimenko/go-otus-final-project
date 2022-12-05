package app

import (
	"context"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	Connect() error
	Close() error
	AddBanner(ctx context.Context, bannerId uuid.UUID, slotId uuid.UUID) error
	RemoveBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error
	ClickBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID, groupId uuid.UUID) error
	ChooseBanner(ctx context.Context, slotId uuid.UUID, groupId uuid.UUID) (*uuid.UUID, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) AddBanner(ctx context.Context, bannerId uuid.UUID, slotId uuid.UUID) error {
	return a.storage.AddBanner(ctx, bannerId, slotId)
}

func (a *App) RemoveBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID) error {
	return a.storage.RemoveBanner(ctx, slotId, bannerId)
}

func (a *App) ClickBanner(ctx context.Context, slotId uuid.UUID, bannerId uuid.UUID, groupId uuid.UUID) error {
	return a.storage.ClickBanner(ctx, slotId, bannerId, groupId)
}
func (a *App) ChooseBanner(ctx context.Context, slotId uuid.UUID, groupId uuid.UUID) (*uuid.UUID, error) {
	return a.storage.ChooseBanner(ctx, slotId, groupId)
}
