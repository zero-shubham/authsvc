package internal_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	driver "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "github.com/zero-shubham/authsvc/db/orm"
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

	go func(ctx context.Context) {
		<-ctx.Done()
		if !TestDbConn.IsClosed() {
			TestDbConn.Close(ctx)
		}
	}(ctx)

}

func RunMigration() error {
	DB_URL := os.Getenv(DbUrlEnv)

	db, err := sql.Open("pgx/v5", DB_URL)
	if err != nil {
		return err
	}

	dr, err := driver.WithInstance(db, &driver.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///authsvc/db/migrations",
		"pgx/v5", dr)
	if err != nil {
		return err
	}
	return m.Up()
}

func CleanUp() error {
	if TestDbConn.IsClosed() {
		return errors.New("DB connection closed")
	}
	return db.New(TestDbConn).DbCleanUp(context.Background())
}

func TestService(t *testing.T) {
	SetupTest(t.Context())

	require.NoError(t, RunMigration())

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

		require.NoError(t, CleanUp())
	})

}
