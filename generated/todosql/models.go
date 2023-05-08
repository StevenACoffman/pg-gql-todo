// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package todosql

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID             uuid.UUID `db:"id" json:"id"`
	Description    string    `db:"description" json:"description"`
	Done           bool      `db:"done" json:"done"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	LastModifiedAt time.Time `db:"last_modified_at" json:"last_modified_at"`
}