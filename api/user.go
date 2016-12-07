package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/anthonynsimon/parrot/model"
	"github.com/anthonynsimon/parrot/render"
	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	errs := decodeAndValidate(r.Body, &user)
	if errs != nil {
		render.Error(w, http.StatusUnprocessableEntity, errs)
		return
	}

	existingUser, err := store.GetUserByEmail(user.Email)
	if err == nil && existingUser.Email == user.Email {
		handleError(w, ErrAlreadyExists)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(w, err)
		return
	}

	user.Password = string(hashed)

	result, err := store.CreateUser(user)
	if err != nil {
		handleError(w, err)
		return
	}

	// Hide password
	result.Password = ""
	render.JSON(w, http.StatusCreated, result)
}

func getUserID(ctx context.Context) (string, error) {
	v := ctx.Value("userID")
	if v == nil {
		return "", ErrBadRequest
	}
	id, ok := v.(string)
	if id == "" || !ok {
		return "", ErrInternal
	}
	return id, nil
}

func decodeAndValidate(r io.Reader, m model.Validatable) error {
	if err := json.NewDecoder(r).Decode(m); err != nil {
		return ErrBadRequest
	}
	return m.Validate()
}
