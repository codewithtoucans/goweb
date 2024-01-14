package middleware

import (
	"github.com/codewithtoucans/goweb/context"
	"github.com/codewithtoucans/goweb/controllers"
	"github.com/codewithtoucans/goweb/models"
	"net/http"
)

type UserMiddleWare struct {
	SessionService *models.SessionService
}

func (um UserMiddleWare) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenValue, err := controllers.ReadCookie(r, controllers.CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := um.SessionService.User(tokenValue)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithUser(r.Context(), user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (um UserMiddleWare) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
