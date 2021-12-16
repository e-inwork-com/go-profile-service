package api

import (
	"errors"
	"net/http"

	"github.com/e-inwork-com/golang-profile-microservice/internal/data"
	"github.com/e-inwork-com/golang-profile-microservice/internal/validator"
)

// Function to create a Profile
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

// Function to get a Profile of the current User
func (app *Application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current user as the owner of the Profile
	owner := app.contextGetUser(r)

	// Get profile by owner
	profile, err := app.Models.Profiles.GetByOwner(owner.ID)

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
	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// patchProfileHandler function to update a Profile record
func (app *Application) patchProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get ID from the request parameters
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get Profile from the database
	profile, err := app.Models.Profiles.GetByID(id)
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

	// Check if the Profile has a related to the owner
	// Only the Owner of the Profile can update the own profile
	if profile.Owner != owner.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Profile input
	var input struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
	}

	// Read JSON from input
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Assign input FirstName if exist
	if input.FirstName != nil {
		profile.FirstName = *input.FirstName
	}

	// Assign input LastName if exist
	if input.LastName != nil {
		profile.LastName = *input.LastName
	}

	// Create a Validator
	v := validator.New()
	// Check if the Profile is valid
	if data.ValidateProfile(v, profile); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Update the Profile
	err = app.Models.Profiles.Update(profile)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Send back the Profile to the request response
	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
