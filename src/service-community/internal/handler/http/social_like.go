package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary      Like a post
// @Description  Likes a post on behalf of the current user
// @Tags         social
// @Param        X-User-ID header string true "UUID of the current user (mock authorization)"
// @Param        id path int true "Post ID"
// @Success      200 "Successfully liked"
// @Failure      400 "Bad request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal server error"
// @Router       /api/v1/posts/{id}/like [post]
func (h *SocialHandler) like(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.service.ToggleLike(r.Context(), postID, userID, true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Unlike a post
// @Description  Removes the current user's like from a post
// @Tags         social
// @Param        X-User-ID header string true "UUID of the current user (mock authorization)"
// @Param        id path int true "Post ID"
// @Success      200 "Successfully unliked"
// @Failure      400 "Bad request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal server error"
// @Router       /api/v1/posts/{id}/like [delete]
func (h *SocialHandler) unlike(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	postID, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	if err := h.service.ToggleLike(r.Context(), postID, userID, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
