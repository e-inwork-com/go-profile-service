package api

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/health", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/profiles", app.requireAuthenticated(app.createProfileHandler))
	router.HandlerFunc(http.MethodGet, "/api/profile", app.requireAuthenticated(app.getProfileHandler))
	router.HandlerFunc(http.MethodPatch, "/api/profiles/:id", app.requireAuthenticated(app.patchProfileHandler))
	router.HandlerFunc(http.MethodPost, "/api/addresses", app.requireAuthenticated(app.createAddressHandler))
	router.HandlerFunc(http.MethodGet, "/api/address", app.requireAuthenticated(app.getAddressHandler))

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
