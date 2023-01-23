package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/e-inwork-com/go-profile-service/internal/data"
	"github.com/e-inwork-com/go-profile-service/internal/jsonlog"
	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	// Cofiguration
	var cfg Config
	cfg.Db.Dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	cfg.Auth.Secret = "secret"
	cfg.Db.MaxOpenConn = 25
	cfg.Db.MaxIdleConn = 25
	cfg.Db.MaxIdleTime = "15m"
	cfg.Limiter.Enabled = true
	cfg.Limiter.Rps = 2
	cfg.Limiter.Burst = 6
	cfg.Uploads = "../local/test/uploads"
	cfg.GRPCProfile = "localhost:5002"

	// Logger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Database
	db, err := OpenDB(cfg)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	// Application
	app := Application{
		Config: cfg,
		Logger: logger,
		Models: data.InitModels(db, cfg.GRPCProfile),
	}

	// API Routes
	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()

	// Read SQL file
	script, err := os.ReadFile("./test/sql/delete_all.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Delete Records
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Initial
	email := "jon@doe.com"
	password := "pa55word"

	// Initial
	var userResponse map[string]data.User

	t.Run("Register User", func(t *testing.T) {
		data := fmt.Sprintf(
			`{"email": "%v", "password": "%v", "first_name": "Jon", "last_name": "Doe"}`,
			email,
			password)
		req, _ := http.NewRequest(
			"POST",
			"http://localhost:8000/service/users",
			bytes.NewReader([]byte(data)))
		req.Header.Add("Content-Type", "application/json")

		res, err := ts.Client().Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &userResponse)
		assert.Nil(t, err)
		assert.Equal(t, email, userResponse["user"].Email)
	})

	// Initial
	type authType struct {
		Token string `json:"token"`
	}
	var authentication authType

	t.Run("Login User", func(t *testing.T) {
		data := fmt.Sprintf(
			`{"email": "%v", "password": "%v"}`,
			email,
			password)
		req, _ := http.NewRequest(
			"POST",
			"http://localhost:8000/service/users/authentication",
			bytes.NewReader([]byte(data)))
		req.Header.Add("Content-Type", "application/json")

		res, err := ts.Client().Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &authentication)
		assert.Nil(t, err)
		assert.NotNil(t, authentication.Token)
	})

	// Initial
	var profileResponse map[string]data.Profile

	t.Run("Create Profile", func(t *testing.T) {
		tBody, tContentType := app.testFormProfile(t)
		req, _ := http.NewRequest(
			"POST",
			ts.URL+"/service/profiles",
			tBody)
		req.Header.Add("Content-Type", tContentType)

		bearer := fmt.Sprintf("Bearer %v", authentication.Token)
		req.Header.Set("Authorization", bearer)

		res, err := ts.Client().Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &profileResponse)
		assert.Nil(t, err)
		assert.Equal(t,
			profileResponse["profile"].ProfileUser,
			userResponse["user"].ID)
	})

	t.Run("Patch Profile", func(t *testing.T) {
		tBody, tContentType := app.testFormProfile(t)
		req, _ := http.NewRequest(
			"PATCH",
			ts.URL+"/service/profiles/"+profileResponse["profile"].ID.String(),
			tBody)
		req.Header.Add("Content-Type", tContentType)

		bearer := fmt.Sprintf("Bearer %v", authentication.Token)
		req.Header.Set("Authorization", bearer)

		res, err := ts.Client().Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &profileResponse)
		assert.Nil(t, err)
		assert.Equal(t,
			profileResponse["profile"].ProfileUser,
			userResponse["user"].ID)
	})
}
