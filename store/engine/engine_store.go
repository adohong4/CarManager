package engine

import (
	"context"
	"database/sql"

	"github.com/adohong4/carZone/models"
)

type EngineSstore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineSstore {
	return &EngineSstore{db: db}
}

func (e EngineSstore) EngineById(ctx context.Context, id string) (models.Engine, error) {

}

func (e EngineSstore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {

}

func (e EngineSstore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {

}

func (e EngineSstore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {

}
