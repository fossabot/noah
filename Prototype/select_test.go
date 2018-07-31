package Prototype

import (
	"testing"
)

func Test_Selects(t *testing.T) {
	queries := []string{
		"SELECT 1",
		"SELECT 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id=1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = 1;",
		"SELECT products.product_id,products.sku FROM products WHERE account_id = '1' OR product_id = '2';",
		"SELECT products.product_id FROM products WHERE (account_id = '1') AND (product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id FROM products WHERE (products.account_id = '1') AND (products.product_id = '2') LIMIT 10 OFFSET 0;",
		"SELECT products.product_id,variations.variation_id FROM products INNER JOIN variations ON variations.product_id=products.product_id INNER JOIN users ON users.id=products.id WHERE (account_id = '2') LIMIT 10 OFFSET 0;",
	}
	for _, q := range queries {
		if err := InjestQuery(q); err != nil {
			t.Error(err)
		}
	}
}
