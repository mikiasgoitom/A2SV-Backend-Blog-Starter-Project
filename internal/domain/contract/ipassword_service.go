package contract

type IHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswordHash(password, hash string) error
}
