package integration

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/database"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/router"
	"github.com/nubank/pismo-code-assessment/internal/usecase/account"
)

type TestServer struct {
	Server *httptest.Server
	DB     *sql.DB
}

func (ts *TestServer) Close() {
	ts.Server.Close()
	ts.DB.Close()
}

func SetupTestServer(t *testing.T, ctx context.Context) *TestServer {
	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("pismo_test"),
		postgres.WithUsername("pismo"),
		postgres.WithPassword("pismo"),
		postgres.WithInitScripts(filepath.Join("..", "..", "internal", "infrastructure", "database", "migrations", "001_create_accounts.up.sql")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		postgresContainer.Terminate(ctx)
	})

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	accountRepo := database.NewAccountRepository(db)
	createAccount := account.NewCreateAccount(accountRepo)
	getAccount := account.NewGetAccount(accountRepo)
	accountHandler := handler.NewAccountHandler(createAccount, getAccount)
	r := router.New(accountHandler)

	server := httptest.NewServer(r)
	t.Cleanup(func() {
		server.Close()
		db.Close()
	})

	return &TestServer{
		Server: server,
		DB:     db,
	}
}
