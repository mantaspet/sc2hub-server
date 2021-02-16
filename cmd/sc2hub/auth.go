package main

import (
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/mantaspet/sc2hub-server/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var signingKey []byte

func (app *application) generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = username
	claims["exp"] = time.Now().Add(time.Hour * 360).Unix()

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

var isAuthenticated = func(app *application, endpoint func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("token is invalid")
				}
				return signingKey, nil
			})

			if err != nil {
				app.clientError(w, http.StatusUnauthorized, err)
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			app.clientError(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		}
	})
}

func (app *application) getAccessToken(w http.ResponseWriter, r *http.Request) {
	var req models.User

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		app.clientError(w, http.StatusBadRequest, errors.New("specify a valid username and password"))
		return
	}

	errs := req.Validate(app.db)
	if errs != nil {
		app.validationError(w, errs)
		return
	}

	user, err := app.users.SelectOne(req.Username)
	if err == models.ErrNotFound {
		errs = map[string]string{"Username": "Username not found"}
		app.validationError(w, errs)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		errs = map[string]string{"Password": "Incorrect password"}
		app.validationError(w, errs)
		return
	}

	token, err := app.generateJWT(user.Username)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.json(w, token)
}
