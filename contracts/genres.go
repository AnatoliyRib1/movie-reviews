package contracts

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type GetGenreRequest struct {
	GenreId int `param:"genreId" validate:"nonzero"`
}
type CreateGenreRequest struct {
	Name string `json:"name" validate:"min=3,max=32"`
}
type UpdateGenreRequest struct {
	GenreId int    `param:"genreId" validate:"nonzero"`
	Name    string `json:"name" validate:"min=3,max=32"`
}
type DeleteGenreRequest struct {
	GenreId int `param:"genreId" validate:"nonzero"`
}
