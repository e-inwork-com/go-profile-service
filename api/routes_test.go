package api

import (
	"bytes"
	"database/sql"
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
	"time"

	apiUser "github.com/e-inwork-com/go-user-service/api"
	dataUser "github.com/e-inwork-com/go-user-service/pkg/data"
	jsonLogUser "github.com/e-inwork-com/go-user-service/pkg/jsonlog"

	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
)

var db *sql.DB

func TestRoutes(t *testing.T) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_USER=postgres",
			"POSTGRES_DB=postgres",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	// Tell docker to hard kill the container in 120 seconds
	resource.Expire(120)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Server Setup
	var cfgUser apiUser.Config
	cfgUser.Db.Dsn = databaseUrl
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

	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	assert.Nil(t, err)

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
		"profile_name char varying(100) NOT NULL," +
		"profile_picture char varying(512)," +
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

	// Profile Service
	var cfgProfile Config
	cfgProfile.Db.Dsn = databaseUrl
	cfgProfile.Auth.Secret = "secret"
	cfgProfile.Db.MaxOpenConn = 25
	cfgProfile.Db.MaxIdleConn = 25
	cfgProfile.Db.MaxIdleTime = "15m"
	cfgProfile.Limiter.Enabled = true
	cfgProfile.Limiter.Rps = 2
	cfgProfile.Limiter.Burst = 6
	cfgProfile.Uploads = "../uploads"

	loggerProfile := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	appProfile := Application{
		Config: cfgProfile,
		Logger: loggerProfile,
		Models: data.InitModels(db),
	}

	tsProfile := httptest.NewTLSServer(appProfile.routes())
	defer tsProfile.Close()

	// Create body buffer
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// Add profile name
	bodyWriter.WriteField("profile_name", "Jon Doe")

	// Upload file
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
	assert.Equal(t, mProfile["profile"].ProfilePicture, mUser["user"].ID.String()+filepath.Ext(filename))

	// Get profile picture
	req, _ = http.NewRequest("GET", tsProfile.URL+"/api/profiles/pictures/"+mProfile["profile"].ProfilePicture, nil)

	res, err = tsProfile.Client().Do(req)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	// Read response
	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	// Check type of file
	filetype := http.DetectContentType(body)
	assert.NotEqual(t, filetype, "")
}
