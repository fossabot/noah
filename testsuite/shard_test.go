/*
 * Copyright (c) 2019 Ready Stock
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

package testsuite

import (
	"fmt"
	"github.com/readystock/golog"
	"testing"
)

// Create an accounts table, and a sharded table and then insert into the shard table.
func Test_Shard_Create_Simple(t *testing.T) {
	DoExecTest(t, ExecTest{
		Query: `DROP TABLE IF EXISTS public.accounts CASCADE;`,
	})
	DoExecTest(t, ExecTest{
		Query: `CREATE TABLE public.accounts (account_id BIGSERIAL PRIMARY KEY, account_number BIGINT) TABLESPACE "noah.account";`,
	})
	defer DoExecTest(t, ExecTest{
		Query: `DROP TABLE public.accounts CASCADE;`,
	})

	DoExecTest(t, ExecTest{
		Query: `DROP TABLE IF EXISTS public.users CASCADE;`,
	})
	DoExecTest(t, ExecTest{
		Query: `CREATE TABLE public.users (user_id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES public.accounts (account_id), user_number BIGINT) TABLESPACE "noah.shard";`,
	})
	defer DoExecTest(t, ExecTest{
		Query: `DROP TABLE public.users CASCADE;`,
	})

	numberOfAccountsToMake := 10
	numberOfUsersPerAccount := 20

	for i := 0; i < numberOfAccountsToMake; i++ {
		accountResult := DoQueryTest(t, QueryTest{
			Query: fmt.Sprintf(`INSERT INTO public.accounts (account_number) VALUES(%d) RETURNING account_id;`, i),
		})

		accountId := accountResult[0][0].(int64)

		golog.Warnf("created new account with ID [%d]", accountId)

		for u := 0; u < numberOfUsersPerAccount; u++ {
			userResult := DoQueryTest(t, QueryTest{
				Query: fmt.Sprintf(`INSERT INTO public.users (account_id, user_number) VALUES(%d, %d) RETURNING user_id;`, accountId, i),
			})

			userId := userResult[0][0].(int64)

			golog.Warnf("created new user with ID [%d] for account ID [%d]", userId, accountId)
		}
	}
}
