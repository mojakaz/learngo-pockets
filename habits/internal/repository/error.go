package repository

import "fmt"

// DBError is returned when data is not present in DB.
type DBError struct {
	reason string
}

// Error implements error.
func (e DBError) Error() string {
	return fmt.Sprintf("database error: %s", e.reason)
}
