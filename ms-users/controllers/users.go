package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/bsanzhiev/bahamas/ms-users/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserController struct {
	DBPool *pgxpool.Pool
	Ctx    context.Context
}

// GetUsers How Kafka push to work this code?
func (uc *UserController) GetUsers() ([]types.User, error) {
	ctx := uc.Ctx
	dbPool := uc.DBPool
	// Делаем запрос
	rows, errRows := dbPool.Query(ctx, "SELECT id, username, first_name, last_name, email FROM users")
	if errRows != nil {
		return nil, errRows
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		if errUsers := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email); errUsers != nil {
			return nil, errUsers
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *UserController) UserByID(userID int) (types.User, error) {
	ctx := uc.Ctx
	dbPool := uc.DBPool
	var user types.User

	query := "SELECT id, username, first_name, last_name, email FROM users WHERE id=$1"
	err := dbPool.QueryRow(ctx, query, userID).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return types.User{}, fmt.Errorf("no user found with id: %d", userID)
		}
		return types.User{}, err
	}
	return user, nil
}

func (uc *UserController) UserCreate(userData types.UserRequestData) error {
	query := "INSERT INTO users (username, first_name, last_name, email) VALUES ($1, $2, $3, $4)"
	_, err := uc.DBPool.Exec(uc.Ctx, query, userData.Username, userData.FirstName, userData.LastName, userData.Email)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}
	return nil
}

func (uc *UserController) UserUpdate(userID int, userData types.UserRequestData) error {
	query := `
	UPDATE users
	SET username = COALESCE(NULLIF($1, ''), username),
	    first_name = COALESCE(NULLIF($2, ''), first_name),
	    last_name = COALESCE(NULLIF($3, ''), last_name),
	    email = COALESCE(NULLIF($4, ''), email)
	WHERE id = $5`

	cmdTag, err := uc.DBPool.Exec(uc.Ctx, query, userData.Username, userData.FirstName, userData.LastName, userData.Email, userID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %d", userID)
	}
	return nil
}

func (uc *UserController) UserDelete(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	cmdTag, err := uc.DBPool.Exec(uc.Ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Check if the deleting was successful
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user with id %d", userID)
	}
	return nil
}
