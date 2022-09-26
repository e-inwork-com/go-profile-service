package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/e-inwork-com/go-profile-service/internal/data"
	"github.com/e-inwork-com/go-profile-service/internal/jsonlog"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	apiUser "github.com/e-inwork-com/go-user-service/api"
	dataUser "github.com/e-inwork-com/go-user-service/pkg/data"
	jsonLogUser "github.com/e-inwork-com/go-user-service/pkg/jsonlog"

	"github.com/cockroachdb/cockroach-go/v2/testserver"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	// CockroachDB Test Server Setup
	tsDB, err := testserver.NewTestServer()
	assert.Nil(t, err)
	urlDB := tsDB.PGURL()

	// User Microservice
	var cfgUser apiUser.Config
	cfgUser.Db.Dsn = urlDB.String()
	cfgUser.Auth.Secret = "secret"
	cfgUser.Db.MaxOpenConn = 25
	cfgUser.Db.MaxIdleConn = 25
	cfgUser.Db.MaxIdleTime = "15m"
	cfgUser.Limiter.Enabled = true
	cfgUser.Limiter.Rps = 2
	cfgUser.Limiter.Burst = 4

	loggerUser := jsonLogUser.New(os.Stdout, jsonLogUser.LevelInfo)

	db, err := apiUser.OpenDB(cfgUser)
	assert.Nil(t, err)
	defer db.Close()

	_, err = db.Exec("" +
		"CREATE TABLE IF NOT EXISTS users (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()," +
		"email text UNIQUE NOT NULL," +
		"password_hash bytea NOT NULL," +
		"first_name char varying(100) NOT NULL," +
		"last_name char varying(100) NOT NULL," +
		"activated bool NOT NULL DEFAULT false," +
		"version integer NOT NULL DEFAULT 1);")
	assert.Nil(t, err)

	_, err = db.Exec("" +
		"CREATE TABLE IF NOT EXISTS profiles (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()," +
		"profile_user UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE," +
		"profile_picture char varying(512) NOT NULL," +
		"version integer NOT NULL DEFAULT 1);")
	assert.Nil(t, err)

	appUser := &apiUser.Application{
		Config: cfgUser,
		Logger: loggerUser,
		Models: dataUser.InitModels(db),
	}

	tsUser := httptest.NewTLSServer(appUser.Routes())
	defer tsUser.Close()

	// Register an user
	email := "test@example.com"
	password := "pa55word"
	user := fmt.Sprintf(`{"email": "%v", "password":  "%v", "first_name": "Jon", "last_name": "Doe"}`, email, password)
	res, err := tsUser.Client().Post(tsUser.URL+"/api/users", "application/json", bytes.NewReader([]byte(user)))
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var mUser map[string]dataUser.User
	err = json.Unmarshal(body, &mUser)
	assert.Nil(t, err)
	assert.Equal(t, mUser["user"].Email, email)

	// User sign in to get a token
	login := fmt.Sprintf(`{"email": "%v", "password":  "%v"}`, email, password)
	res, err = tsUser.Client().Post(tsUser.URL+"/api/authentication", "application/json", bytes.NewReader([]byte(login)))
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	type authType struct {
		Token string `json:"token"`
	}
	var authResult authType
	err = json.Unmarshal(body, &authResult)

	assert.Nil(t, err)
	assert.NotNil(t, authResult.Token)

	// Profile Microservice
	var cfgProfile Config
	cfgProfile.Db.Dsn = urlDB.String()
	cfgProfile.Auth.Secret = "secret"
	cfgProfile.Db.MaxOpenConn = 25
	cfgProfile.Db.MaxIdleConn = 25
	cfgProfile.Db.MaxIdleTime = "15m"
	cfgProfile.Limiter.Enabled = true
	cfgProfile.Limiter.Rps = 2
	cfgProfile.Limiter.Burst = 6

	loggerProfile := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	appProfile := Application{
		Config: cfgProfile,
		Logger: loggerProfile,
		Models: data.InitModels(db),
	}

	tsProfile := httptest.NewTLSServer(appProfile.routes())
	defer tsProfile.Close()

	// Upload file
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// this step is very important
	filename := "./test/profile.jpg"
	fileWriter, err := bodyWriter.CreateFormFile("profile_picture", filename)
	if err != nil {
		fmt.Println("Error writing to buffer")
	}

	// open file handle
	fileHandler, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
	}

	// Copy file to file handler
	_, err = io.Copy(fileWriter, fileHandler)
	if err != nil {
		fmt.Println("Error copy file")
	}

	// Put on body
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create a new profile
	req, _ := http.NewRequest("POST", tsProfile.URL+"/api/profiles", bodyBuf)
	req.Header.Add("Content-Type", contentType)

	bearer := fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	// Read response
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var mProfile map[string]data.Profile
	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].ProfileUser, mUser["user"].ID)

	// Get a profile of the current user
	req, _ = http.NewRequest("GET", tsProfile.URL+"/api/profiles/me", nil)

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// Read response
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].ProfileUser, mUser["user"].ID)

	// this step is very important
	filename = "./test/profile.png"
	fileWriter, err = bodyWriter.CreateFormFile("profile_picture", filename)
	if err != nil {
		fmt.Println("Error writing to buffer")
	}

	// open file handle
	fileHandler, err = os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file")
	}

	// Copy file to file handler
	_, err = io.Copy(fileWriter, fileHandler)
	if err != nil {
		fmt.Println("Error copy file")
	}

	// Put on body
	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Patch user profile
	req, _ = http.NewRequest("PATCH", tsProfile.URL+"/api/profiles/"+mProfile["profile"].ID.String(), bodyBuf)
	req.Header.Add("Content-Type", contentType)

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// Read response
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].ProfilePicture, mUser["user"].ID.String() + filepath.Ext(filename))
}
