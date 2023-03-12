BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id       uuid PRIMARY KEY,
    username text NOT NULL,
    password text NOT NULL,
    name     text,
    gender   numeric,
    height   numeric,
    created  timestamp NOT NULL,
    updated  timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS exercises
(
    id                uuid PRIMARY KEY,
    name              text NOT NULL,
    short_description text,
    owner             uuid references users (id) NOT NULL,
    complex           uuid[],
    created           timestamp NOT NULL,
    updated           timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS workout_plans
(
    id                uuid PRIMARY KEY,
    name              text NOT NULL,
    short_description text,
    repeatable        boolean NOT NULL,
    owner             uuid references users (id) NOT NULL,
    created           timestamp NOT NULL,
    updated           timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts
(
    id            uuid PRIMARY KEY,
    plan_id       uuid references workout_plans (id) NOT NULL,
    order_in_plan numeric NOT NULL,
    complex       uuid[],
    created       timestamp NOT NULL,
    updated       timestamp NOT NULL,
    UNIQUE (plan_id, order_in_plan)
);

CREATE TABLE IF NOT EXISTS user_workouts
(
    id           uuid PRIMARY KEY,
    user_id      uuid references users (id) NOT NULL,
    workout_plan uuid references workout_plans (id) NOT NULL,
    active       boolean NOT NULL,
    schedule     int[] NOT NULL,
    created      timestamp NOT NULL,
    updated      timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts_statistic
(
    id           uuid PRIMARY KEY,
    user_workout uuid references user_workouts (id) NOT NULL,
    workout      uuid references workouts (id) NOT NULL,
    skipped      boolean NOT NULL,
    workout_date timestamp NOT NULL,
    created      timestamp NOT NULL,
    updated      timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS weight_statistic
(
    id      uuid PRIMARY KEY,
    user_id uuid references users (id) NOT NULL,
    weight  decimal NOT NULL,
    created timestamp NOT NULL,
    updated timestamp NOT NULL
);

CREATE INDEX IF NOT EXISTS workouts_statistic_user_workout_idx ON workouts_statistic (user_workout);
CREATE INDEX IF NOT EXISTS workouts_plan_id_order_in_plan_idx ON workouts (plan_id, order_in_plan);

COMMIT;