package api

import (
	"errors"
	"fmt"
	"github.com/e-inwork-com/go-profile-service/internal/data"
	"github.com/e-inwork-com/go-profile-service/internal/validator"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Function to create a Profile
func (app *Application) createProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Read a file attachment
	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	// Get the current user
	user := app.contextGetUser(r)

	// Set profile picture
	profilePicture := fmt.Sprintf("%s%s",  user.ID.String(), filepath.Ext(fileHeader.Filename))

	// Set Profile
	profile := &data.Profile{
		ProfileUser: 	user.ID,
		ProfilePicture:	profilePicture,
	}

	// Validate Profile
	v := validator.New()
	if data.ValidateProfile(v, profile); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Check type of file
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	// Read a file from the beginning offset
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create an uploading folder if it doesn't
	// already exist
	err = os.MkdirAll("../uploads", os.ModePerm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("../uploads/%s", profilePicture))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
	err = app.writeJSON(w, http.StatusCreated, envelope{"profile": profile}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getProfileHandler function to get a Profile of the current user
func (app *Application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get the current user as the owner of the Profile
	user := app.contextGetUser(r)

	// Get profile by user
	profile, err := app.Models.Profiles.GetByProfileUser(user.ID)

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
	user := app.contextGetUser(r)

	// Check if the Profile has a related to the current user
	// Only the owner of the Profile can update the own profile
	if profile.ProfileUser != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	// Read a file attachment
	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	// Set profile picture
	profilePicture := fmt.Sprintf("%s%s",  user.ID.String(), filepath.Ext(fileHeader.Filename))

	// Set a new Profile
	newProfile := &data.Profile{
		ProfilePicture:	profilePicture,
	}

	// Create a Validator
	v := validator.New()
	// Check if the Profile is valid
	if data.ValidateProfile(v, newProfile); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Check type of file
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" {
		http.Error(w, "Please upload a JPEG or PNG image", http.StatusBadRequest)
		return
	}

	// Read a file from the beginning offset
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create an uploading folder if it doesn't
	// already exist
	err = os.MkdirAll("../uploads", os.ModePerm)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Delete the old profile picture
	err = os.Remove(fmt.Sprintf("../uploads/%s", profile.ProfilePicture))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Create a new file in the uploads directory
	dst, err := os.Create(fmt.Sprintf("../uploads/%s", profilePicture))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer dst.Close()

	// Copy the uploaded file to the filesystem
	// at the specified destination
	_, err = io.Copy(dst, file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Update the old profile picture with a new one
	profile.ProfilePicture = newProfile.ProfilePicture

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
