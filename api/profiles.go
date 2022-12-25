package api

import (
	"errors"
	"fmt"
	"github.com/e-inwork-com/go-profile-service/internal/data"
	"github.com/e-inwork-com/go-profile-service/internal/validator"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

// Function to create a Profile
func (app *Application) createProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get a profile name
	profileName := r.FormValue("profile_name")

	// Read a file attachment
	file, fileHeader, err := r.FormFile("profile_picture")
	if err == nil {
		defer file.Close()
	}

	// Get the current user
	user := app.contextGetUser(r)

	// Set profile picture
	profilePicture := ""
	if file != nil {
		profilePicture = fmt.Sprintf("%s%s",  user.ID.String(), filepath.Ext(fileHeader.Filename))
	}

	// Set Profile
	profile := &data.Profile{
		ProfileUser: 	user.ID,
		ProfileName: 	profileName,
		ProfilePicture:	profilePicture,
	}

	// Validate Profile
	v := validator.New()
	if data.ValidateProfile(v, profile); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Check type of file
	if profilePicture != "" {
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
		err = os.MkdirAll(app.Config.Uploads, os.ModePerm)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Create a new file in the uploads directory
		dst, err := os.Create(fmt.Sprintf("%s/%s", app.Config.Uploads, profilePicture))
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

	// Get a profile name
	profileName := r.FormValue("profile_name")

	// Read a file attachment
	file, fileHeader, err := r.FormFile("profile_picture")
	if err == nil {
		defer file.Close()
	}

	// Set profile picture
	profilePicture := ""
	if file != nil {
		profilePicture = fmt.Sprintf("%s%s",  user.ID.String(), filepath.Ext(fileHeader.Filename))
	}

	// Set a new Profile
	newProfile := &data.Profile{
		ProfileName: profileName,
		ProfilePicture:	profilePicture,
	}

	if profilePicture != "" {
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
		err = os.MkdirAll(app.Config.Uploads, os.ModePerm)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Delete the old profile picture
		if _, err := os.Stat(fmt.Sprintf("%s/%s", app.Config.Uploads, profile.ProfilePicture)); err == nil {
			err = os.Remove(fmt.Sprintf("%s/%s", app.Config.Uploads, profile.ProfilePicture))
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		// Create a new file in the uploads directory
		dst, err := os.Create(fmt.Sprintf("%s/%s", app.Config.Uploads, profilePicture))
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
	}

	// Update the old profile picture with a new one
	if newProfile.ProfileName != "" {
		profile.ProfileName = newProfile.ProfileName
	}

	if newProfile.ProfilePicture != "" {
		profile.ProfilePicture = newProfile.ProfilePicture
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


// getProfilePictureHandler function to get a profile picture
func (app *Application) getProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	// Get file from the request parameters
	file, err := app.readFileParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Read file
	buffer, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", app.Config.Uploads, file))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Check type of file
	filetype := http.DetectContentType(buffer)

	w.Header().Set("Content-Type", filetype)
	w.Write(buffer)
}