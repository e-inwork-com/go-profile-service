package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/e-inwork-com/golang-profile-microservice/internal/data"
	"github.com/e-inwork-com/golang-profile-microservice/internal/jsonlog"

	apiUser "github.com/e-inwork-com/golang-user-microservice/api"
	dataUser "github.com/e-inwork-com/golang-user-microservice/pkg/data"
	jsonLogUser "github.com/e-inwork-com/golang-user-microservice/pkg/jsonlog"

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
		"name text NOT NULL," +
		"email text UNIQUE NOT NULL," +
		"password_hash bytea NOT NULL," +
		"activated bool NOT NULL DEFAULT false," +
		"version integer NOT NULL DEFAULT 1);")
	assert.Nil(t, err)

	_, err = db.Exec("" +
		"CREATE TABLE IF NOT EXISTS profiles (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()," +
		"owner UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE," +
		"first_name text NOT NULL," +
		"last_name text NOT NULL," +
		"version integer NOT NULL DEFAULT 1);")
	assert.Nil(t, err)

	_, err = db.Exec("" +
		"CREATE TABLE IF NOT EXISTS addresses (" +
		"id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()," +
		"created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()," +
		"owner UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE UNIQUE," +
		"street text NOT NULL," +
		"post_code text NOT NULL," +
		"city text NOT NULL," +
		"country_code text NOT NULL," +
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
	user := fmt.Sprintf(`{"name": "Test", "email": "%v", "password": "pa55word"}`, email)
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
	login := `{"email": "test@example.com", "password": "pa55word"}`
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

	// Create Profile with the Authorization JSON Web token from the User Microservice
	firstName := "Jon"
	profile := fmt.Sprintf(`{"first_name": "%v", "last_name": "Doe"}`, firstName)
	req, _ := http.NewRequest("POST", tsProfile.URL+"/api/profiles", bytes.NewReader([]byte(profile)))

	bearer := fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var mProfile map[string]data.Profile
	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].FirstName, firstName)

	// Get Profile with the Authorization JSON Web token from the User Microservice
	req, _ = http.NewRequest("GET", tsProfile.URL+"/api/profile", nil)

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].FirstName, firstName)

	// Patch Profile with the Authorization JSON Web token from the User Microservice
	firstName = "Test"
	profile = fmt.Sprintf(`{"first_name": "%v"}`, firstName)
	req, _ = http.NewRequest("PATCH", tsProfile.URL+"/api/profiles/"+mProfile["profile"].ID.String(), bytes.NewReader([]byte(profile)))

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mProfile)
	assert.Nil(t, err)
	assert.Equal(t, mProfile["profile"].FirstName, firstName)

	// Create Address with the Authorization JSON Web token from the User Microservice
	street := "Kenduruan"
	address := fmt.Sprintf(`{"street": "%v", "post_code": "089", "city": "Tuban", "country_code": "ID"}`, street)
	req, _ = http.NewRequest("POST", tsProfile.URL+"/api/addresses", bytes.NewReader([]byte(address)))

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusAccepted)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	var mAddress map[string]data.Address
	err = json.Unmarshal(body, &mAddress)
	assert.Nil(t, err)
	assert.Equal(t, mAddress["address"].Street, street)

	// Get Address with the Authorization JSON Web token from the User Microservice
	req, _ = http.NewRequest("GET", tsProfile.URL+"/api/address", nil)

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mAddress)
	assert.Nil(t, err)
	assert.Equal(t, mAddress["address"].Street, street)

	// Patch Address with the Authorization JSON Web token from the User Microservice
	street = "Test"
	address = fmt.Sprintf(`{"street": "%v"}`, street)
	req, _ = http.NewRequest("PATCH", tsProfile.URL+"/api/addresses/"+mAddress["address"].ID.String(), bytes.NewReader([]byte(address)))

	bearer = fmt.Sprintf("Bearer %v", authResult.Token)
	req.Header.Set("Authorization", bearer)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(body, &mAddress)
	assert.Nil(t, err)
	assert.Equal(t, mAddress["address"].Street, street)
}
