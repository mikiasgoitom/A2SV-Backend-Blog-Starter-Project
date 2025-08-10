package usecasecontract

import "context"

type IAIUseCase interface {
	GenerateBlogContent(ctx context.Context, keywords string) (string, error)
	SuggestAndModifyContent(ctx context.Context, keywords, blog string) (string, error)
}
