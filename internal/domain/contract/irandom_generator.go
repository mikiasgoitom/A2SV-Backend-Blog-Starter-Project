package contract

type IRandomGenerator interface {
	GenerateRandomToken(n int) (string, error)
}
