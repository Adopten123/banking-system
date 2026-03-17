package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// @Summary      Follow a user
// @Description  Follows another user from the current user's account
// @Tags         social
// @Param        X-User-ID header string true "UUID of the current user (mock authorization)"
// @Param        id path string true "UUID of the user to follow"
// @Success      200 "Successfully followed"
// @Failure      400 "Bad request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal server error"
// @Router       /api/v1/users/{id}/follow [post]
func (h *SocialHandler) follow(w http.ResponseWriter, r *http.Request) {
	followerID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	followingID, _ := uuid.Parse(chi.URLParam(r, "id"))

	if err := h.service.FollowUser(r.Context(), followerID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary      Unfollow a user
// @Description  Unfollows another user from the current user's account
// @Tags         social
// @Param        X-User-ID header string true "UUID of the current user (mock authorization)"
// @Param        id path string true "UUID of the user to unfollow"
// @Success      200 "Successfully unfollowed"
// @Failure      400 "Bad request"
// @Failure      401 "Unauthorized"
// @Failure      500 "Internal server error"
// @Router       /api/v1/users/{id}/follow [delete]
func (h *SocialHandler) unfollow(w http.ResponseWriter, r *http.Request) {
	followerID, err := getUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	followingID, _ := uuid.Parse(chi.URLParam(r, "id"))

	if err := h.service.UnfollowUser(r.Context(), followerID, followingID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
