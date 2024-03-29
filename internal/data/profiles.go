package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/e-inwork-com/go-profile-service/internal/validator"

	"github.com/google/uuid"
)

type ProfileModelInterface interface {
	Insert(profile *Profile) error
	GetByID(id uuid.UUID) (*Profile, error)
	GetByProfileUser(profileUser uuid.UUID) (*Profile, error)
	Update(profile *Profile) error
}

type Profile struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at_dt"`
	ProfileUser    uuid.UUID `json:"profile_user_s"`
	ProfileName    string    `json:"profile_name_t"`
	ProfilePicture string    `json:"profile_picture_s"`
	Version        int       `json:"-"`
}

func ValidateProfile(v *validator.Validator, profile *Profile) {
	v.Check(profile.ProfileName != "", "profile_name_t", "must be provided")
}

type ProfileModel struct {
	DB *sql.DB
}

func (m ProfileModel) Insert(profile *Profile) error {
	query := `
        INSERT INTO profiles (profile_user_s, profile_name_t, profile_picture_s)
        VALUES ($1, $2, $3)
        RETURNING id, created_at_dt, version`

	args := []interface{}{profile.ProfileUser, profile.ProfileName, profile.ProfilePicture}

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
        SELECT id, created_at_dt, profile_user_s, profile_name_t, profile_picture_s, version
        FROM profiles
        WHERE id = $1`

	var profile Profile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.ProfileUser,
		&profile.ProfileName,
		&profile.ProfilePicture,
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

// GetByProfileUser Function to get a Profile by Owner
func (m ProfileModel) GetByProfileUser(profileUser uuid.UUID) (*Profile, error) {
	// Select query by owner
	query := `
        SELECT id, created_at_dt, profile_user_s, profile_name_t, profile_picture_s, version
        FROM profiles
        WHERE profile_user_s = $1`

	// Define a Profile variable
	var profile Profile

	// Create a context background
	// to use it with a query to database
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query Profile by owner to the database,
	// and the assign the row result to the profile variable
	err := m.DB.QueryRowContext(ctx, query, profileUser).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.ProfileUser,
		&profile.ProfileName,
		&profile.ProfilePicture,
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

// Update function to update the Profile
func (m ProfileModel) Update(profile *Profile) error {
	// SQL Update
	query := `
        UPDATE profiles
        SET profile_name_t = $1, profile_picture_s = $2, version = version + 1
        WHERE id = $3 AND version = $4
        RETURNING version`

	// Assign arguments
	args := []interface{}{
		profile.ProfileName,
		profile.ProfilePicture,
		profile.ID,
		profile.Version,
	}

	// Create a context of the SQL Update
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Run SQL Update
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
