package http

import (
	"encoding/json"
	"net/http"

	"github.com/Adopten123/banking-system/service-community/internal/domain"
	"github.com/google/uuid"
)

// @Summary      Создать новый чат
// @Description  Создает чат (приватный или группу) и добавляет в него участников
// @Tags         chats
// @Accept       json
// @Produce      json
// @Param        X-User-ID header string true "UUID текущего пользователя (имитация авторизации)"
// @Param        input body CreateChatRequest true "Данные чата"
// @Success      201 {object} domain.Chat
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/chats [post]
func (h *ChatHandler) createChat(w http.ResponseWriter, r *http.Request) {
	currentUserIDStr := r.Header.Get("X-User-ID")
	currentUserID, err := uuid.Parse(currentUserIDStr)
	if err != nil {
		http.Error(w, "Unauthorized: missing or invalid X-User-ID header", http.StatusUnauthorized)
		return
	}

	var req CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var parsedMemberIDs []uuid.UUID
	for _, idStr := range req.MemberIDs {
		parsedID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid member UUID: "+idStr, http.StatusBadRequest)
			return
		}
		parsedMemberIDs = append(parsedMemberIDs, parsedID)
	}

	isCreatorInList := false
	for _, memberID := range parsedMemberIDs {
		if memberID == currentUserID {
			isCreatorInList = true
			break
		}
	}
	if !isCreatorInList {
		parsedMemberIDs = append(parsedMemberIDs, currentUserID)
	}

	input := domain.CreateChatInput{
		TypeID:    req.TypeID,
		Title:     req.Title,
		MemberIDs: parsedMemberIDs,
	}

	chat, err := h.chatService.CreateChat(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}
