package core

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/bot"
	"github.com/rexagod/newman/core/queries"
	"github.com/rexagod/newman/internal"
	"k8s.io/klog/v2"
)

type Bot struct {
	Ctx *bot.Context
}

type runner struct {
	database        *sql.DB
	databaseContext context.Context
	loader          *internal.Loader
	prefix          string
	token           string
	server          string
	userId          string
	password        string
	databaseName    string
}

var (
	R           = &runner{}
	databaseCtx = context.Background()
)

func initializeDatabase() (*sql.DB, error) {
	db, err := sql.Open("mssql",
		fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s", R.server, R.userId, R.password, R.databaseName))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create pre-requisites.
	_, err = db.ExecContext(databaseCtx, queries.Q[queries.CREATENOTESTABLE])
	if err != nil {
		klog.Infof("failed to create table: %w", err)
	}
	_, err = db.ExecContext(databaseCtx, queries.Q[queries.CREATEDELETEDMESSAGESTABLE])
	if err != nil {
		klog.Infof("failed to create table: %w", err)
	}

	// Ping database.
	if err := db.PingContext(databaseCtx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return db, nil
}

func Start(l *internal.Loader) (*state.State, error) {
	var err error

	// Initialize the runner object.
	R.loader = l
	R.prefix = l.PublicFields["prefix"]
	R.token = l.PrivateFields["token"]
	R.server = l.PrivateFields["server"]
	R.userId = l.PrivateFields["user_id"]
	R.password = l.PrivateFields["password"]
	R.databaseName = l.PrivateFields["database"]
	R.database, err = initializeDatabase()
	R.databaseContext = databaseCtx
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	var s *state.State

	// Start the bot.
	bot.Run(R.token, &Bot{},
		func(ctx *bot.Context) error {

			// Set prefix.
			ctx.HasPrefix = bot.NewPrefix(R.prefix)

			// Initiate handlers.
			s = undoDelete()
			if err := s.Open(context.Background()); err != nil {
				klog.Fatalf("failed to open state: %v", err)
			}
			u, err := s.Me()
			if err != nil {
				klog.Fatalf("failed to get user: %v", err)
			}
			klog.Infof("Starting %s.", u.Username)

			return nil
		},
	)
	return s, nil
}
