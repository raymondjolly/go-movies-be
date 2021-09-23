package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pascaldekloe/jwt"
	"go-movies-be/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var validUser = models.User{
	Id:       5,
	Email:    "me@here.com",
	Password: "$2a$12$YZmO3zxVXaKGXORRDxMleOD8COPtz85eSfuxB3ulSwfZmQ6uNzmE2",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
	}

	hashedPassword := validUser.Password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
	}

	var claims jwt.Claims
	claims.Subject = fmt.Sprint(validUser.Id)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 24))
	claims.Issuer = "mydomain.com"
	claims.Audiences = []string{"mydomain.com"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
	}
	app.writeJSON(w, http.StatusOK, string(jwtBytes), "response")

}
