package common

import "math/big"

type AppMode string

const (
	DevMode  AppMode = "development"
	ProdMode AppMode = "production"
)

func BigZero() *big.Int {
	return big.NewInt(0)
}
