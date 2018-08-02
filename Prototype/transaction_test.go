package Prototype

import (
	"testing"
	"fmt"
)

func Test_TransactionSimple(t *testing.T) {
	context := Start()
	fmt.Println("STATE:", context.TransactionState)
	if err := InjestQuery(&context,"BEGIN;"); err != nil {
		t.Error(err)
	}
	fmt.Println("STATE:", context.TransactionState)
}