package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"

	"github.com/google/uuid"
)

type Profile struct {
	ID        	uuid.UUID	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	FirstName 	string    	`json:"first_name"`
	LastName  	string    	`json:"last_name"`
	Version   	int       	`json:"-"`
}

func ValidateProfile(v *validator.Validator, profile *Profile) {
	v.Check(profile.FirstName != "", "first_name", "must be provided")
	v.Check(profile.LastName != "", "last_name", "must be provided")
}

type ProfileModel struct {
	DB *sql.DB
}

func (m ProfileModel) Insert(profile *Profile) error {
	query := `
        INSERT INTO profiles (first_name, last_name) 
        VALUES ($1, $2)
        RETURNING id, created_at, version`

	args := []interface{}{profile.FirstName, profile.LastName}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&profile.ID, &profile.CreatedAt, &profile.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m ProfileModel) GetByID(id uuid.UUID) (*Profile, error) {
	query := `
        SELECT id, created_at, first_name, last_name, version
        FROM profiles
        WHERE id = $1`

	var profile Profile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.FirstName,
		&profile.LastName,
		&profile.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &profile, nil
}

func (m ProfileModel) Update(profile *Profile) error {
	query := `
        UPDATE profiles 
        SET first_name = $1, last_name = $2, version = version + 1
        WHERE id = $3, AND version = $4
        RETURNING version`

	args := []interface{}{
		profile.FirstName,
		profile.LastName,
		profile.ID,
		profile.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&profile.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}