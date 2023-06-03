package client

import "github.com/AnatoliyRib1/movie-reviews/contracts"

func (c *Client) GetStar(starID int) (*contracts.Star, error) {
	var g contracts.Star

	_, err := c.client.R().
		SetResult(&g).
		Get(c.path("/api/stars/%d", starID))

	return &g, err
}

/*
	func (c *Client) GetGenres() ([]*contracts.Genre, error) {
		var g []*contracts.Genre

		_, err := c.client.R().
			SetResult(&g).
			Get(c.path("/api/genres"))

		return g, err
	}
*/
func (c *Client) CreateStar(req *contracts.AuthenticatedRequest[*contracts.CreateStarRequest]) (*contracts.Star, error) {
	var g *contracts.Star
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&g).
		Post(c.path("/api/stars"))

	return g, err
}

/*
func (c *Client) UpdateGenre(req *contracts.AuthenticatedRequest[*contracts.UpdateGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/genres/%d", req.Request.GenreID))

	return err
}

func (c *Client) DeleteGenre(req *contracts.AuthenticatedRequest[*contracts.DeleteGenreRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Delete(c.path("/api/genres/%d", req.Request.GenreID))

	return err
}
*/
