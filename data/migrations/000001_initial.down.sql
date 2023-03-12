BEGIN;

DROP INDEX IF EXISTS workouts_statistic_user_workout_idx;
DROP INDEX IF EXISTS workouts_plan_id_order_in_plan_idx;

DROP TABLE IF EXISTS exercises, workout_plans, workouts, user_workouts,
    weight_statistic, workouts_statistic, users CASCADE;

COMMIT;