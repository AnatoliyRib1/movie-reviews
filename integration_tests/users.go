package integration_tests

import (
	"net/http"
	"testing"

	"github.com/AnatoliyRib1/movie-reviews/client"
	"github.com/AnatoliyRib1/movie-reviews/contracts"
	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/stretchr/testify/require"
)

func usersApiChecks(t *testing.T, c *client.Client, cfg *config.Config) {
	t.Run("users.GetUserByUserName: admin", func(t *testing.T) {
		u, err := c.GetUserByUserName(cfg.Admin.Username)
		require.NoError(t, err)
		admin = u

		require.Equal(t, cfg.Admin.Username, u.Username)
		require.Equal(t, cfg.Admin.Email, u.Email)
		require.Equal(t, users.AdminRole, u.Role)
	})

	t.Run("users.GetUserByUserName: not found", func(t *testing.T) {
		_, err := c.GetUserByUserName("notfound")
		requireNotFoundError(t, err, "user", "username", "notfound")
	})

	t.Run("users.GetUserById: admin", func(t *testing.T) {
		u, err := c.GetUserById(admin.ID)
		require.NoError(t, err)

		require.Equal(t, admin, u)
	})

	t.Run("users.GetUserById: not found", func(t *testing.T) {
		nonExistingId := 1000
		_, err := c.GetUserById(nonExistingId)
		requireNotFoundError(t, err, "user", "id", nonExistingId)
	})

	t.Run("users.UpdateUser: success", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateUserRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)

		johnDoe = getUser(t, c, johnDoe.ID)
		require.Equal(t, bio, *johnDoe.Bio)
	})

	t.Run("users.UpdateUser: non-authenticated", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateUserRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("users.UpdateUser: another user", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateUserRequest{
			UserId: johnDoe.ID + 1,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.UpdateUser: by admin", func(t *testing.T) {
		bio := "Updated by admin"
		req := &contracts.UpdateUserRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}

		err := c.UpdateUser(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)

		johnDoe = getUser(t, c, johnDoe.ID)
		require.Equal(t, bio, *johnDoe.Bio)
	})

	t.Run("users.SetUserRole: John Doe to editor", func(t *testing.T) {
		req := &contracts.SetUserRoleRequest{
			UserId: johnDoe.ID,
			Role:   users.EditorRole,
		}
		err := c.SetUserRole(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)

		johnDoe = getUser(t, c, johnDoe.ID)
		require.Equal(t, users.EditorRole, johnDoe.Role)

		// Have to re-login to become an editor
		johnDoeToken = login(t, c, johnDoe.Email, johnDoePass)
	})
	t.Run("users.SetUserRole: bad role", func(t *testing.T) {
		req := &contracts.SetUserRoleRequest{
			UserId: johnDoe.ID,
			Role:   "superuser",
		}
		err := c.SetUserRole(contracts.NewAuthenticated(req, adminToken))
		requireBadRequestError(t, err, "Role")
	})

	randomUser := registerRandomUser(t, c)
	t.Run("users.DeleteUser: another user", func(t *testing.T) {
		req := &contracts.DeleteUserRequest{
			UserId: randomUser.ID,
		}
		err := c.DeleteUser(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")

		randomUser = getUser(t, c, randomUser.ID)
		require.NotNil(t, randomUser)
	})

	t.Run("users.DeleteUser: by admin", func(t *testing.T) {
		req := &contracts.DeleteUserRequest{
			UserId: randomUser.ID,
		}
		err := c.DeleteUser(contracts.NewAuthenticated(req, adminToken))
		require.NoError(t, err)

		randomUser = getUser(t, c, randomUser.ID)
		require.Nil(t, randomUser)
	})
}

func getUser(t *testing.T, c *client.Client, id int) *contracts.User {
	u, err := c.GetUserById(id)
	if err != nil {
		cerr, ok := err.(*client.Error)
		require.True(t, ok)
		require.Equal(t, http.StatusNotFound, cerr.Code)
		return nil
	}

	return u
}
