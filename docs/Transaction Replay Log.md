# Transaction Replay Log
**(Noah does not have a transaction replay log, this was an idea early on, but ended up not being viable)**

**(This documentation only applies for database nodes that are replicated via statement replication and not PostgreSQL's built in replication)**

The transaction replay log is an array of transactions and their queries that have been sent to database nodes by a coordinator.

As each coordinator receives a query that modifies any data it will keep track of whether or not that transaction was successful, what database nodes were effected, and what time the transaction begun.
This log is not distributed throughout the cluster constantly. It is stored on the same server as the coordinator that handled the transaction.

All coordinators in a cluster are aware of the existence and status (alive | dead) of all other coordinators in the cluster. 

# A database node goes offline
When a database node goes offline, all of the coordinators in the cluster are notified. And any write query that effects accounts that have some data replicated to that offline node will be logged to the transaction replay log.

From here a few different things can happen:
1. The database node comes back online after a period less than 8 hours, and steps are taken to bring that node back up to date with the rest of the cluster.
2. The database node never comes back online, and after 8 hours is decommissioned by the cluster and considered completely dead.
3. The database node comes back online after after 8 hours; the node is then completely reset by the cluster and the under-replicated data is transferred.

## The database node comes back online before it is decommissioned
Transaction replay logs are stored for 12 hours, allowing for 4 hours to retrieve any required logs from the coordinators if the database node were to come back online at the 8 hour point.
The master coordinator will request transaction logs for all of the accounts that were on the recovering node from all coordinators in the cluster.
If there are any coordinators that have gone offline after the database node went offline then the master coordinator will not try to retrieve the transaction logs; since transaction logs could be missing due to the offline coordinator.
If not all the transaction logs are retrievable for any reason then the database node is completely reset and recovered.

If all of the transaction logs are retrievable, the master coordinator will start attempting to replay the transaction logs on that server in the order they were received by the cluster.
Any write queries that would effect this node during this process are also logged and sent to the master node right away so they can be sent at the end.
If all of the transactions were sent to the recovering node successfully then then the database rejoins the cluster and begins serving normal queries again.
If any of the transactions fail the database will be treated as corrupt and will be completely reset and recovered.

## The database node needs to be completely reset
When this happens it is preferable that the server itself be completely reset and the data directory of PostgreSQL be completely deleted and recreated.
Once the database is clean a base database is created called `ready`. The schema is created and any extensions needed are installed.
The extension `postgres_fdw` is installed in order to communicate with other database nodes in the cluster.
Once the schema is created the following query is sent to the recovering node. 
The query will give us a hierarchy of data within the database; also giving us the order in which we need to recreate the data.

```sql
WITH RECURSIVE ref (tbl, reftbl, depth) AS (
  SELECT pg_class.oid, NULL::oid, 0
  FROM pg_class
  JOIN pg_namespace ON
    pg_namespace.oid = pg_class.relnamespace
  WHERE 
    relkind = 'r' AND
    nspname = 'public' AND
    NOT EXISTS (
      SELECT 1 FROM pg_constraint
      WHERE 
        conrelid = pg_class.oid AND
        contype = 'f'
    )
  UNION ALL
  SELECT conrelid, ref.tbl, ref.depth + 1
  FROM ref
  JOIN pg_constraint ON
    confrelid = ref.tbl AND
    contype = 'f'
)
SELECT
  tbl::regclass::text as tablename,
  string_agg(DISTINCT reftbl::regclass::text, ',') as reftables
FROM ref
GROUP BY tablename
ORDER BY max(depth)
```

Starting with the highest table in the hierarchy, and taking into account which accounts need to have their data on this node; we spawn a connection using `postgres_fdw` to a node that contains the data for that account.

```sql
CREATE SERVER account_data_wrapper_0X FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host '10.0.0.2', dbname 'ready', port '5432');
```

Then a temporary foreign table is created to mirror the data on the other database node.

```sql
CREATE FOREIGN TABLE users_temp_0X (
    user_id     BIGINT,
    email       CITEXT NOT NULL,
    password    bytea[] NOT NULL
)
SERVER account_data_wrapper_0X;
```

Now the data is transferred to the recovering server in chunks of 1000.

```postgresql
INSERT INTO users (user_id, email, password)
SELECT
    ut.user_id,
    ut.email,
    ut.password
FROM users_temp_0X ut
LEFT JOIN users u ON u.user_id=ut.user_id
WHERE
    u.user_id IS NULL
LIMIT 1000;
```

This process is repeated until the rows effected is 0. Then it moves onto the next table in the hierarchy.
If this process fails at any point then an error is reported and the process is put on hold for an hour. 
After an hour the server is reset and the data transfer is attempted again.