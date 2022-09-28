package core

import (
	"database/sql"
	"fmt"
	"github.com/rexagod/newman/core/queries"
	"strings"
	"time"
)

func addRow(addQuery string, params ...interface{}) (string, error) {
	_, err := R.database.ExecContext(R.databaseContext, addQuery, params...)
	if err != nil {
		return "", fmt.Errorf("failed to execute query: %v", err)
	}
	return "`Entry added!`", nil
}

func showRows(showQuery string) (string, error) {
	rows, err := R.database.QueryContext(R.databaseContext, showQuery)
	if err != nil {
		return "", fmt.Errorf("failed to query database: %v", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)
	var out []string
	outer := []string{"", ""}
	for rows.Next() {
		switch showQuery {
		case queries.Q[queries.SHOWDELETEDMESSAGES]:
			var id int
			var author, timestamp, message string
			err := rows.Scan(&id, &author, &timestamp, &message)
			if err != nil {
				return "", fmt.Errorf("failed to scan row: %v", err)
			}
			formatTimestamp, err := time.Parse(time.RFC3339, timestamp)
			if err != nil {
				return "", fmt.Errorf("failed to parse timestamp: %v", err)
			}
			// TODO: handle non-textual data for message
			s := strings.Split(formatTimestamp.String(), "+")[0]
			out = append(out, fmt.Sprintf("%s: \"%s\" at _%s_.", author, message, s[:len(s)-1]))
		case queries.Q[queries.SHOW]:
			var id int
			var content string
			err := rows.Scan(&id, &content)
			if err != nil {
				return "", fmt.Errorf("failed to scan row: %v", err)
			}
			out = append(out, fmt.Sprintf("%d: %s", id, content))
			outer = []string{"```yaml\n", "\n```"}
		}
	}
	if rows.Err() != nil {
		return "", fmt.Errorf("failed to iterate over rows: %v", rows.Err())
	}
	outs := outer[0] + strings.Join(out, "\n") + outer[1]
	if len(out) == 0 {
		outs = "`No entries found.`"
	}
	return outs, nil
}
