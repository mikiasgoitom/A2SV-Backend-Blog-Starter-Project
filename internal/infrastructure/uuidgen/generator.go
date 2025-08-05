package uuidgen

import (
	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/domain/contract"
)

// StandardUUIDGenerator implements the usecase.UUIDGenerator interface
type StandardUUIDGenerator struct{}

// NewUUID generates a new UUID using the standard library's uuid package.
func (g *StandardUUIDGenerator) NewUUID() string {
	return uuid.New().String()
}

// Ensure StandardUUIDGenerator implements the usecase.UUIDGenerator interface
var _ contract.IUUIDGenerator = (*StandardUUIDGenerator)(nil)
