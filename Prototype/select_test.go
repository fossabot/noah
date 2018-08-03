package Prototype

import (
	"testing"
	"fmt"
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
	FunctionQueries = []string {
		"SELECT do_test('test');",
		"SELECT CURRENT_TIMESTAMP;",
	}
)

func Test_Select1(t *testing.T) {
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
	for _, q := range SelectQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}

func Test_SelectFunctions(t *testing.T) {
	for _, q := range FunctionQueries {
		context := Start()
		if _, err := InjestQuery(&context, q); err != nil {
			t.Error(err)
		}
	}
}



func Benchmark_Select(b *testing.B) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id INNER JOIN temp ON temp.id=products.product_id WHERE (product.id = 1 or product.test = true) and (products.account_id = '1') LIMIT 10 OFFSET 0;"); err != nil {
		b.Error(err)
	}
}
