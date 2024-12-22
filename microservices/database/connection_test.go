package database

import (
	"fmt"
	"testing"
)

func TestUserQuery(t *testing.T) {
	SetDBFile("file:./local/dbtest.sqlite3")

	rows, err := DB().Query("SELECT username, email FROM users LIMIT 10")
	if err != nil {
		t.Fatal(err)
	}

	var users []string

	for rows.Next() {
		var username string
		var email string

		if err := rows.Scan(&username, &email); err != nil {
			t.Fatal(err)
		}

		users = append(users, fmt.Sprintf("%v, %v", username, email))
	}

	if len(users) == 0 {
		t.Log("Query was successful but returned no rows")
		return
	}

	t.Logf("Query returned %v records", len(users))
}
