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

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"testing"
	"time"
)

func Test_Select(t *testing.T) {
	go StartCoordinator()
	// defer StopCoordinator()
	time.Sleep(15 * time.Second)
	connStr := "postgres://postgres:password@localhost:5433/pqgotest?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	res, err := db.Query("SELECT 1")
	if err != nil {
		panic(err)
	}
	for res.Next() {
		resultString := ""
		if err := res.Scan(&resultString); err != nil {
			panic(err)
		}
		fmt.Printf("Result: [%s]\n", resultString)
	}

}

func Test_Prepare(t *testing.T) {
	go StartCoordinator()
	// defer StopCoordinator()
	time.Sleep(15 * time.Second)
	connStr := "postgres://postgres:Spring!2016@localhost:5433/app?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test 1")
	stmt, err := db.Prepare("SELECT $1::citext")
	fmt.Println("test 2")
	if err != nil {
		fmt.Println("test error")
		t.Error(err)
		t.Fail()
		return
	}
	fmt.Println("test 3")
	rows, err := stmt.Query("test")
	for rows.Next() {
		i := ""
		rows.Scan(&i)
		if i != "2" {
			t.Error("result does not equal")
			t.Fail()
		}
		return
	}
	t.Error("no rows returned")
	t.Fail()
}
