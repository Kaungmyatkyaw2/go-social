package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/Kaungmyatkyaw2/go-social/internal/mailer"
	"github.com/Kaungmyatkyaw2/go-social/internal/store"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,alphanum,min=3,max=30"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type UserWithToken struct {
	*store.User
	Token string `json:"token"`
}

// RegisterUser godoc
//
//	@Summary		Registers a new user
//	@Description	Registers a new user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User credentials"
//	@Success		201		{object}	UserWithToken		"User registered"
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/auth/register [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()

	hash := sha256.Sum256([]byte(plainToken))

	hashToken := hex.EncodeToString(hash[:])

	err := app.store.Users.CreateAndInvite(r.Context(), user, hashToken, app.config.mail.expDuration)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail, store.ErrDuplicateUsername:
			app.badRequestResponse(w, r, err)

		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}

	isProduction := app.config.env == "production"

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontEndURL, plainToken)

	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	err = app.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProduction)

	if err != nil {

		app.logger.Errorw("error sending welcome email", "error", err)

		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.logger.Errorw("error deleting user after failed email send", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}



type CreateUserTokenPayload struct {
	Email string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// createTokenHandler godoc
//
//	@Summary		Create a token
//	@Description	Create a token for a user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User credentials"
//	@Success		200		{string}	string		"User Token"
//	@Failure		400		{object}	error
//	@Failure		401	{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter,r *http.Request) {
	// parse payload credentials 

		var payload CreateUserTokenPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}


	user, err := app.store.Users.GetByEmail(r.Context(),payload.Email)




	if err != nil {
		switch err{
		case store.ErrorNotFound:
			app.unauthorizedResponse(w,r,err)
			return

		default:
			app.internalServerError(w,r,err)
			return 
		}
	}


	claims := jwt.MapClaims{
		"sub" : user.ID,
		"expo" : time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat" :  time.Now().Unix(),
		"nbf" :  time.Now().Unix(),
		"iss" : app.config.auth.token.iss,
		"aud" : app.config.auth.token.iss,
	}


	token, err := app.authenticator.GenerateToken(claims)


	if err != nil {
		app.internalServerError(w,r,err)
		return 
	}


	if err := app.jsonResponse(w,http.StatusCreated,token); err != nil {
		app.internalServerError(w,r,err)
	}
}