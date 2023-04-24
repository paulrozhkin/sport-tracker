BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id       uuid PRIMARY KEY,
    username text      NOT NULL,
    password text      NOT NULL,
    name     text,
    gender   numeric,
    height   numeric,
    created  timestamp NOT NULL,
    updated  timestamp NOT NULL,
    UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS exercises
(
    id                uuid PRIMARY KEY,
    name              text                       NOT NULL,
    short_description text,
    owner             uuid references users (id) NOT NULL,
    complex           uuid[],
    created           timestamp                  NOT NULL,
    updated           timestamp                  NOT NULL
);

CREATE TABLE IF NOT EXISTS workout_plans
(
    id                uuid PRIMARY KEY,
    name              text                       NOT NULL,
    short_description text,
    repeatable        boolean                    NOT NULL,
    owner             uuid references users (id) NOT NULL,
    workouts          uuid[],
    created           timestamp                  NOT NULL,
    updated           timestamp                  NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts
(
    id                 uuid PRIMARY KEY,
    custom_name        text,
    custom_description text,
    complex            uuid[],
    owner              uuid references users (id) NOT NULL,
    created            timestamp                  NOT NULL,
    updated            timestamp                  NOT NULL
);

CREATE TABLE IF NOT EXISTS user_workouts
(
    id           uuid PRIMARY KEY,
    user_id      uuid references users (id)         NOT NULL,
    workout_plan uuid references workout_plans (id) NOT NULL,
    active       boolean                            NOT NULL,
    schedule     int[]                              NOT NULL,
    created      timestamp                          NOT NULL,
    updated      timestamp                          NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts_statistic
(
    id             uuid PRIMARY KEY,
    user_workout   uuid references user_workouts (id) NOT NULL,
    workout        uuid references workouts (id)      NOT NULL,
    workout_date   timestamp,
    scheduled_date timestamp                          NOT NULL,
    created        timestamp                          NOT NULL,
    updated        timestamp                          NOT NULL
);

CREATE TABLE IF NOT EXISTS weight_statistic
(
    id      uuid PRIMARY KEY,
    user_id uuid references users (id) NOT NULL,
    weight  decimal                    NOT NULL,
    created timestamp                  NOT NULL,
    updated timestamp                  NOT NULL
);

CREATE INDEX IF NOT EXISTS users_username_idx ON users (username);
CREATE INDEX IF NOT EXISTS workouts_statistic_user_workout_idx ON workouts_statistic (user_workout);

COMMIT;