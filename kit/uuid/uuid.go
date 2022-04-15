package uuid

import (
	"github.com/google/uuid"
)

// New uuid
func New() string {
	return uuid.New().String()
}
