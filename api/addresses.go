package api

import (
	"net/http"

	"github.com/e-inwork-com/golang-profile-microservice/internal/data"
	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"
)

func (app *Application) createAddressHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Street 		string    	`json:"street"`
		PostCode  	string    	`json:"post_code"`
		City  		string    	`json:"city"`
		CountryCode string    	`json:"country_code"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	address := &data.Address{
		Street:		input.Street,
		PostCode:	input.PostCode,
		City:		input.City,
		CountryCode:input.CountryCode,
	}

	v := validator.New()

	if data.ValidateAddress(v, address); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Addresses.Insert(address)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"address": address}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
