package models

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (m *AttendeesModel) Insert(attendee *Attendee) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Insert("attendees").
		Columns("user_id", "event_id").
		Values(attendee.UserID, attendee.EventID).
		Suffix("RETURNING id, user_id, event_id, created_at").
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	// Scan the returned row
	return m.DB.QueryRowContext(ctx, sqlStr, args...).
		Scan(&attendee.ID, &attendee.UserID, &attendee.EventID, &attendee.CreatedAt)
}

func (m *AttendeesModel) GetAll() ([]*Attendee, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select("id", "user_id", "event_id", "created_at").
		From("attendees").
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

	var attendees []*Attendee

	for rows.Next() {
		var attendee Attendee

		if err := rows.Scan(&attendee.ID, &attendee.UserID, &attendee.EventID, &attendee.CreatedAt); err != nil {
			return nil, err
		}

		attendees = append(attendees, &attendee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

func (m *AttendeesModel) Get(id int) (*Attendee, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"a.id", "a.user_id", "a.event_id", "a.created_at",
		// "u.id", "u.name", "u.email",
	).
		From("attendees a").
		// LeftJoin("users u ON a.user_id = u.id").
		Where(sq.Eq{"a.id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var attendee Attendee
	attendee.User = &User{}

	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&attendee.ID, &attendee.UserID, &attendee.EventID, &attendee.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// For other errors
		return nil, err
	}

	return &attendee, nil
}

func (m *AttendeesModel) Update(id int, attendee *UpdateAttendeeDto) (*Attendee, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Update("attendees").PlaceholderFormat(sq.Dollar)

	if attendee.UserID != 0 {
		query = query.Set("user_id", attendee.UserID)
	}
	if attendee.EventID != 0 {
		query = query.Set("event_id", attendee.EventID)
	}

	query = query.Where(sq.Eq{"id": id}).Suffix("RETURNING id, user_id, event_id, created_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var updated Attendee
	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.EventID,
		&updated.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (m *AttendeesModel) Delete(id int) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Delete("attendees").
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

func (m *AttendeesModel) GetAttendeesByEventId(eventId int64) ([]*Attendee, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"a.id", "a.user_id", "a.event_id", "a.created_at",
		"u.id", "u.name", "u.email", "u.created_at",
		"e.id", "e.user_id", "e.name", "e.description", "e.date", "e.location", "e.created_at",
	).
		From("attendees a").
		LeftJoin("users u ON a.user_id = u.id").
		LeftJoin("events e ON a.event_id = e.id").
		Where(sq.Eq{"a.event_id": eventId}).
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

	var attendees []*Attendee

	for rows.Next() {
		var attendee Attendee
		attendee.User = &User{}
		attendee.Event = &Event{}

		if err := rows.Scan(&attendee.ID, &attendee.UserID, &attendee.EventID, &attendee.CreatedAt,
			&attendee.User.ID, &attendee.User.Name, &attendee.User.Email, &attendee.User.CreatedAt,
			&attendee.Event.ID, &attendee.Event.UserID, &attendee.Event.Name, &attendee.Event.Description,
			&attendee.Event.Date, &attendee.Event.Location, &attendee.Event.CreatedAt); err != nil {
			return nil, err
		}

		attendees = append(attendees, &attendee)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return attendees, nil
}

func (m *AttendeesModel) GetByEventAndAttendee(eventId, userId int64) (*Attendee, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	query := sq.Select(
		"a.id", "a.user_id", "a.event_id", "a.created_at",
		// "u.id", "u.name", "u.email",
	).
		From("attendees a").
		// LeftJoin("users u ON a.user_id = u.id").
		Where(sq.Eq{"a.user_id": userId, "a.event_id": eventId}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var attendee Attendee

	err = m.DB.QueryRowContext(ctx, sqlStr, args...).Scan(&attendee.ID, &attendee.UserID, &attendee.EventID, &attendee.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		// For other errors
		return nil, err
	}

	return &attendee, nil
}
