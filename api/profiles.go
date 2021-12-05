package api

import (
	"net/http"

	"github.com/e-inwork-com/golang-profile-microservice/internal/data"
	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"
)

func (app *Application) createProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Profile input
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
	}

	// Assign data from HTTP request to the Profile inout
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the current user as the owner of the new Profile
	owner := app.contextGetUser(r)

	// Set input to the Profile record
	profile := &data.Profile{
		Owner: 		owner.ID,
		FirstName:	input.FirstName,
		LastName:	input.LastName,
	}

	// Validate Profile data
	v := validator.New()
	if data.ValidateProfile(v, profile); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Insert data to Profile
	err = app.Models.Profiles.Insert(profile)
	if err != nil {
		switch {
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send a Profile data as response of the HTTP request
	err = app.writeJSON(w, http.StatusAccepted, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
