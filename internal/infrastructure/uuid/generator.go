package uuid

import (
	"github.com/google/uuid"
	"github.com/mikiasgoitom/A2SV-Backend-Blog-Starter-Project/internal/usecase"
)

// Generator implements the usecase.UUIDGenerator interface.
type Generator struct{}

// NewGenerator creates a new UUID generator.
func NewGenerator() usecase.UUIDGenerator {
	return &Generator{}
}

// NewUUID generates a new UUID.
func (g *Generator) NewUUID() uuid.UUID {
	return uuid.New()
}

// Ensure Generator implements the usecase.UUIDGenerator interface
var _ usecase.UUIDGenerator = (*Generator)(nil)
