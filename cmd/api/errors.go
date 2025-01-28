package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "err", err.Error())

	writeErrJSON(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error", "method", r.Method, "path", r.URL.Path, "err", err.Error())

	writeErrJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "err", err.Error())

	writeErrJSON(w, http.StatusNotFound, "not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("conflict error", "method", r.Method, "path", r.URL.Path, "err", err.Error())

	writeErrJSON(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "err", err.Error())

	writeErrJSON(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "err")

	writeErrJSON(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)
	w.Header().Set("Retry-After", retryAfter)

	writeErrJSON(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}

func (app *application) unauthorizedBasicResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "err", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

	writeErrJSON(w, http.StatusUnauthorized, "unauthorized")
}
