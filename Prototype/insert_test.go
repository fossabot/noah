package Prototype

import (
	"testing"
	"fmt"
)

var (
	InsertQueries = [][]string {
		// {
		// 	"BEGIN;",
		// 	"INSERT INTO products (account_id,sku,title) VALUES(1,'test','test');",
		// 	"ROLLBACK;",
		// },
		// {
		// 	"BEGIN;",
		// 	"INSERT INTO products (account_id,sku,title) VALUES(1,'test','test');",
		// 	"COMMIT;",
		// },
		{
			"BEGIN;",
			"INSERT INTO accounts (account_id,account_name) VALUES(NEW_ID(),'Elliots Account');",
			"COMMIT;",
		},
	}
)

func Test_Inserts(t *testing.T) {
	for _, QuerySet := range InsertQueries {
		context := Start()
		for _, Query := range QuerySet {
			if _, err := InjestQuery(&context, Query); err != nil {
				t.Error(err)
				t.Fail()
				break
			}
		}
		fmt.Println("")
	}
}
