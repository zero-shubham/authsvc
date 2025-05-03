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
			Email:      "test@user.com",
			Password:   "test123",
			OrgId:      "",
			AppGroupId: testAppGrp.Id,
		})
		assert.NoError(t, err)
		assert.Equal(t, "test@user.com", resp.Email)

	})

}
