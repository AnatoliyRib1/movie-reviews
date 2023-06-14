package client

import "github.com/AnatoliyRib1/movie-reviews/contracts"

func (c *Client) GetStar(starID int) (*contracts.StarDetails, error) {
	var g contracts.StarDetails

	_, err := c.client.R().
		SetResult(&g).
		Get(c.path("/api/stars/%d", starID))

	return &g, err
}

func (c *Client) GetStars(req *contracts.GetStarsRequest) (*contracts.PaginatedResponse[contracts.Star], error) {
	var res contracts.PaginatedResponse[contracts.Star]

	_, err := c.client.R().
		SetResult(&res).
		SetQueryParams(req.ToQueryParams()).
		Get(c.path("/api/stars"))

	return &res, err
}

func (c *Client) CreateStar(req *contracts.AuthenticatedRequest[*contracts.CreateStarRequest]) (*contracts.StarDetails, error) {
	var g *contracts.StarDetails
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetBody(req.Request).
		SetResult(&g).
		Post(c.path("/api/stars"))

	return g, err
}

func (c *Client) UpdateStar(req *contracts.AuthenticatedRequest[*contracts.UpdateStarRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/stars/%d", req.Request.StarID))

	return err
}

func (c *Client) DeleteStar(req *contracts.AuthenticatedRequest[*contracts.DeleteStarRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Delete(c.path("/api/stars/%d", req.Request.StarID))

	return err
}
