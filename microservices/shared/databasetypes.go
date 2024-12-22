package shared

import (
	"bytes"
	"encoding/json"
)

type DatabaseQuery struct {
	SqlStatement string `json:"sqlStatement"`
	Args         []any  `json:"args"`
}

func ConstructQuery(statement string, args ...any) *DatabaseQuery {
	return &DatabaseQuery{
		SqlStatement: statement,
		Args:         args,
	}
}

func (q *DatabaseQuery) JSON() []byte {
	marshalled, _ := json.Marshal(q)
	return marshalled
}

func (q *DatabaseQuery) JSONReader() *bytes.Reader {
	marshalled, _ := json.Marshal(q)
	return bytes.NewReader(marshalled)
}
