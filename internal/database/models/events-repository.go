package models

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (m *EventModel) Insert(event *CreateEventDto) (*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Insert("events").
		Columns("user_id", "name", "description", "date", "location").
		Values(event.UserID, event.Name, event.Description, event.Date, event.Location).
		Suffix("RETURNING id, user_id, name, description, date, location, created_at").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// Scan the returned row
	var newEvent Event
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&newEvent.ID, &newEvent.UserID, &newEvent.Name, &newEvent.Description, &newEvent.Date, &newEvent.Location, &newEvent.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &newEvent, nil
}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"e.id", "e.user_id", "e.name", "e.description", "e.date", "e.location", "e.created_at",
		"u.id", "u.name", "u.email",
	).
		From("events e").
		LeftJoin("users u ON e.user_id = u.id").
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
		event.User = &User{}

		if err := rows.Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location, &event.CreatedAt, &event.User.ID, &event.User.Name, &event.User.Email); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Get(id int64) (*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"e.id", "e.user_id", "e.name", "e.description", "e.date", "e.location", "e.created_at",
		"u.id", "u.name", "u.email",
	).
		From("events e").
		LeftJoin("users u ON e.user_id = u.id").
		Where(sq.Eq{"e.id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var event Event
	event.User = &User{}

	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location, &event.CreatedAt, &event.User.ID, &event.User.Name, &event.User.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// For other errors
		return nil, err
	}

	return &event, nil
}

func (m *EventModel) GetEventsByAttendeeId(attendeeId int64) ([]*Event, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"e.id", "e.user_id", "e.name", "e.description", "e.date", "e.location", "e.created_at",
		"u.id", "u.name", "u.email",
	).
		From("events e").
		LeftJoin("users u ON e.user_id = u.id").
		Where(sq.Eq{"e.user_id": attendeeId}).
		OrderBy("e.created_at ASC").
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
		event.User = &User{}

		if err := rows.Scan(&event.ID, &event.UserID, &event.Name, &event.Description, &event.Date, &event.Location, &event.CreatedAt, &event.User.ID, &event.User.Name, &event.User.Email); err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Update(id int64, event *UpdateEventDto) (*Event, error) {
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

	query = query.Where(sq.Eq{"id": id}).Suffix("RETURNING id, user_id, name, description, date, location, created_at, updated_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var updated Event
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Name,
		&updated.Description,
		&updated.Date,
		&updated.Location,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (m *EventModel) Delete(id int64) error {
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
