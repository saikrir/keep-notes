package api

import (
	"errors"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/saikrir/keep-notes/internal/env"
	"github.com/saikrir/keep-notes/internal/logger"
)

func validateToken(accessToken string) bool {

	signingKey := env.GetEnvValAsString("SIGNING_KEY")
	var mySigningKey = []byte(signingKey)
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("could not validate auth token")
		}
		return mySigningKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header["Authorization"]
		if len(authHeader) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("No AuthHeader found")
			return
		}

		logger.Info("KEY ", authHeader[0])

		authHeaderParts := strings.Split(authHeader[0], " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("authorization header could not be parsed")
			return
		}

		if !validateToken(authHeaderParts[1]) {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("could not validate incoming token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
