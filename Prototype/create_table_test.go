/*
 * Copyright (c) 2018 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package Prototype

import (
	"fmt"
	"github.com/readystock/noah/Prototype/cluster"
	"testing"
)

var (
	createQueries = [][]string{
		// {
		//     "BEGIN;",
		//     "INSERT INTO products (account_id,sku,title) VALUES(1,'test','test');",
		//     "ROLLBACK;",
		// },
		// {
		//     "BEGIN;",
		//     "INSERT INTO products (account_id,sku,title) VALUES(1,'test','test');",
		//     "COMMIT;",
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
	t.Skip("skipping prototype tests")
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
