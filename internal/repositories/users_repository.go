package repositories

import (
	"database/sql"
	"errors"
	"github.com/paulrozhkin/sport-tracker/internal/infrastructure"
	"github.com/paulrozhkin/sport-tracker/internal/models"
	"go.uber.org/zap"
)

type UsersRepository struct {
	store *infrastructure.Store
	log   *zap.SugaredLogger
}

func NewUsersRepository(store *infrastructure.Store,
	logger *zap.SugaredLogger) (*UsersRepository, error) {
	return &UsersRepository{store: store, log: logger}, nil
}

func (ur *UsersRepository) CreateUser(user models.User) (*models.User, error) {
	user.FillForCreate()
	query := `INSERT INTO users (id, created, updated, username, password, name, gender, height) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := ur.store.Db.Exec(query, user.Id, user.Created, user.Updated,
		user.Username, user.Password, user.Name, user.Gender, user.Height)
	if err != nil {
		ur.log.Error("Failed to create user", err)
		return nil, err
	}
	return ur.GetUserById(user.Id)
}
func (ur *UsersRepository) GetUserById(id string) (*models.User, error) {
	query := `SELECT (id, created, updated, username, password, name, gender, height) 
				FROM users WHERE id=$1`
	row := ur.store.Db.QueryRow(query, id)
	user, err := rowToUser(row)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, models.NewNotFoundByIdError("user", id)
	} else if err != nil {
		ur.log.Errorf("Failed to get user by id %s due to: %v", id, err)
		return nil, err
	}
	return user, nil
}

func (ur *UsersRepository) GetUserByUsername(username string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username=$1"
	row := ur.store.Db.QueryRow(query, username)
	user, err := rowToUser(row)
	if err != nil && errors.Is(sql.ErrNoRows, err) {
		return nil, models.NewNotFoundError("user", username, "username")
	} else if err != nil {
		ur.log.Errorf("Failed to get user by username %s due to: %v", username, err)
		return nil, err
	}
	return user, nil
}

func rowToUser(row *sql.Row) (*models.User, error) {
	user := &models.User{}

	err := row.Scan(&user.Id, &user.Created,
		&user.Updated, &user.Username,
		&user.Password, &user.Name, &user.Gender,
		&user.Height)
	if err != nil {
		return nil, err
	}
	return user, nil
}
