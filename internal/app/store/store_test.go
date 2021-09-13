package store

import (
	"os"
	"testing"
)

var (
	databaseURL string
)
// TestMain func
func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "sqlserver://kursUser:kursPswd@localhost:1400"
	}
	os.Exit(m.Run())
}
