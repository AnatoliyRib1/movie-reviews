package client

import "github.com/AnatoliyRib1/movie-reviews/contracts"

func (c *Client) GetGenreById(genreId int) (*contracts.Genre, error) {
	var g contracts.Genre

	_, err := c.client.R().
		SetResult(&g).
		Get(c.path("/api/genres/%d", genreId))

	return &g, err
}

func (c *Client) GetGenres() ([]*contracts.Genre, error) {
	var g []*contracts.Genre

	_, err := c.client.R().
		SetResult(&g).
		Get(c.path("/api/genres"))

	return g, err
}

func (c *Client) CreateGenre(req *contracts.AuthenticatedRequest[*contracts.CreateGenreRequest]) (*contracts.Genre, error) {
	var g *contracts.Genre
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&g).
		Post(c.path("/api/genres"))

	return g, err
}

func (c *Client) UpdateGenre(req *contracts.AuthenticatedRequest[*contracts.UpdateGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/genres/%d", req.Request.GenreId))

	return err
}

func (c *Client) DeleteGenre(req *contracts.AuthenticatedRequest[*contracts.DeleteGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Delete(c.path("/api/genres/%d", req.Request.GenreId))

	return err
}
