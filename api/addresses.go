package api

import (
	"errors"
	"net/http"

	"github.com/e-inwork-com/golang-profile-microservice/internal/data"
	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"
)

func (app *Application) createAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Address input
	var input struct {
		Street 		string    	`json:"street"`
		PostCode  	string    	`json:"post_code"`
		City  		string    	`json:"city"`
		CountryCode	string    	`json:"country_code"`
	}

	// Assign data from HTTP request to the Address inout
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the current user as the owner of the new Address
	owner := app.contextGetUser(r)

	// Set input to the Address record
	address := &data.Address{
		Owner: 			owner.ID,
		Street:			input.Street,
		PostCode:		input.PostCode,
		City:			input.City,
		CountryCode:	input.CountryCode,
	}

	// Validate Address data
	v := validator.New()
	if data.ValidateAddress(v, address); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert data to Address
	err = app.Models.Addresses.Insert(address)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a Address data as response of the HTTP request
	err = app.writeJSON(w, http.StatusAccepted, envelope{"address": address}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Function to get a Address of the current User
func (app *Application) getAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current user as the owner of the new Address
	owner := app.contextGetUser(r)

	// Get address by owner
	address, err := app.Models.Addresses.GetByOwner(owner.ID)

	// Check error
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a request response
	err = app.writeJSON(w, http.StatusOK, envelope{"address": address}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// patchAddressHandler function to update a Address record
func (app *Application) patchAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from the request parameters
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get Address from the database
	address, err := app.Models.Addresses.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Get the current user
	owner := app.contextGetUser(r)

	// Check if the Address has a related to the owner
	// Only the Owner of the Address can update the own address
	if address.Owner != owner.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Address input
	var input struct {
		Street 			*string    	`json:"street"`
		PostCode  		*string    	`json:"post_code"`
		City  			*string    	`json:"city"`
		CountryCode		*string    	`json:"country_code"`
	}

	// Read JSON from input
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Assign input Street if exist
	if input.Street != nil {
		address.Street = *input.Street
	}

	// Assign input PostCode if exist
	if input.PostCode != nil {
		address.PostCode = *input.PostCode
	}

	// Assign input City if exist
	if input.City != nil {
		address.City = *input.City
	}

	// Assign input CountryCode if exist
	if input.CountryCode != nil {
		address.CountryCode = *input.CountryCode
	}

	// Create a Validator
	v := validator.New()
	// Check if the Address is valid
	if data.ValidateAddress(v, address); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the Address
	err = app.Models.Addresses.Update(address)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send back the Address to the request response
	err = app.writeJSON(w, http.StatusOK, envelope{"address": address}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
