package main

import (
	"github.com/jimsmart/schema"

	"github.com/pkg/errors"

	"github.com/sjhitchner/graphql-resolver/internal/config"
	"github.com/sjhitchner/graphql-resolver/lib/db/psql"
)

func GenerateConfigFromDb() error {
	db, err := psql.NewPSQLDBHandler(dbHost, dbName, dbUser, dbPass, dbPort, false)
	if err != nil {
		return err
	}

	tableName, err := schema.Tables(db)
	if err != nil {
		return err
	}

}
