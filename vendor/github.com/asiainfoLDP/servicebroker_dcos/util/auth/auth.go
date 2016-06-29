package auth

import (
	"net/http"
)

type Wrapper struct {
	username string
	password string
}

func NewWrapper(username, password string) *Wrapper {
	return &Wrapper{
		username: username,
		password: password,
	}
}

func (wrapper *Wrapper) Wrap(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !wrapper.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (wrapper *Wrapper) WrapFunc(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !wrapper.isAuthorized(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handlerFunc(w, r)
	})
}

func (wrapper *Wrapper) isAuthorized(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	return ok && wrapper.username == username && wrapper.password == password
}
