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
