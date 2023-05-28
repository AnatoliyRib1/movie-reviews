package client

import "github.com/AnatoliyRib1/movie-reviews/contracts"

func (c *Client) GetUser(userId int) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().
		SetResult(&u).
		Get(c.path("/api/users/%d", userId))

	return &u, err
}

func (c *Client) GetUserByUserName(userName string) (*contracts.User, error) {
	var u contracts.User

	_, err := c.client.R().
		SetResult(&u).
		Get(c.path("/api/users/username/%s", userName))

	return &u, err
}

func (c *Client) UpdateUser(req *contracts.AuthenticatedRequest[*contracts.UpdateUserRequest]) error {
	_, err := c.client.R().
		SetAuthToken(req.AccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(req.Request).
		Put(c.path("/api/users/%d", req.Request.UserId))

	return err
}
