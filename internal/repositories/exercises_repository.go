package repositories

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
)

type ExercisesRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewExercisesRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*ExercisesRepository, error) {
	return &ExercisesRepository{store: store, log: logger}, nil
}

// CreateExercise Create new exercise
func (er *ExercisesRepository) CreateExercise(exercise models.Exercise) (*models.Exercise, error) {
	exercise.FillForCreate()
	query := `INSERT INTO exercises (id, created, updated, name, short_description, owner, complex) 
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	var complexIds []string
	if exercise.Complex != nil {
		for _, complexExercise := range exercise.Complex {
			complexIds = append(complexIds, complexExercise.Id)
		}
	}
	_, err := er.store.Pool.Exec(context.Background(), query, exercise.Id, exercise.Created, exercise.Updated,
		exercise.Name, exercise.ShortDescription, exercise.Owner, complexIds)
	if err != nil {
		er.log.Error("Failed to create exercise", err)
		return nil, err
	}
	return er.GetExerciseById(exercise.Id)
}

func (er *ExercisesRepository) UpdateExercise(exercise models.Exercise) (*models.Exercise, error) {
	exercise.FillForUpdate()
	query := `UPDATE exercises SET updated=$2, name=$3, short_description=$4, complex=$5 WHERE id=$1`
	var complexIds []string
	if exercise.Complex != nil {
		for _, complexExercise := range exercise.Complex {
			complexIds = append(complexIds, complexExercise.Id)
		}
	}
	_, err := er.store.Pool.Exec(context.Background(), query, exercise.Id, exercise.Updated,
		exercise.Name, exercise.ShortDescription, complexIds)
	if err != nil {
		er.log.Error("Failed to update exercise", err)
		return nil, err
	}
	return er.GetExerciseById(exercise.Id)
}

// GetExerciseById Get exercise by id with filled complex
func (er *ExercisesRepository) GetExerciseById(id string) (*models.Exercise, error) {
	query := `WITH RECURSIVE descendents AS (
				  SELECT id, unnest(complex) AS child
				  FROM exercises o
				  WHERE id=$1
					UNION ALL
				  SELECT o.id, unnest(o.complex) AS child
					FROM exercises o
						JOIN descendents d ON o.id = d.child
				
				)
				SELECT child, e.created, e.updated, e.name, e.short_description, e.owner,
						ARRAY [descendents.id] as parentId FROM descendents
					JOIN exercises e on descendents.child=e.id
				UNION ALL
				SELECT id, e.created, e.updated, e.name, e.short_description, e.owner, ARRAY [id] as parentId
				  FROM exercises e
				  WHERE id=$1
				`
	rows, err := er.store.Pool.Query(context.Background(), query, id)
	if err != nil {
		er.log.Errorf("Failed to get exercise by id %s due to: %v", id, err)
		return nil, err
	}

	var allExercises []*models.Exercise
	for rows.Next() {
		exercise, rowScanErr := rowToExercise(rows)
		if rowScanErr != nil && errors.Is(pgx.ErrNoRows, rowScanErr) {
			return nil, models.NewNotFoundByIdError("exercise", id)
		} else if rowScanErr != nil {
			er.log.Errorf("Failed to get user by id %s due to: %v", id, rowScanErr)
			return nil, rowScanErr
		}
		allExercises = append(allExercises, exercise)
	}
	if len(allExercises) == 0 {
		return nil, models.NewNotFoundByIdError("exercise", id)
	}
	return buildExerciseTree(createCopyWithoutLink(getRoot(allExercises)), allExercises), nil
}

// GetExercisesByIdWithoutComplex Get exercise  by id without filled complex (only ids)
func (er *ExercisesRepository) GetExercisesByIdWithoutComplex(id string) (*models.Exercise, error) {
	query := `SELECT id, created, updated, name, short_description, owner, complex
				FROM exercises WHERE id=$1`
	row := er.store.Pool.QueryRow(context.Background(), query, id)
	result, err := rowToExercise(row)
	if err != nil && errors.Is(pgx.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("exercise", id)
	} else if err != nil {
		er.log.Errorf("Failed to get exercise by id %s due to: %v", id, err)
		return nil, err
	}
	return result, nil
}

// GetExercises Get exercises without filled complex (only ids)
func (er *ExercisesRepository) GetExercises() ([]*models.Exercise, error) {
	query := `SELECT id, created, updated, name, short_description, owner, complex
				FROM exercises`
	rows, err := er.store.Pool.Query(context.Background(), query)
	if err != nil {
		er.log.Errorf("Failed to get exercises due to: %v", err)
		return nil, err
	}
	var result []*models.Exercise
	for rows.Next() {
		exercise, rowScanErr := rowToExercise(rows)
		if rowScanErr != nil {
			er.log.Errorf("Failed to scan exercises due to: %v", rowScanErr)
			continue
		}
		result = append(result, exercise)
	}
	return result, nil
}

func (er *ExercisesRepository) DeleteExerciseById(id string) error {
	query := `DELETE FROM exercises WHERE id = $1;`
	res, err := er.store.Pool.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return models.NewNotFoundByIdError("exercise", id)
	}
	return nil
}

func rowToExercise(row pgx.Row) (*models.Exercise, error) {
	exercise := &models.Exercise{}
	shortDescription := sql.NullString{}
	var exerciseComplex []string
	err := row.Scan(&exercise.Id, &exercise.Created,
		&exercise.Updated, &exercise.Name,
		&shortDescription, &exercise.Owner, &exerciseComplex)
	if err != nil {
		return nil, err
	}
	if shortDescription.Valid {
		exercise.ShortDescription = &shortDescription.String
	}
	for _, id := range exerciseComplex {
		internalExercise := new(models.Exercise)
		internalExercise.Id = id
		exercise.Complex = append(exercise.Complex, internalExercise)
	}
	return exercise, nil
}

func buildExerciseTree(root *models.Exercise, allExercises []*models.Exercise) *models.Exercise {
	for _, exerciseTreeItem := range allExercises {
		// exerciseTreeItem.Complex[0] id is parent value in tree, can't be null
		// exerciseTreeItem.Complex[0].Id != exerciseTreeItem.Id - is root
		if exerciseTreeItem.Complex[0].Id == root.Id && exerciseTreeItem.Complex[0].Id != exerciseTreeItem.Id {
			child := createCopyWithoutLink(exerciseTreeItem)
			root.Complex = append(root.Complex, child)
			buildExerciseTree(child, allExercises)
		}
	}
	return root
}

func getRoot(allExercises []*models.Exercise) *models.Exercise {
	for _, exerciseTreeItem := range allExercises {
		if exerciseTreeItem.Complex[0].Id == exerciseTreeItem.Id {
			return exerciseTreeItem
		}
	}
	return nil
}

func createCopyWithoutLink(exerciseForCopy *models.Exercise) *models.Exercise {
	result := *exerciseForCopy
	result.Complex = nil
	return &result
}
