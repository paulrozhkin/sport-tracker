package services

import (
	"context"
	"github.com/paulrozhkin/sport-tracker/config"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type DataSeedingService struct {
	needSeedData bool
	logger       *zap.SugaredLogger
	users        *UsersService
	exercises    *ExercisesService
	workouts     *WorkoutsService
	workoutPlans *WorkoutPlansService
	store        *infrastructure.Store
	userId       string
}

func NewDataSeedingService(lc fx.Lifecycle,
	cfg *config.Configuration,
	logger *zap.SugaredLogger,
	users *UsersService,
	exercises *ExercisesService,
	workouts *WorkoutsService,
	workoutPlans *WorkoutPlansService,
	store *infrastructure.Store) (*DataSeedingService, error) {
	seeder := &DataSeedingService{
		logger:       logger,
		needSeedData: cfg.Database.DataSeeding,
		users:        users,
		store:        store,
		exercises:    exercises,
		workouts:     workouts,
		workoutPlans: workoutPlans,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return seeder.seed()
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})
	return seeder, nil
}

func (s *DataSeedingService) seed() error {
	if s.isEmpty("users") && s.isEmpty("exercises") &&
		s.isEmpty("workouts") && s.isEmpty("workout_plans") {
		err := s.seedWorkouts()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DataSeedingService) seedUsers() error {
	name := "user"
	gender := models.UserGenderMale
	height := 180
	testUser := &models.User{Username: name, Name: &name, Gender: &gender, Height: &height, Password: "password"}
	testUser, err := s.users.CreateUser(*testUser)
	s.userId = testUser.Id
	return err
}

func (s *DataSeedingService) seedWorkouts() error {
	s.logger.Info("Start seeding data")
	err := s.seedUsers()
	if err != nil {
		return err
	}

	cardio := &models.Exercise{Name: "Кардио 3 минуты"}
	cardio, err = s.exercises.CreateExercise(s.modify(cardio))
	if err != nil {
		return err
	}

	benchPress := &models.Exercise{Name: "Жим лежа 3x10/5x8"}
	benchPress, err = s.exercises.CreateExercise(s.modify(benchPress))
	if err != nil {
		return err
	}

	sitUps := &models.Exercise{Name: "Приседания 3x10/5x8"}
	sitUps, err = s.exercises.CreateExercise(s.modify(sitUps))
	if err != nil {
		return err
	}

	deadlift := &models.Exercise{Name: "Становая 3x10/5x8"}
	deadlift, err = s.exercises.CreateExercise(s.modify(deadlift))
	if err != nil {
		return err
	}

	hitchEx1 := &models.Exercise{Name: "Планка (табата)"}
	hitchEx1, _ = s.exercises.CreateExercise(s.modify(hitchEx1))

	hitchEx2 := &models.Exercise{Name: "Пресс (табата)"}
	hitchEx2, _ = s.exercises.CreateExercise(s.modify(hitchEx2))

	hitchEx3 := &models.Exercise{Name: "Французский жим на трицепс"}
	hitchEx3, _ = s.exercises.CreateExercise(s.modify(hitchEx3))

	hitchEx4 := &models.Exercise{Name: "Предплечья. Подъем с упором о скамью"}
	hitchEx4, _ = s.exercises.CreateExercise(s.modify(hitchEx4))

	hitch := &models.Exercise{Name: "Заминка"}
	hitch.Complex = []*models.Exercise{hitchEx1, hitchEx2, hitchEx3, hitchEx4}
	hitch, err = s.exercises.CreateExercise(s.modify(hitch))
	if err != nil {
		return err
	}

	// Day 1
	// Warm up
	warmUpDay1Ex1 := &models.Exercise{Name: "10 Приседаний"}
	warmUpDay1Ex1, _ = s.exercises.CreateExercise(s.modify(warmUpDay1Ex1))

	warmUpDay1Ex2 := &models.Exercise{Name: "10 Выпады"}
	warmUpDay1Ex2, _ = s.exercises.CreateExercise(s.modify(warmUpDay1Ex2))

	warmUpDay1Ex3 := &models.Exercise{Name: "10 Ситапы"}
	warmUpDay1Ex3, _ = s.exercises.CreateExercise(s.modify(warmUpDay1Ex3))

	warmUpDay1Ex4 := &models.Exercise{Name: "6 Экстензий"}
	warmUpDay1Ex4, _ = s.exercises.CreateExercise(s.modify(warmUpDay1Ex4))

	warmUpDay1 := &models.Exercise{Name: "Разминка 2 круга"}
	warmUpDay1.Complex = []*models.Exercise{warmUpDay1Ex1, warmUpDay1Ex2, warmUpDay1Ex3, warmUpDay1Ex4}
	warmUpDay1, err = s.exercises.CreateExercise(s.modify(warmUpDay1))
	if err != nil {
		return err
	}

	// WOD
	cindyEx1 := &models.Exercise{Name: "5 подтягиваний"}
	cindyEx1, _ = s.exercises.CreateExercise(s.modify(cindyEx1))

	cindyEx2 := &models.Exercise{Name: "10 отжиманий"}
	cindyEx2, _ = s.exercises.CreateExercise(s.modify(cindyEx2))

	cindyEx3 := &models.Exercise{Name: "15 воздушных приседаний"}
	cindyEx3, _ = s.exercises.CreateExercise(s.modify(cindyEx3))

	cindyWOD := &models.Exercise{Name: "WOD: Синди (Cindy)"}
	cindyWOD.Complex = []*models.Exercise{cindyEx1, cindyEx2, cindyEx3}
	cindyWOD, err = s.exercises.CreateExercise(s.modify(cindyWOD))
	if err != nil {
		return err
	}

	// Workout 1
	workoutDay1 := &models.Workout{}
	workoutDay1.Complex = []*models.Exercise{cardio, warmUpDay1, sitUps, cindyWOD, hitch}
	workoutDay1, err = s.workouts.CreateWorkout(s.modifyW(workoutDay1))
	if err != nil {
		return err
	}

	// Day 2
	// Warm up
	warmUpDay2Ex1 := &models.Exercise{Name: "6 Становых тяг"}
	warmUpDay2Ex1, _ = s.exercises.CreateExercise(s.modify(warmUpDay2Ex1))

	warmUpDay2Ex2 := &models.Exercise{Name: "6 Взятий с виса"}
	warmUpDay2Ex2, _ = s.exercises.CreateExercise(s.modify(warmUpDay2Ex2))

	warmUpDay2Ex3 := &models.Exercise{Name: "6 Фронтальных приседаний"}
	warmUpDay2Ex3, _ = s.exercises.CreateExercise(s.modify(warmUpDay2Ex3))

	warmUpDay2Ex4 := &models.Exercise{Name: "8 Подъемов ног к перекладине"}
	warmUpDay2Ex4, _ = s.exercises.CreateExercise(s.modify(warmUpDay2Ex4))

	warmUpDay2 := &models.Exercise{Name: "Разминка 2 круга"}
	warmUpDay2.Complex = []*models.Exercise{warmUpDay2Ex1, warmUpDay2Ex2, warmUpDay2Ex3, warmUpDay2Ex4}
	warmUpDay2, err = s.exercises.CreateExercise(s.modify(warmUpDay2))
	if err != nil {
		return err
	}

	// WOD
	wod2Ex1 := &models.Exercise{Name: "Гребля 3 минуты"}
	wod2Ex1, _ = s.exercises.CreateExercise(s.modify(wod2Ex1))

	wod2Ex2 := &models.Exercise{Name: "10 Подтягиваний"}
	wod2Ex2, _ = s.exercises.CreateExercise(s.modify(wod2Ex2))

	wod2Ex3 := &models.Exercise{Name: "15 Махов гирей"}
	wod2Ex3, _ = s.exercises.CreateExercise(s.modify(wod2Ex3))

	wod2WOD := &models.Exercise{Name: "WOD: Гребля"}
	wod2WOD.Complex = []*models.Exercise{wod2Ex1, wod2Ex2, wod2Ex3}
	wod2WOD, err = s.exercises.CreateExercise(s.modify(wod2WOD))
	if err != nil {
		return err
	}

	// Workout 2
	workoutDay2 := &models.Workout{}
	workoutDay2.Complex = []*models.Exercise{cardio, warmUpDay2, deadlift, wod2WOD, hitch}
	workoutDay2, err = s.workouts.CreateWorkout(s.modifyW(workoutDay2))
	if err != nil {
		return err
	}

	// Day 3
	// Warm up
	warmUpDay3Ex1 := &models.Exercise{Name: "6 Отжиманий"}
	warmUpDay3Ex1, _ = s.exercises.CreateExercise(s.modify(warmUpDay3Ex1))

	warmUpDay3Ex2 := &models.Exercise{Name: "6 Экстенизий"}
	warmUpDay3Ex2, _ = s.exercises.CreateExercise(s.modify(warmUpDay3Ex2))

	warmUpDay3Ex3 := &models.Exercise{Name: "6 Ситапов"}
	warmUpDay3Ex3, _ = s.exercises.CreateExercise(s.modify(warmUpDay3Ex3))

	warmUpDay3Ex4 := &models.Exercise{Name: "6 Жим стоя"}
	warmUpDay3Ex4, _ = s.exercises.CreateExercise(s.modify(warmUpDay3Ex4))

	warmUpDay3Ex5 := &models.Exercise{Name: "6 Тяга в наклоне"}
	warmUpDay3Ex5, _ = s.exercises.CreateExercise(s.modify(warmUpDay3Ex5))

	warmUpDay3 := &models.Exercise{Name: "Разминка 1 круг"}
	warmUpDay3.Complex = []*models.Exercise{warmUpDay3Ex1, warmUpDay3Ex2, warmUpDay3Ex3, warmUpDay3Ex4, warmUpDay3Ex5}
	warmUpDay3, err = s.exercises.CreateExercise(s.modify(warmUpDay3))
	if err != nil {
		return err
	}

	// WOD
	franEx1 := &models.Exercise{Name: "Трастеры"}
	franEx1, _ = s.exercises.CreateExercise(s.modify(franEx1))

	franEx2 := &models.Exercise{Name: "Подтягивания"}
	franEx2, _ = s.exercises.CreateExercise(s.modify(franEx2))

	franWOD := &models.Exercise{Name: "WOD: Фрэн (Fran) 21-15-9"}
	franWOD.Complex = []*models.Exercise{franEx1, franEx2}
	franWOD, err = s.exercises.CreateExercise(s.modify(franWOD))
	if err != nil {
		return err
	}

	// Workout 3
	workoutDay3 := &models.Workout{}
	workoutDay3.Complex = []*models.Exercise{cardio, warmUpDay3, benchPress, franWOD, hitch}
	workoutDay3, err = s.workouts.CreateWorkout(s.modifyW(workoutDay3))
	if err != nil {
		return err
	}

	// Day 4
	// WOD
	lasEx1 := &models.Exercise{Name: "10+10 бег в упоре лежа"}
	lasEx1, _ = s.exercises.CreateExercise(s.modify(lasEx1))

	lasEx2 := &models.Exercise{Name: "10 воздушных приседаний"}
	lasEx2, _ = s.exercises.CreateExercise(s.modify(lasEx2))

	lasEx3 := &models.Exercise{Name: "10 отжиманий \"паук\"(отведение ноги в сторону)"}
	lasEx3, _ = s.exercises.CreateExercise(s.modify(lasEx3))

	lasEx4 := &models.Exercise{Name: "10 бёрпи"}
	lasEx4, _ = s.exercises.CreateExercise(s.modify(lasEx4))

	lasEx5 := &models.Exercise{Name: "10 пресс (Ситап-пресс)"}
	lasEx5, _ = s.exercises.CreateExercise(s.modify(lasEx5))

	lasEx6 := &models.Exercise{Name: "10 прыжков на тумбу (50-60 см)"}
	lasEx6, _ = s.exercises.CreateExercise(s.modify(lasEx6))

	lasWOD := &models.Exercise{Name: "WOD: ЛАС (LAS) 7 раундов"}
	lasWOD.Complex = []*models.Exercise{lasEx1, lasEx2, lasEx3, lasEx4, lasEx5, lasEx6}
	lasWOD, err = s.exercises.CreateExercise(s.modify(lasWOD))
	if err != nil {
		return err
	}

	// Workout 4
	workoutDay4 := &models.Workout{}
	workoutDay4.Complex = []*models.Exercise{cardio, warmUpDay1, sitUps, lasWOD, hitch}
	workoutDay4, err = s.workouts.CreateWorkout(s.modifyW(workoutDay4))
	if err != nil {
		return err
	}

	// Day 5
	// WOD
	fgbEx1 := &models.Exercise{Name: "1 минута - броски мяча в цель, 9 кг"}
	fgbEx1, _ = s.exercises.CreateExercise(s.modify(fgbEx1))

	fgbEx2 := &models.Exercise{Name: "1 минута - тяга штанги сумо, 35 кг"}
	fgbEx2, _ = s.exercises.CreateExercise(s.modify(fgbEx2))

	fgbEx3 := &models.Exercise{Name: "1 минута - прыжки на коробку"}
	fgbEx3, _ = s.exercises.CreateExercise(s.modify(fgbEx3))

	fgbEx4 := &models.Exercise{Name: "1 минута - толчковый швунг штанги, 35 кг"}
	fgbEx4, _ = s.exercises.CreateExercise(s.modify(fgbEx4))

	fgbEx5 := &models.Exercise{Name: "1 минута гребли"}
	fgbEx5, _ = s.exercises.CreateExercise(s.modify(fgbEx5))

	fgbWOD := &models.Exercise{Name: "WOD: ФГБ (Fight Gone Bad) 3 раунда"}
	fgbWOD.Complex = []*models.Exercise{fgbEx1, fgbEx2, fgbEx3, fgbEx4, fgbEx5}
	fgbWOD, err = s.exercises.CreateExercise(s.modify(fgbWOD))
	if err != nil {
		return err
	}

	// Workout 4
	workoutDay5 := &models.Workout{}
	workoutDay5.Complex = []*models.Exercise{cardio, warmUpDay2, deadlift, fgbWOD, hitch}
	workoutDay5, err = s.workouts.CreateWorkout(s.modifyW(workoutDay5))
	if err != nil {
		return err
	}

	// Day 6
	// WOD
	elizabethEx1 := &models.Exercise{Name: "Взятие штанги на грудь 60 кг"}
	elizabethEx1, _ = s.exercises.CreateExercise(s.modify(elizabethEx1))

	elizabethEx2 := &models.Exercise{Name: "Отжимания на кольцах"}
	elizabethEx2, _ = s.exercises.CreateExercise(s.modify(elizabethEx2))

	elizabethWOD := &models.Exercise{Name: "WOD: Элизабет (Elizabeth)  21-15-9"}
	elizabethWOD.Complex = []*models.Exercise{elizabethEx1, elizabethEx2}
	elizabethWOD, err = s.exercises.CreateExercise(s.modify(elizabethWOD))
	if err != nil {
		return err
	}
	// Workout 6
	workoutDay6 := &models.Workout{}
	workoutDay6.Complex = []*models.Exercise{cardio, warmUpDay3, benchPress, elizabethWOD, hitch}
	workoutDay6, err = s.workouts.CreateWorkout(s.modifyW(workoutDay6))
	if err != nil {
		return err
	}

	// Workout plan
	plan := &models.WorkoutPlan{Name: "Силовая и кроссфит"}
	plan.Repeatable = true
	plan.Workouts = []*models.Workout{workoutDay1, workoutDay2, workoutDay3, workoutDay4, workoutDay5, workoutDay6}
	_, err = s.workoutPlans.CreateWorkoutPlan(s.modifyP(plan))
	if err != nil {
		return err
	}

	s.logger.Info("Data seeded")
	return err
}

func (s *DataSeedingService) modify(exercise *models.Exercise) models.Exercise {
	exercise.Owner = s.userId
	return *exercise
}

func (s *DataSeedingService) modifyW(workout *models.Workout) models.Workout {
	workout.Owner = s.userId
	return *workout
}

func (s *DataSeedingService) modifyP(plan *models.WorkoutPlan) models.WorkoutPlan {
	plan.Owner = s.userId
	return *plan
}

func (s *DataSeedingService) isEmpty(tableName string) bool {
	query := "SELECT COUNT(*) FROM " + tableName
	var count int64
	err := s.store.Pool.QueryRow(context.Background(), query).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}
