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
	Owner		uuid.UUID	`json:"owner"`
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
        INSERT INTO profiles (owner, first_name, last_name) 
        VALUES ($1, $2, $3)
        RETURNING id, created_at, version`

	args := []interface{}{profile.Owner, profile.FirstName, profile.LastName}

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

// GetByOwner Function to get a Profile by Owner
func (m ProfileModel) GetByOwner(owner uuid.UUID) (*Profile, error) {
	// Select query by owner
	query := `
        SELECT id, created_at, owner, first_name, last_name, version
        FROM profiles
        WHERE owner = $1`

	// Define a Profile variable
	var profile Profile

	// Create a context background
	// to use it with a query to database
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query Profile by owner to the database,
	// and the assign the row result to the profile variable
	err := m.DB.QueryRowContext(ctx, query, owner).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.Owner,
		&profile.FirstName,
		&profile.LastName,
		&profile.Version,
	)

	// Check error
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Return the result
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