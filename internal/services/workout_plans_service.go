package services

import (
	"fmt"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"github.com/paulrozhkin/sport-tracker/internal/repositories"
)

type WorkoutPlansService struct {
	workoutPlanRepository *repositories.WorkoutPlansRepository
}

func NewWorkoutPlansService(workoutPlanRepository *repositories.WorkoutPlansRepository) (*WorkoutPlansService, error) {
	return &WorkoutPlansService{workoutPlanRepository: workoutPlanRepository}, nil
}

func (ws *WorkoutPlansService) CreateWorkoutPlan(workoutPlan models.WorkoutPlan) (*models.WorkoutPlan, error) {
	if workoutPlan.Owner == "" {
		return nil, fmt.Errorf("workoutPlan owner %s in CreateWorkoutPlan", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutPlanRepository.CreateWorkoutPlan(workoutPlan)
}

func (ws *WorkoutPlansService) UpdateWorkoutPlan(workoutPlan models.WorkoutPlan) (*models.WorkoutPlan, error) {
	if workoutPlan.Id == "" {
		return nil, fmt.Errorf("workoutPlan id %s in UpdateWorkoutPlan", models.ArgumentNullOrEmptyError)
	}
	originalWorkoutPlan, err := ws.workoutPlanRepository.GetWorkoutPlansByIdWithoutComplex(workoutPlan.Id)
	if err != nil {
		return nil, err
	}
	if originalWorkoutPlan.Owner != workoutPlan.Owner {
		return nil, models.NewNoRightsOnEntityError("workoutPlan", workoutPlan.Id)
	}
	return ws.workoutPlanRepository.UpdateWorkoutPlan(workoutPlan)
}

func (ws *WorkoutPlansService) GetWorkoutPlanById(workoutPlanId string) (*models.WorkoutPlan, error) {
	if workoutPlanId == "" {
		return nil, fmt.Errorf("workoutPlan name %s in GetWorkoutPlanById", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutPlanRepository.GetWorkoutPlanById(workoutPlanId)
}

func (ws *WorkoutPlansService) GetWorkoutPlans() ([]*models.WorkoutPlan, error) {
	return ws.workoutPlanRepository.GetWorkoutPlans()
}

func (ws *WorkoutPlansService) DeleteWorkoutPlanById(workoutPlanId string) error {
	if workoutPlanId == "" {
		return fmt.Errorf("workoutPlan name %s in DeleteWorkoutPlanById", models.ArgumentNullOrEmptyError)
	}
	return ws.workoutPlanRepository.DeleteWorkoutPlanById(workoutPlanId)
}
