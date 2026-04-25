package main

import (
	"errors"
	"net/http"

	"github.com/Kaungmyatkyaw2/go-social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
	UserID  int64  `json:"user_id" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {

	post := getPostFromCtx(r)

	var payload CreateCommentPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		UserID:  payload.UserID,
		PostID:  post.ID,
	}

	err := app.store.Comments.Create(r.Context(), comment)

	if err != nil {

		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundResponse(w, r, err)

		default:
			app.internalServerError(w, r, err)

		}

		return
	}

	if err := app.jsonResponse(w, http.StatusOK, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
