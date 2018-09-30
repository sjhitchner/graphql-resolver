package psqlloader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sjhitchner/graphql-resolver/lib/db/psql"
)

type Config struct {
	Files []struct {
		Database      string `json:"databaes"`
		Table         string `json:"table"`
		Schema        string `json:"schema"`
		Data          string `json:"data"`
		ExpectedCount int    `json:"count"`
	} `json:"files"`
}

func (t Config) String() string {
	b, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func InitializeDB(configPath string) error {

	f, err := os.Open(configPath)
	if err != nil {
		return err
	}

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return err
	}

	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		return err
	}
	if port == 0 {
		port = 15432
	}

	dbh, err := psql.NewPSQLDBHandler(
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		port,
		psql.SSLModeDisable)
	if err != nil {
		return err
	}

	for _, file := range config.Files {
		results, err := dbh.DBConnection().Exec(
			fmt.Sprintf(`SELECT EXISTS(
                          SELECT * 
                          FROM information_schema.tables 
                          WHERE 
                           table_schema = 'public' AND 
                           table_name = '%s')`,
				file.Table))
		if err != nil {
			return err
		}
		if count, err := results.RowsAffected(); err == nil && count > 0 {
			fmt.Println("TABLE " + file.Table + " EXISTS")
			continue
		}

		loadSchema(dbh.DBConnection(), file.Table, file.Schema)
		loadData(dbh.DBConnection(), file.Table, file.Data)
	}

	return nil
}

func loadSchema(conn *sqlx.DB, table, filepath string) error {
	conn.Exec("DROP TABLE " + table + " IF EXISTS")

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	cmds := make([]string, 0, 100)

	buf := &bytes.Buffer{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "--") {
			continue
		}

		buf.WriteString(text)

		if strings.HasSuffix(text, ";") {
			fmt.Println(buf.String())
			cmds = append(cmds, buf.String())
			buf.Reset()
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	for _, cmd := range cmds {
		_, err = conn.Exec(cmd)
		fmt.Println(err)
	}

	return err
}

func loadData(conn *sqlx.DB, table, file string) error {

	stmt := fmt.Sprintf("COPY %s FROM '%s';", table, file)

	fmt.Println(stmt)

	_, err := conn.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}
