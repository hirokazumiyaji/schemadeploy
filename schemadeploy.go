package schemadeploy

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/schemalex/schemalex"
	"github.com/schemalex/schemalex/diff"
)

type Runner struct {
	Dryrun bool
	DSN    string
	Schema string
}

func (r *Runner) Run(ctx context.Context) error {
	db, err := r.DB()
	if err != nil {
		return err
	}

	defer db.Close()

	return r.DeploySchema(ctx, db)
}

func (r *Runner) DB() (*sql.DB, error) {
	return sql.Open("mysql", strings.Replace(r.DSN, "mysql://", "", 1))
}

func (r *Runner) DeploySchema(ctx context.Context, db *sql.DB) error {
	fromSource, err := schemalex.NewSchemaSource(r.DSN)
	if err != nil {
		return err
	}

	toSource, err := schemalex.NewSchemaSource(r.Schema)
	if err != nil {
		return err
	}

	stmts := &bytes.Buffer{}
	p := schemalex.New()
	err = diff.Sources(
		stmts,
		fromSource,
		toSource,
		diff.WithTransaction(true),
		diff.WithParser(p),
	)
	if err != nil {
		return err
	}

	queries := queryListFromString(stmts.String())

	return r.execSql(ctx, db, queries)
}

func (r *Runner) execSql(ctx context.Context, db *sql.DB, queries queryList) error {
	if r.Dryrun {
		return queries.dump(os.Stdout)
	}
	return queries.execute(ctx, db)
}
