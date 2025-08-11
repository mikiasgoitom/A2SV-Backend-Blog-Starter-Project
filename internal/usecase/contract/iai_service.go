package usecasecontract

import "context"

type IAIService interface {
	GenerateContent(ctx context.Context, prompt string) (string, error)
}
