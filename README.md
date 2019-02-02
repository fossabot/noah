# Noah 
[![Build Status](https://travis-ci.com/readystock/noah.svg?token=QvXZjJzgiir2JHLaKFrG&branch=master)](https://travis-ci.com/readystock/noah)
[![CodeFactor](https://www.codefactor.io/repository/github/readystock/noah/badge)](https://www.codefactor.io/repository/github/readystock/noah)

Noah is a database sharding tool for Postgres inspired by Vitess and Citus. 
At it's core is a gutted version of CockroachDB, rewritten to use multiple PostgreSQL databases as it's primary storage engine.
The query parsing of Cockroach has also been replaced by PostgreSQL's actual query parser in C. 
It communicates with the database servers using a bare-bone and customized version of PGX, with all non-essential features removed and two-phase transaction handling added.

PLEASE NOTE: Noah is still a work in progress, while some of it's functionality does work currently; it is nowhere near a production worthy state. 

## Building

Noah requires `protoc` to be built. Protoc will generate the protos needed for
some of the key value store data and wire protocol. Cockroach had some protos
that have been "hard-coded" into noah to make builds easier. These protos also
won't be modified so there should be no need to regenerate them.

Any tools that you might need can be found in [tools](./docs/Tools.md).

To build noah.
```bash
go get -d -u github.com/readystock/noah
cd $GOPATH/src/github.com/readystock/noah
make
```

This will create a `bin` folder if it does not already exist and create the noah executable within 
that folder.

#### Clean build
Please note: Noah cannot be built on Windows at this time.


If you need to do a completely clean build of noah run the following:
```bash
make fresh
```
This will recompile all of noah's dependencies and may take a few minutes.

## Example

#### 1. Create an accounts table:

```postgresql
CREATE TABLE public.accounts (
  account_id   BIGSERIAL NOT NULL PRIMARY KEY,
  account_name TEXT      NOT NULL
) TABLESPACE "noah.account";
```

Noah will recognize that this table will be a shard key for sharded tables. 
It will look for tables that have a foreign key reference to the primary key of the created
`accounts` table and then map the data for that ID to a node in the cluster.

#### 2. Create some global tables:

```postgresql
CREATE TABLE public.logins (
  login_id BIGSERIAL NOT NULL PRIMARY KEY,
  email    TEXT      NOT NULL UNIQUE,
  name     TEXT      NULL
);
```

This table will be used when a user initially attempts to login. By not specifying
a Noah tablespace; it will assume the table should be global. All of the data in this
table will be mirrored on every node in the cluster.

Then we can create another global table to map logins to an account, allowing a single
login to be associated with multiple accounts.

```postgresql
CREATE TABLE public.users (
  user_id    BIGSERIAL NOT NULL PRIMARY KEY,
  login_id   BIGINT    NOT NULL REFERENCES public.logins (login_id),
  account_id BIGINT    NOT NULL REFERENCES public.accounts (account_id)
) TABLESPACE "noah.global";
```

You can also manually specify the global tablespace when creating a table.

#### 3. Create some sharded tables:

```postgresql
CREATE TABLE public.products (
  product_id BIGSERIAL NOT NULL PRIMARY KEY,
  account_id BIGINT    NOT NULL REFERENCES public.accounts (account_id),
  sku        TEXT      NOT NULL,
  title      TEXT      NULL,
  CONSTRAINT uq_products_sku_per_account UNIQUE (account_id, sku)
) TABLESPACE "noah.shard";
```

This table stores all of the products for an account. Because it has a foreign key referencing
the `accounts` table Noah will know how to handle any operations that target the `products` 
table. 

```postgresql
CREATE TABLE public.orders (
  order_id     BIGSERIAL NOT NULL PRIMARY KEY,
  account_id   BIGINT    NOT NULL REFERENCES public.accounts (account_id),
  order_number TEXT      NOT NULL,
  marketplace  TEXT      NOT NULL
) TABLESPACE "noah.shard";
```

```postgresql
CREATE TABLE public.order_lines (
  order_line_id BIGSERIAL NOT NULL PRIMARY KEY,
  account_id    BIGINT    NOT NULL REFERENCES public.accounts (account_id),
  order_id      BIGINT    NOT NULL REFERENCES public.orders (order_id),
  product_id    BIGINT    NOT NULL REFERENCES public.products (product_id),
  qty           NUMERIC   NOT NULL
) TABLESPACE "noah.shard";
```

We can then create an `orders` and `order_lines` table that will have their data co-located
in the cluster with data for their given `account_id`. This allows for strong consistency
with foreign keys in massive data sets, but also extremely fast performance since the 
data is split into chunks (a smaller index is a better index).

### 4. Insert some data:

Create a few accounts:

```postgresql
INSERT INTO public.accounts (account_name) 
VALUES
       ('Do Stuff Inc.'),        -- account_id = 1
       ('Les Vegetables LLC.'),  -- account_id = 2
       ('Dunder Mifflin');       -- account_id = 3
```
(Although it is not present in this example, I highly recommend using the `RETURNING` clause as much
as possible in your code where you need to get the ID of a created record. It saves significant
performance from round-trip queries and allows follow-up records to be created immediately; and 
for sharded tables, it will save a potentially expensive query plan.)

Create a few users and logins:

```postgresql
INSERT INTO public.users (email, name) 
VALUES
       ('elliot@elliot.elliot', 'Not Elliot'), -- user_id = 1
       ('bob@bob.com', 'Billy'),               -- user_id = 2
       ('ron@gmail.com', 'Ron Swanson');       -- user_id = 3
```

```postgresql
INSERT INTO public.logins (user_id, account_id) 
VALUES
       (1, 1),
       (1, 2),
       (2, 2),
       (3, 3);
```

You can now run queries for a login. Like if we wanted to see the accounts for a given
set of credentials.

```postgresql
SELECT accounts.* FROM public.logins JOIN public.accounts ON  logins.account_id = accounts.account_id JOIN public.users ON logins.user_id = users.user_id WHERE users.email = 'elliot@elliot.elliot' LIMIT 1;
```

This will give us a list of account_ids and account_names for a given login.

Now we can create some products and orders for some sharded queries.

```postgresql
INSERT INTO public.products (account_id, sku, title) 
VALUES 
       (1, 'ABC123', 'A title for a product'),  -- product_id = 1
       (1, 'WHEELS-4', 'A set of 4 wheels'),    -- product_id = 2
       (2, 'JACKET-LG', 'A large jacket'),      -- product_id = 3
       (3, 'LAPTOP-CRAP', 'A crappy laptop.');  -- product_id = 4
```

Although the use case is probably rare, Noah does support inserting into multiple accounts from a
single query. It will parse this query and break it up into separate queries if the accounts happen
to reside on different nodes.

We can then query our inserted data by providing the `account_id` column in our select query.

```postgresql
SELECT * FROM public.products WHERE account_id = 1;
```

If you try to query a sharded table without providing the shard column then Noah will return an 
error at this time. (More info will be provided in docs later)

Add some order data for account 1.

```postgresql
INSERT INTO public.orders (account_id, order_number, marketplace) 
VALUES
       (1, '123-4271841-32141', 'Amazon'),  -- order_id = 1
       (1, '312-8231982-98392', 'Amazon');  -- order_id = 2
```

```postgresql
INSERT INTO public.order_lines (account_id, order_id, product_id, qty)
VALUES 
       (1, 1, 1, 3),
       (1, 1, 2, 4),
       (1, 2, 2, 1);
```

Now that there is some data in the sharded tables we can query and join those too.

```postgresql
SELECT 
       orders.order_number, 
       orders.marketplace, 
       products.sku, 
       order_lines, qty
FROM public.orders
JOIN public.order_lines ON order_lines.order_id = orders.order_id
JOIN public.products ON products.product_id = order_lines.product_id
WHERE orders.account_id = 1;
```