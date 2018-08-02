package Prototype

import (
	"testing"
)

var (
	InsertQueries = [][]string {
		{
			"BEGIN;",
			"INSERT INTO products (account_id,sku,title) VALUES(1,'test','test');",
			"ROLLBACK;",
		},
	}
)

func Test_Inserts(t *testing.T) {
	for _, QuerySet := range InsertQueries {
		context := Start()
		for _, Query := range QuerySet {
			if err := InjestQuery(&context, Query); err != nil {
				t.Error(err)
				t.Fail()
			}
		}
	}
}
