package api

import (
	"io"
	"net/http"
	"testing"

	"github.com/e-inwork-com/go-profile-service/internal/data/mocks"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	app := testApplication(t)

	ts := testServer(t, app.Routes())
	defer ts.Close()

	firstToken := app.testFirstToken(t)
	secondToken := app.testSecondToken(t)
	tBody, tContentType := app.testFormProfile(t)

	tests := []struct {
		name         string
		method       string
		urlPath      string
		contentType  string
		token        string
		body         io.Reader
		expectedCode int
	}{
		{
			name:         "Create Profile",
			method:       "POST",
			urlPath:      "/service/profiles",
			contentType:  tContentType,
			token:        firstToken,
			body:         tBody,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "Get Profile",
			method:       "GET",
			urlPath:      "/service/profiles/me",
			contentType:  "",
			token:        firstToken,
			body:         nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Get Profile Picture",
			method:       "GET",
			urlPath:      "/service/profiles/pictures/" + mocks.MockFirstUUID().String() + ".jpg",
			contentType:  "",
			token:        "",
			body:         nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Patch Profile",
			method:       "PATCH",
			urlPath:      "/service/profiles/" + mocks.MockFirstUUID().String(),
			contentType:  tContentType,
			token:        firstToken,
			body:         tBody,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Patch Profile Forbidden",
			method:       "PATCH",
			urlPath:      "/service/profiles/" + mocks.MockFirstUUID().String(),
			contentType:  tContentType,
			token:        secondToken,
			body:         tBody,
			expectedCode: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualCode, _, _ := ts.request(t, tt.method, tt.urlPath, tt.contentType, tt.token, tt.body)
			assert.Equal(t, tt.expectedCode, actualCode)
		})
	}

}
