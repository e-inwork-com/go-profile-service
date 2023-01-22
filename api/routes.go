package api

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/service/profiles/health", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/service/profiles", app.requireAuthenticated(app.createProfileHandler))
	router.HandlerFunc(http.MethodGet, "/service/profiles/me", app.requireAuthenticated(app.getProfileHandler))
	router.HandlerFunc(http.MethodPatch, "/service/profiles/:id", app.requireAuthenticated(app.patchProfileHandler))
	router.HandlerFunc(http.MethodGet, "/service/profiles/pictures/:file", app.getProfilePictureHandler)

	router.Handler(http.MethodGet, "/service/profiles/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
