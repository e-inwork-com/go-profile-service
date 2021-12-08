package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"

	"github.com/google/uuid"
)

type Address struct {
	ID        	uuid.UUID	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	Owner		uuid.UUID	`json:"owner"`
	Street 		string    	`json:"street"`
	PostCode  	string    	`json:"post_code"`
	City  		string    	`json:"city"`
	CountryCode string    	`json:"country_code"`
	Version   	int       	`json:"-"`
}

func ValidateAddress(v *validator.Validator, address *Address) {
	v.Check(address.Street != "", "street", "must be provided")
	v.Check(address.PostCode != "", "post_code", "must be provided")
	v.Check(address.City != "", "city", "must be provided")
	v.Check(address.CountryCode != "", "country_code", "must be provided")
}

type AddressModel struct {
	DB *sql.DB
}

func (m AddressModel) Insert(address *Address) error {
	query := `
        INSERT INTO addresses (owner, street, post_code, city, country_code) 
        VALUES ($1, $2,$3, $4, $5)
        RETURNING id, created_at, version`

	args := []interface{}{address.Owner, address.Street, address.PostCode, address.City, address.CountryCode}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&address.ID, &address.CreatedAt, &address.Version)
	if err != nil {
		return err
	}

	return nil
}

func (m AddressModel) GetByID(id uuid.UUID) (*Address, error) {
	query := `
        SELECT id, created_at, street, post_code, city, country_code, version
        FROM addresses
        WHERE id = $1`

	var address Address

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&address.ID,
		&address.CreatedAt,
		&address.Street,
		&address.PostCode,
		&address.City,
		&address.CountryCode,
		&address.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &address, nil
}

// GetByOwner Function to get a Address by Owner
func (m AddressModel) GetByOwner(owner uuid.UUID) (*Address, error) {
	// Select query by owner
	query := `
        SELECT id, created_at, owner, street, post_code, city, country_code, version
        FROM addresses
        WHERE owner = $1`

	// Define a Address variable
	var address Address

	// Create a context background
	// to use it with a query to database
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query Address by owner to the database,
	// and the assign the row result to the address variable
	err := m.DB.QueryRowContext(ctx, query, owner).Scan(
		&address.ID,
		&address.CreatedAt,
		&address.Owner,
		&address.Street,
		&address.PostCode,
		&address.City,
		&address.CountryCode,
		&address.Version,
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
	return &address, nil
}

func (m AddressModel) Update(address *Address) error {
	query := `
        UPDATE addresses 
        SET street = $1, post_code = $2, city = $3, country_code = $4, version = version + 1
        WHERE id = $5, AND version = $6
        RETURNING version`

	args := []interface{}{
		address.Street,
		address.PostCode,
		address.City,
		address.CountryCode,
		address.ID,
		address.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&address.Version)
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