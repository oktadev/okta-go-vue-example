package http

import (
	"context"
	"log"
	"net/http"
	"strings"

	jwtverifier "github.com/okta/okta-jwt-verifier-golang"
	"github.com/rs/cors"
)

func JSONApi(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}

func AccsessLog(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s: %s", r.Method, r.RequestURI)
		h.ServeHTTP(w, r)
	})
}

func OktaAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header["Authorization"]
		jwt, _ := validateAccessToken(accessToken)
		ctx := context.WithValue(r.Context(), "userId", jwt.Claims["sub"].(string))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateAccessToken(accessToken []string) (*jwtverifier.Jwt, error) {
	parts := strings.Split(accessToken[0], " ")
	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           "https://dev-509836.oktapreview.com/oauth2/default",
		ClaimsToValidate: map[string]string{"aud": "api://default", "cid": "0oagcbm1o6GTTB9Da0h7"},
	}
	verifier := jwtVerifierSetup.New()
	return verifier.VerifyIdToken(parts[1])
}

func UseMiddlewares(h http.Handler) http.Handler {
	h = JSONApi(h)
	h = OktaAuth(h)
	corsConfig := cors.New(cors.Options{
		AllowedHeaders: []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "Authorization"},
		AllowedMethods: []string{"POST", "PUT", "GET", "PATCH", "OPTIONS", "HEAD", "DELETE"},
		Debug:          true,
	})
	h = corsConfig.Handler(h)
	return AccsessLog(h)
}
