package models

import (
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Insert("events").
		Columns("userId", "name", "description", "date", "location").
		Values(event.UserID, event.Name, event.Description, event.Date, event.Location).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	// Scan the returned row
	return m.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location)
}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("*").
		From("events").
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

	var events []*Event

	for rows.Next() {
		var event Event

		if err := rows.Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("*").
		From("events").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var event Event
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		// For other errors
		return nil, err
	}

	return &event, nil
}

func (m *EventModel) Update(id int, event *UpdateEventDto) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Update("events").PlaceholderFormat(sq.Dollar)

	if event.Name != "" {
		query = query.Set("name", event.Name)
	}
	if event.Description != "" {
		query = query.Set("description", event.Description)
	}
	if !event.Date.IsZero() {
		query = query.Set("date", event.Date)
	}
	if event.Location != "" {
		query = query.Set("location", event.Location)
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

func (m *EventModel) Delete(id int) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Delete("events").
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
