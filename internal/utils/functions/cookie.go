package functions

import (
	"errors"
	"net/http"

	ent "github.com/cantylv/authorization-service/internal/entity"
	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
)

// SetCookie sets cookies to the http response
func SetCookie(w http.ResponseWriter, r *http.Request, token *ent.Session) {
	cookieRefreshToken := http.Cookie{
		Name:     mc.RefreshToken,
		Value:    token.RefreshToken,
		Expires:  token.ExpiresAt,
		HttpOnly: true,
		Path:     "/api/auth",
	}
	http.SetCookie(w, &cookieRefreshToken)
}

func CookieExpired(w http.ResponseWriter, r *http.Request) {
	cookieRefreshToken := http.Cookie{
		Name:   mc.RefreshToken,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &cookieRefreshToken)
}

func IsCookieExist(r *http.Request, cookieName string) (bool, error) {
	_, err := r.Cookie(cookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
