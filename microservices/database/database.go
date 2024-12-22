package database

import (
	"errors"

	"github.com/vicentebravocabezas/streamingUIDE/microservices/shared"
)

func ExecuteQueryNoRows(query *shared.DatabaseQuery) (bool, error) {
	result, err := DB().Exec(query.SqlStatement, query.Args...)
	rowsAffected, _ := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, errors.New("no records were affected")
	}

	return true, nil
}

func ExecuteQuery(query *shared.DatabaseQuery) ([]map[string]any, error) {
	var result []map[string]any

	rows, err := DB().Query(query.SqlStatement, query.Args...)
	if err != nil {
		return nil, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		resultRow := make(map[string]any)
		vals := make([]any, len(cols))
		pointers := make([]any, len(cols))

		for i := range vals {
			pointers[i] = &vals[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		for i := range vals {
			resultRow[cols[i]] = vals[i]
		}

		result = append(result, resultRow)
	}

	return result, nil
}
