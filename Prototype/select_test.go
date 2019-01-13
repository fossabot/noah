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
	"testing"
)

var (
	SelectQueries = []string{
		"SELECT 1",
		"SELECT 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id=1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = '1' OR product_id = '2';",
		"SELECT products.product_id FROM products WHERE (account_id = '1') AND (product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id FROM products WHERE (products.account_id = '1') AND (products.product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id WHERE (account_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (account_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (product.id = 1 or product.test = true) and (products.account_id = '1') LIMIT 10 OFFSET 0;",
	}
	FunctionQueries = []string{
		"SELECT do_test('test');",
		"SELECT CURRENT_TIMESTAMP;",
	}
)

func Test_Select1(t *testing.T) {
	t.Skip("skipping prototype tests")
	context := Start()
	if rows, err := InjestQuery(&context, "SELECT 1;"); err != nil {
		t.Error(err)
	} else {
		for rows.Next() {
			var col = ""
			rows.Scan(&col)
			fmt.Printf("Row (%s)", col)
		}
	}
}

func Test_Selects(t *testing.T) {
	t.Skip("skipping prototype tests")
	for _, q := range SelectQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}

func Test_SelectFunctions(t *testing.T) {
	t.Skip("skipping prototype tests")
	for _, q := range FunctionQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}

func Benchmark_Select(b *testing.B) {
	t.Skip("skipping prototype tests")
	context := Start()
	if _, err := InjestQuery(&context, "SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (product.id = 1 or product.test = true) and (products.account_id = '1') LIMIT 10 OFFSET 0;"); err != nil {
		b.Error(err)
	}
}
