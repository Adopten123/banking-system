package http

type ErrorResponse struct {
	Message string `json:"message" example:"invalid request body"`
}

type CreatePostRequest struct {
	AuthorID           string  `json:"author_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	TypeID             int32   `json:"type_id" example:"2"`
	Content            string  `json:"content" example:"Привет, это мой первый пост!"`
	RelatedAssetTicker *string `json:"related_asset_ticker,omitempty" example:"AAPL"`
}
