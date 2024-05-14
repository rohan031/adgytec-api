package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/rohan031/adgytec-api/firebase"
	"github.com/rohan031/adgytec-api/helper"
	"github.com/rohan031/adgytec-api/v1/custom"
)

func TokenAuthetication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check for authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			message := "The request lacks an authorization header."
			err := &custom.MalformedRequest{Status: http.StatusUnauthorized, Message: message}
			helper.HandleError(w, err)
			return
		}

		// check for valid header
		authArray := strings.Split(authHeader, " ")
		if len(authArray) != 2 {
			message := "The authentication header provided is invalid."
			err := &custom.MalformedRequest{Status: http.StatusUnauthorized, Message: message}
			helper.HandleError(w, err)
			return
		}

		// check for bearer scheme
		if bearer := authArray[0]; bearer != "Bearer" {
			message := "The authentication scheme provided is invalid."
			err := &custom.MalformedRequest{Status: http.StatusUnauthorized, Message: message}
			helper.HandleError(w, err)
			return
		}

		// verify id token provided
		idToken := authArray[1]
		token, err := firebase.FirebaseClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {

			if auth.IsIDTokenExpired(err) {
				message := "The ID token provided has expired and is no longer valid for authentication."
				err := &custom.MalformedRequest{Status: http.StatusUnauthorized, Message: message}
				helper.HandleError(w, err)
				return
			}

			if auth.IsIDTokenInvalid(err) {
				message := "The provided ID token is invalid and cannot be used for authentication."
				err := &custom.MalformedRequest{Status: http.StatusUnauthorized, Message: message}
				helper.HandleError(w, err)
				return
			}

			log.Fatalf("error verifying ID token: %v\n", err)
			helper.HandleError(w, err)
			return
		}

		// adding values to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, custom.UserID, token.UID)
		ctx = context.WithValue(ctx, custom.UserRole, token.Claims["role"])
		req := r.WithContext(ctx)

		*r = *req

		next.ServeHTTP(w, r)
	})
}

func RoleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting id params for patch or delete req
		idParam := chi.URLParam(r, "id")

		// check if req is for patch or delete: /users/{id}
		if len(idParam) != 0 {
			uid := r.Context().Value(custom.UserID).(string)

			// if user is trying to perform action in their account
			if idParam == uid {
				next.ServeHTTP(w, r)
				return
			}
		}

		// admin and super admin having privilaged rights
		userRole := r.Context().Value(custom.UserRole).(string)
		if userRole == "user" {
			message := "Insufficient privileges to perform requested action."
			err := &custom.MalformedRequest{Status: http.StatusForbidden, Message: message}
			helper.HandleError(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
