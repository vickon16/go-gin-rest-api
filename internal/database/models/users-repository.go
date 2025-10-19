package models

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Insert("users").
		Columns("userId", "name", "description", "date", "location").
		Values(user.Email, user.Name, user.Password).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	// Scan the returned row
	return m.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&user.ID, &user.Email, &user.Name)
}

func (m *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("*").
		From("users").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// Use QueryContext for multiple rows
	rows, err := m.DB.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.ID, &user.Email, &user.Name); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("*").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		// For other errors
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetUserByEmail(email string) (*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("*").
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}

		// For other errors
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Update(id int, user *UpdateUserDto) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Update("users").PlaceholderFormat(sq.Dollar)

	if user.Name != "" {
		query = query.Set("name", user.Name)
	}
	if user.Email != "" {
		query = query.Set("email", user.Email)
	}
	if user.Password != "" {
		query = query.Set("password", user.Password)
	}

	query = query.Where(sq.Eq{"id": id})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}

	_, err = m.DB.ExecContext(ctx, sqlStr, args...)
	return err
}

func (m *UserModel) Delete(id int) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = m.DB.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	return nil
}
