package integration_tests

import (
	"context"
	"fmt"
	"github.com/AnatoliyRib1/movie-reviews/client"
	"github.com/AnatoliyRib1/movie-reviews/contracts"
	"github.com/AnatoliyRib1/movie-reviews/internal/apperrors"
	"github.com/AnatoliyRib1/movie-reviews/internal/config"
	"github.com/AnatoliyRib1/movie-reviews/internal/modules/users"
	"github.com/AnatoliyRib1/movie-reviews/internal/server"
	"github.com/hashicorp/consul/sdk/testutil/retry"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	prepareInfrastructure(t, runServer)
}

func runServer(t *testing.T, pgConnString string) {
	cfg := &config.Config{
		DbUrl: pgConnString,
		Port:  0, // random port
		Jwt: config.JwtConfig{
			Secret:           "secret",
			AccessExpiration: time.Minute * 15,
		},
		Admin: config.AdminConfig{
			Username: "admin",
			Password: "&dm1Npa$$",
			Email:    "admin@mail.com",
		},
		Local:    true,
		LogLevel: "info",
	}

	srv, err := server.New(context.Background(), cfg)
	require.NoError(t, err)
	defer srv.Close()

	go func() {
		if serr := srv.Start(); serr != http.ErrServerClosed {
			require.NoError(t, serr)
		}
	}()

	var port int
	retry.Run(t, func(r *retry.R) {
		port, err = srv.Port()
		if err != nil {
			require.NoError(r, err)
		}
	})

	tests(t, port, cfg)

	err = srv.Shutdown(context.Background())
	require.NoError(t, err)
}

func tests(t *testing.T, port int, cfg *config.Config) {
	addr := fmt.Sprintf("http://localhost:%d", port)
	c := client.New(addr)

	t.Run("users.GetUserByUserName: admin", func(t *testing.T) {
		u, err := c.GetUserByUserName(cfg.Admin.Username)
		require.NoError(t, err)

		require.Equal(t, cfg.Admin.Username, u.Username)
		require.Equal(t, cfg.Admin.Email, u.Email)
		require.Equal(t, users.AdminRole, u.Role)
	})

	t.Run("users.GetUserByUserName: not found", func(t *testing.T) {
		_, err := c.GetUserByUserName("not found")
		requireNotFoundError(t, err, "user", "username", "not found")
	})

	var (
		johnDoe         *contracts.User
		johnDoePass     = standardPassword
		johnDoeId       int
		johnDoeEmail    = "johndoe@example.com"
		johnDoeUserNAme = "johndoe"
	)
	t.Run("auth.Register: success", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johndoe",
			Email:    "johndoe@example.com",
			Password: johnDoePass,
		}
		u, err := c.RegisterUser(req)
		require.NoError(t, err)
		johnDoe = u
		johnDoeId = u.ID

		require.Equal(t, req.Username, u.Username)
		require.Equal(t, req.Email, u.Email)
		require.Equal(t, users.UserRole, u.Role)
	})

	t.Run("auth.Register: notUniq", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "johndoe",
			Email:    "johndoe@example.com",
			Password: johnDoePass,
		}
		_, err := c.RegisterUser(req)
		requireAlreadyExists(t, err, "user email:johndoe already exists")
	})

	t.Run("users.GetUserByUserId", func(t *testing.T) {
		u, err := c.GetUser(johnDoeId)
		require.NoError(t, err)

		require.Equal(t, johnDoeUserNAme, u.Username)
		require.Equal(t, johnDoeEmail, u.Email)
		require.Equal(t, users.UserRole, u.Role)
	})

	t.Run("users.GetUserByUserIdIfIserIsNotInDB", func(t *testing.T) {
		_, err := c.GetUser(johnDoeId * 100)
		requireBadRequestError(t, err, "")

	})

	t.Run("auth.Register: short username", func(t *testing.T) {
		req := &contracts.RegisterUserRequest{
			Username: "joh",
			Email:    "joh@example.com",
			Password: standardPassword,
		}
		_, err := c.RegisterUser(req)
		requireBadRequestError(t, err, "Username")
	})

	var johnDoeToken string
	t.Run("auth.Login: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    johnDoe.Email,
			Password: johnDoePass,
		}
		res, err := c.LoginUser(req)
		require.NoError(t, err)
		require.NotEmpty(t, res.AccessToken)
		johnDoeToken = res.AccessToken
	})

	t.Run("users.UpdateUser: success", func(t *testing.T) {
		bio := "I'm John Doe"
		req := &contracts.UpdateUserRequest{
			UserId: johnDoe.ID,
			Bio:    &bio,
		}
		err := c.UpdateUser(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)
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

	t.Run("users.setRole: success", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    cfg.Admin.Email,
			Password: cfg.Admin.Password,
		}
		res, err := c.LoginUser(req)
		adminToken := res.AccessToken

		req2 := &contracts.SetUserRoleRequest{
			UserId: johnDoe.ID,
			Role:   "editor",
		}
		err = c.SetRole(contracts.NewAuthenticated(req2, adminToken))
		require.NoError(t, err)
	})

	t.Run("users.setRole: unknown role", func(t *testing.T) {
		req := &contracts.LoginUserRequest{
			Email:    cfg.Admin.Email,
			Password: cfg.Admin.Password,
		}
		res, err := c.LoginUser(req)
		adminToken := res.AccessToken

		req2 := &contracts.SetUserRoleRequest{
			UserId: johnDoe.ID,
			Role:   "editors",
		}
		err = c.SetRole(contracts.NewAuthenticated(req2, adminToken))
		requireBadRequestError(t, err, "")
	})

	t.Run("users.DeleteUser: another user", func(t *testing.T) {

		req := &contracts.DeleteUserRequest{
			UserId: johnDoe.ID + 1,
		}
		err := c.DeleteUser(contracts.NewAuthenticated(req, johnDoeToken))
		requireForbiddenError(t, err, "insufficient permissions")
	})

	t.Run("users.DeleteUser: non-authenticated", func(t *testing.T) {

		req := &contracts.DeleteUserRequest{
			UserId: johnDoe.ID,
		}
		err := c.DeleteUser(contracts.NewAuthenticated(req, ""))
		requireUnauthorizedError(t, err, "invalid or missing token")
	})

	t.Run("users.DeleteUser: success", func(t *testing.T) {

		req := &contracts.DeleteUserRequest{
			UserId: johnDoe.ID,
		}
		err := c.DeleteUser(contracts.NewAuthenticated(req, johnDoeToken))
		require.NoError(t, err)
	})
}

const standardPassword = "secuR3P@ss"

func requireNotFoundError(t *testing.T, err error, subject, key string, value any) {
	msg := apperrors.NotFound(subject, key, value).Error()
	requireApiError(t, err, http.StatusNotFound, msg)
}

func requireUnauthorizedError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusUnauthorized, msg)
}

func requireAlreadyExists(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusConflict, msg)
}

func requireForbiddenError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusForbidden, msg)
}

func requireBadRequestError(t *testing.T, err error, msg string) {
	requireApiError(t, err, http.StatusBadRequest, msg)
}

func requireApiError(t *testing.T, err error, statusCode int, msg string) {
	cerr, ok := err.(*client.Error)
	require.True(t, ok, "expected client.Error")
	require.Equal(t, statusCode, cerr.Code)
	require.Contains(t, cerr.Message, msg)
}
