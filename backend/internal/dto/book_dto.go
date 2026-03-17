package dto

type CreateBookRequest struct {
	Title          string   `json:"title" binding:"required"`
	ISBN13         string   `json:"isbn13"`
	CoverURL       string   `json:"cover_url"`
	Description    string   `json:"description"`
	Year           int      `json:"year"`
	Language       string   `json:"language"`
	Authors        []string `json:"authors"`
	Genres         []string `json:"genres"`
	OwnerProfileID uint     `json:"owner_profile_id"`
}

type UpdateBookRequest struct {
	Title       string   `json:"title"`
	ISBN13      string   `json:"isbn13"`
	CoverURL    string   `json:"cover_url"`
	Description string   `json:"description"`
	Year        int      `json:"year"`
	Language    string   `json:"language"`
	Authors     []string `json:"authors"`
	Genres      []string `json:"genres"`
}

type TransferOwnerRequest struct {
	ProfileID uint `json:"profile_id" binding:"required"`
}

type BookListQuery struct {
	Q              string `form:"q"`
	Author         string `form:"author"`
	Genre          string `form:"genre"`
	Year           int    `form:"year"`
	Language       string `form:"language"`
	OwnerProfileID uint   `form:"owner_profile_id"`
	Page           int    `form:"page,default=1"`
	Limit          int    `form:"limit,default=20"`
}

type AuthorResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookReader struct {
	ProfileID   uint   `json:"profile_id"`
	ProfileName string `json:"profile_name"`
	Status      string `json:"status"`
}

type BookResponse struct {
	ID               uint             `json:"id"`
	Title            string           `json:"title"`
	ISBN13           string           `json:"isbn13"`
	CoverURL         string           `json:"cover_url"`
	Description      string           `json:"description"`
	Year             int              `json:"year"`
	Language         string           `json:"language"`
	OwnerProfileID   *uint            `json:"owner_profile_id"`
	OwnerProfileName string           `json:"owner_profile_name"`
	Authors          []AuthorResponse `json:"authors"`
	Genres           []GenreResponse  `json:"genres"`
	Readers          []BookReader     `json:"readers"`
	CreatedAt        string           `json:"created_at"`
	UpdatedAt        string           `json:"updated_at"`
}

type PaginatedBooks struct {
	Data  []BookResponse `json:"data"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
