package Prototype

import (
	"testing"
	"fmt"
	"github.com/Ready-Stock/Noah/Prototype/cluster"
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
			"DROP TABLE IF EXISTS public.asd3;",
			"CREATE TABLE public.asd3 (test_id BIGINT PRIMARY KEY NOT NULL, testint INT, message VARCHAR);",
			"COMMENT ON TABLE public.asd3 IS 'GLOBAL';",
			"COMMIT;",
		},
	}
)

func Test_CreateTable(t *testing.T) {
	fmt.Printf("%d Table(s) in metadata\n", len(cluster.Tables))
	for _, QuerySet := range createQueries {
		context := Start()
		for _, Query := range QuerySet {
			if _, err := InjestQuery(&context, Query); err != nil {
				fmt.Printf("ERROR ON CLUSTER: %s\n", err.Error())
				if _, err := InjestQuery(&context, "ROLLBACK;"); err != nil {
					t.Error(err)
					t.Fail()
				}
				t.Error(err)
				t.Fail()
				break
			}
		}
		fmt.Println("")
	}
	fmt.Printf("%d Table(s) in metadata\n", len(cluster.Tables))
}
