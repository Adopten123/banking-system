package http

type CreateChatRequest struct {
	TypeID    int32    `json:"type_id"` // 1 - private, 2 - group
	Title     *string  `json:"title"`
	MemberIDs []string `json:"member_ids"`
}
