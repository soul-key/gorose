package gorose

import (
	"database/sql"
	"github.com/gohouse/gorose/v3/builder"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func As(table any, alias string) builder.TableClause {
	return builder.TableClause{
		Tables: table,
		Alias:  alias,
	}
}
func GetRandomInt(num int) int {
	return rand.Intn(num)
}

//func GetRandomWeightedIndex(weights []int) int {
//	if len(weights) == 0 {
//		return 0
//	}
//	if len(weights) == 1 {
//		return 0
//	}
//	totalWeight := 0
//	for _, w := range weights {
//		totalWeight += w
//	}
//	if totalWeight == 0 {
//		return rand.Intn(len(weights))
//	}
//
//	rnd := rand.Intn(totalWeight)
//
//	currentWeight := 0
//	for i, w := range weights {
//		currentWeight += w
//		if rnd < currentWeight {
//			return i
//		}
//	}
//	return -1 // 如果权重都为 0，或者总权重为 0，则返回 -1
//}

//////////// struct field ptr 4 orm helpers ////////////

func Ptr[T any](arg T) *T {
	return &arg
}

//////////// sql.Null* type helpers ////////////

func Null[T any](arg T) sql.Null[T] {
	return sql.Null[T]{V: arg, Valid: true}
}
