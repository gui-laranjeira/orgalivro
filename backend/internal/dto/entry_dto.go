package dto

type CreateEntryRequest struct {
	BookID uint   `json:"book_id" binding:"required"`
	Status string `json:"status"`
	Rating *int   `json:"rating"`
	Notes  string `json:"notes"`
}

type UpdateEntryRequest struct {
	Status string `json:"status"`
	Rating *int   `json:"rating"`
	Notes  string `json:"notes"`
}

type EntryListQuery struct {
	Status string `form:"status"`
	Q      string `form:"q"`
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=20"`
}

type EntryResponse struct {
	ID        uint         `json:"id"`
	ProfileID uint         `json:"profile_id"`
	Book      BookResponse `json:"book"`
	Status    string       `json:"status"`
	Rating    *int         `json:"rating"`
	Notes     string       `json:"notes"`
	AddedAt   string       `json:"added_at"`
	UpdatedAt string       `json:"updated_at"`
}

type PaginatedEntries struct {
	Data  []EntryResponse `json:"data"`
	Total int64           `json:"total"`
	Page  int             `json:"page"`
	Limit int             `json:"limit"`
}
