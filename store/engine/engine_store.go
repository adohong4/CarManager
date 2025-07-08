package engine

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/adohong4/carZone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type EngineSstore struct {
	db *sql.DB
}

func New(db *sql.DB) *EngineSstore {
	return &EngineSstore{db: db}
}

func (e EngineSstore) EngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "EngineById-Store")
	defer span.End()

	var engine models.Engine

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return engine, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id = $1", id).Scan(
		&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, nil // No rows found, return empty engine
		}
		return engine, err
	}
	return engine, err
}

func (e EngineSstore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "CreateEngine-Store")
	defer span.End()

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	engineID := uuid.New()

	_, err = tx.ExecContext(ctx,
		"INSERT INTO engine (id, displacement, no_of_cylinders, car_range) VALUES ($1, $2, $3, $4)",
		engineID, engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange,
	)
	if err != nil {
		return models.Engine{}, err
	}

	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil
}

func (e EngineSstore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "EngineUpdate-Store")
	defer span.End()

	engineID, err := uuid.Parse(id)
	if err != nil {
		return models.Engine{}, fmt.Errorf("invalid engine ID: %w", err)
	}

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	results, err := tx.ExecContext(ctx,
		"UPDATE engine SET displacement = $1, no_of_cylinders = $2, car_range = $3 WHERE id = $4",
		engineReq.Displacement, engineReq.NoOfCylinders, engineReq.CarRange)

	if err != nil {
		return models.Engine{}, err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("No Rows Were Updated")
	}

	engine := models.Engine{
		EngineID:      engineID,
		Displacement:  engineReq.Displacement,
		NoOfCylinders: engineReq.NoOfCylinders,
		CarRange:      engineReq.CarRange,
	}

	return engine, nil
}

func (e EngineSstore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")
	ctx, span := tracer.Start(ctx, "EngineDelete-Store")
	defer span.End()

	var engine models.Engine

	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				fmt.Printf("Transaction rollback error: %v\n", rbErr)
			}
		} else {
			if cmErr := tx.Commit(); cmErr != nil {
				fmt.Printf("Transaction commit error: %v\n", cmErr)
			}
		}
	}()

	err = tx.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id = $1", id).Scan(
		&engine.EngineID, &engine.Displacement, &engine.NoOfCylinders, &engine.CarRange,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return engine, nil // No rows found, return empty engine
		}
		return engine, err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM engine WHERE id = $1", id)
	if err != nil {
		return models.Engine{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Engine{}, err
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("No Rows Were Deleted")
	}

	return engine, nil
}
