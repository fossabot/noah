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
	"testing"
)

func TestSelect1(t *testing.T) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT 1;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAll1(t *testing.T) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT 'test';"); err != nil {
		t.Error(err)
	}
}

func Benchmark_InjestQuery_SelectSimpleAll1(t *testing.B) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT products.product_id FROM products WHERE (account_id = '1') AND (product_id = '2') LIMIT 10 OFFSET 0;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAccount1(t *testing.T) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT product_id,sku,title FROM products WHERE (account_id = '2') LIMIT 1;"); err != nil {
		t.Error(err)
	}
}

func TestSelectSimpleAccount2(t *testing.T) {
	context := Start()
	if _, err := InjestQuery(&context, "SELECT product_id,sku,title FROM products WHERE account_id='1232421' LIMIT 1;"); err != nil {
		t.Error(err)
	}
}
