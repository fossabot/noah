package Prototype

import (
	"testing"
	"fmt"
)

var (
	createQueries = [][]string {
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
			"CREATE TABLE public.xyz (test_id BIGINT PRIMARY KEY NOT NULL, testint INT, message VARCHAR);",
			"ROLLBACK;",
		},
	}
)

func Test_CreateTable(t *testing.T) {
	for _, QuerySet := range createQueries {
		context := Start()
		for _, Query := range QuerySet {
			if err := InjestQuery(&context, Query); err != nil {
				t.Error(err)
				t.Fail()
				break
			}
		}
		fmt.Println("")
	}
}
