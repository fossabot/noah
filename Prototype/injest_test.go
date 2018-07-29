package Prototype

import (
	"testing"
)

func TestSelect1(t *testing.T) {
	if err := InjestQuery("SELECT 1;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAll1(t *testing.T) {
	if err := InjestQuery("SELECT product_id,sku,title FROM products LIMIT 1;"); err != nil {
		t.Error(err)
	}
}

func Benchmark_InjestQuery_SelectSimpleAll1(t *testing.B) {
	if err := InjestQuery("SELECT product_id,sku,title FROM products LIMIT 1;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAccount1(t *testing.T) {
	if err := InjestQuery("SELECT product_id,sku,title FROM products WHERE account_id=1 AND product_id=2 LIMIT 1;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAccount2(t *testing.T) {
	if err := InjestQuery("SELECT product_id,sku,title FROM products WHERE account_id='1232421' LIMIT 1;"); err != nil {
		t.Error(err)
	}
}
