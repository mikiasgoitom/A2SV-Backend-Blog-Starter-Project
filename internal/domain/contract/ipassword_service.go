package contract

type IHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswordHash(password, hash string) error
	HashString(s string) string
	CheckHash(s, hash string) bool
}
