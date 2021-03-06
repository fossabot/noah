# Query Routing
Noah has a system that allows it to decide which servers serve which queries, some queries
will be served by multiple servers indirectly.

#
#### No tables
```postgresql
SELECT 'December 1st 2018'::timestamp;
```
If a query does not target any tables then the query will be sent to any node in the cluster.
This can have some strange side effects. If you were to send a query for something like a 
`current_timestamp` then you might get different results each time the query is sent.
Because of this it's recommended to configure all of your database servers to be in the UTC
timezone and the convert that time to the local time from that.

#### Targets only global tables
```postgresql
SELECT users.id FROM users;
```
The query above targets only 1 global table. Since the users table will be identical
on all database servers this query can be served by any node. This query is also safe, the results
returned by this query will be the same regardless of the node that serves it.

Default prostgres tables are also considered global. So certain SQL drivers that query OIDs or other
data from postgres tables can still get the data they need.

![alt text](/docs/images/OID_Query_Diagram.png "OID Diagram")

Because of this a lot of the types will be converted to a different type once Noah goes to
concatenate the data. For example; a `CITEXT` type would be converted to the OID of `TEXT`.
This may cause different issues with the actual driver as well if Noah returns the wrong type for an
object. But it will also make sure that all of the data returned from Noah will have a static OID 
that will be the same on every instance of postgres.

#### Targets only sharded tables
```postgresql
SELECT products.id FROM products;
```
If a query only targets a sharded table then it can only be served by a few different nodes; that is
any node that contains the data be queried. If there is no account filter in the query then the 
query will definitely be slow because a query will need to be sent to each node that MIGHT have the 
data being requested and then concatenated. Any sorting, grouping and aggregation is done after the
data is concatenated. This will make multi shard queries much much slower but allow pretty much
any query to be accepted by the coordinator.

```postgresql
SELECT products.id FROM products WHERE account_id = 12345;
```
This is the easiest query for Noah to handle, because it targets a single account ID and is a 
sharded table. Noah will lookup that account ID and find out which nodes have that data and send the
query to one of them.

#### Targets global and sharded tables
```postgresql
SELECT products.id FROM products INNER JOIN users ON users.id = products.created_by WHERE products.account_id = 12345;
```
It's very similar when joining with a global table. Because we are still only looking at data for a
single account ID and global tables are identical on all nodes, we can still just choose a node with
that account's data and send the query through without modification.


## Basic Query Flow Diagram
![alt text](/docs/images/Query_Flow.png "Query Flow Diagram")
