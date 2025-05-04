package internal_test

import (
	"context"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zero-shubham/authsvc/config"
	"github.com/zero-shubham/authsvc/internal"
	"github.com/zero-shubham/authsvc/transport/grpc"
)

const (
	DbUrlEnv = "DATABASE_URL"
)

var TestDbConn *pgx.Conn

func SetupTest(ctx context.Context) {
	DB_URL := os.Getenv(DbUrlEnv)

	if TestDbConn != nil {
		return
	}

	cfg, err := pgx.ParseConfig(DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("error parsing db url")
	}
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error creating db conn")
	}

	TestDbConn = conn
}

func CleanUp(ctx context.Context) error {
	if !TestDbConn.IsClosed() {
		err := TestDbConn.Close(ctx)
		if err != nil {
			return err
		}
	}
	return config.MigrateDown()
}

func dbMigrate() error {
	err := config.MigrateUp()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func TestService(t *testing.T) {
	SetupTest(t.Context())
	require.NoError(t, dbMigrate())

	t.Cleanup(func() {
		require.NoError(t, CleanUp(t.Context()))
	})

	const (
		testEmail    = "test@user.com"
		testPassword = "test123"
	)

	testSvc := internal.NewServer(TestDbConn)
	t.Run("create app group and user", func(t *testing.T) {

		testAppGrp, err := testSvc.CreateAppGroup(t.Context(), &grpc.AppGroupRequest{
			OrgId:  uuid.NewString(),
			Name:   "test-app-group",
			Scopes: []string{"*:*"},
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-app-group", testAppGrp.Name)

		resp, err := testSvc.CreateUser(t.Context(), &grpc.UserRequest{
			Email:      testEmail,
			Password:   testPassword,
			OrgId:      "",
			AppGroupId: testAppGrp.Id,
		})
		assert.NoError(t, err)
		assert.Equal(t, testEmail, resp.Email)

	})

	t.Run("get user token", func(t *testing.T) {
		// Test successful token generation
		tokenResp, err := testSvc.UserToken(t.Context(), &grpc.UserTokenRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenResp.AccessToken)

		// Test invalid password
		_, err = testSvc.UserToken(t.Context(), &grpc.UserTokenRequest{
			Email:    testEmail,
			Password: "wrongpassword",
		})
		assert.Error(t, err)

		// Test non-existent user
		_, err = testSvc.UserToken(t.Context(), &grpc.UserTokenRequest{
			Email:    "nonexistent@user.com",
			Password: testPassword,
		})
		assert.Error(t, err)
	})

	t.Run("get app group", func(t *testing.T) {
		// Create a test app group first
		testAppGrp, err := testSvc.CreateAppGroup(t.Context(), &grpc.AppGroupRequest{
			OrgId:  uuid.NewString(),
			Name:   "test-get-app-group",
			Scopes: []string{"test:scope"},
		})
		assert.NoError(t, err)
		assert.Equal(t, "test-get-app-group", testAppGrp.Name)

		// Test getting app group by ID
		byIDResp, err := testSvc.GetAppGroup(t.Context(), &grpc.GetAppGroupRequest{
			Id: testAppGrp.Id,
		})
		assert.NoError(t, err)
		assert.Equal(t, testAppGrp.Id, byIDResp.Id)
		assert.Equal(t, testAppGrp.Name, byIDResp.Name)
		assert.Equal(t, testAppGrp.Scopes, byIDResp.Scopes)

		// Test getting app group by name
		byNameResp, err := testSvc.GetAppGroup(t.Context(), &grpc.GetAppGroupRequest{
			Name: testAppGrp.Name,
		})
		assert.NoError(t, err)
		assert.Equal(t, testAppGrp.Id, byNameResp.Id)
		assert.Equal(t, testAppGrp.Name, byNameResp.Name)
		assert.Equal(t, testAppGrp.Scopes, byNameResp.Scopes)

		// Test getting non-existent app group by ID
		_, err = testSvc.GetAppGroup(t.Context(), &grpc.GetAppGroupRequest{
			Id: uuid.NewString(),
		})
		assert.Error(t, err)

		// Test getting non-existent app group by name
		_, err = testSvc.GetAppGroup(t.Context(), &grpc.GetAppGroupRequest{
			Name: "non-existent-app-group",
		})
		assert.Error(t, err)

		// Test getting app group with neither ID nor name
		_, err = testSvc.GetAppGroup(t.Context(), &grpc.GetAppGroupRequest{})
		assert.Error(t, err)
	})

}
