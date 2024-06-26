package db

import (
	"context"
)

var (
	DBURI      string
	DBNAME     string
	TestDBNAME string
)

// const (
// 	DBURI      = "mongodb://localhost:27017"
// 	DBNAME     = "hotel-reserv"
// 	TestDBNAME = "hotel-reserv-test"
// )

type Dropper interface {
	Drop(ctx context.Context) error
}
