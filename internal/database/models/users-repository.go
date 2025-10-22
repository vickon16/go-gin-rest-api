package models

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Insert("users").
		Columns("email", "name", "password").
		Values(user.Email, user.Name, user.Password).
		Suffix("RETURNING id, email, name, created_at").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		log.Printf("Error converting to sql: %v", err)
		return err
	}

	// Scan the returned row
	return m.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
}

func (m *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("id, email, name, created_at").
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

		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *UserModel) Get(id int64) (*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("id, email, name, created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// For other errors
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) GetUserByEmail(email string, withPassword ...bool) (*User, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	includePassword := false
	if len(withPassword) > 0 && withPassword[0] {
		includePassword = true
	}

	// Define base columns
	columns := []string{"id", "email", "name", "created_at"}
	if includePassword {
		columns = append(columns, "password")
	}

	query := sq.Select(columns...).
		From("users").
		Where(sq.Eq{"email": email}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	if includePassword {
		err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
			&user.ID, &user.Email, &user.Name,
			&user.CreatedAt, &user.Password,
		)
	} else {
		err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
			&user.ID, &user.Email, &user.Name, &user.CreatedAt,
		)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// For other errors
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) Update(id int64, user *UpdateUserDto) (*User, error) {
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

	query = query.Where(sq.Eq{"id": id}).Suffix("RETURNING id, email, name, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var updated User
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
		&updated.ID,
		&updated.Email,
		&updated.Name,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (m *UserModel) Delete(id int64) error {
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
