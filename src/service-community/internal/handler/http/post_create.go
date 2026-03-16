package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

// CreatePost godoc
// @Summary      Create new post
// @Description  Creating new post in community (news, user-post или invest-idea)
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        request body CreatePostRequest true "Data"
// @Success      201  {object}  domain.Post
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/v1/posts [post]
func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid request body"})
		return
	}

	authorUUID, err := uuid.Parse(req.AuthorID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: "invalid author_id format"})
		return
	}

	input := domain.CreatePostInput{
		AuthorID:           authorUUID,
		TypeID:             req.TypeID,
		Content:            req.Content,
		RelatedAssetTicker: req.RelatedAssetTicker,
		Status:             "published",
	}

	post, err := h.postService.CreatePost(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
