# Open Source Components
Noah uses several community maintained libraries, the following libraries play a significant role
within noah's code base; or have been forked and modified to be used within noah.

## CockroachDB
**Project:** https://github.com/cockroachdb/cockroach

Copyright 2018 The Cockroach Authors.

License (Apache License 2.0) https://github.com/cockroachdb/cockroach/blob/master/LICENSE

> Noah's internal wire protocol and message handling is inherited directly from CockroachDB,
the connExecutor is also used to handle query processing. But the query parser, key value store
and raft implementations have all been removed and replaced with a custom version to route queries
to an external Postgres database.

## Vitess
**Project:** https://github.com/vitessio/vitess

Copyright 2018 Google Inc.

License (Apache License 2.0) https://github.com/vitessio/vitess/blob/master/LICENSE

> Some approaches to sharding, scaling and query handling were inspired heavily by the behavior of Vitess. 
However at this time no Vitess code exists within Noah's code base, but credit is given due to the
heavily mirrored behavior in some operations.

## Citus
**Project:** https://github.com/citusdata/citus

Copyright 2018 Citus Data, Inc.

License (GNU Affero General Public License v3.0) https://github.com/citusdata/citus/blob/master/LICENSE

> Noah is designed to compete directly with Citus, it uses a similar sharding system to Citus in
that they both use a key value lookup to determine which shard data exists on for a single query.

## pg_query_go
**Project:** https://github.com/lfittl/pg_query_go

Copyright 2018 Lukas Fittl

License (3-Clause BSD) https://github.com/lfittl/pg_query_go/blob/master/LICENSE

> pg_query_go is one of the largest tools at the heart of noah. It provides all of noah's query 
parsing and compiling. Allowing noah to interpret queries precisely and rewrite them on the fly 
without syntax errors. 

## pgx
**Project:** https://github.com/jackc/pgx

Copyright 2018 Jack Christensen

License (MIT) https://github.com/jackc/pgx/blob/master/LICENSE

> A heavily modified version of pgx is used to allow noah to talk to the Postgres nodes in the 
cluster. PGX was stripped down, removing all non simple query functions (noah does not pass 
arguments to prepared queries when communicating to database nodes). 2-phase commit was also added 
to noah's implementation of pgx to allow for consistency when writing to multiple database nodes.
The method of returning rows was also tweaked slightly to allow for a much easier method of merging
results from multiple servers into a single response to the client.

## BadgerDB
**Project:** https://github.com/dgraph-io/badger

Copyright 2018 Dgraph Labs, Inc. and Contributors

License (MIT) https://github.com/dgraph-io/badger/blob/master/LICENSE

> BadgerDB powers noah's internal key-value store. It is controlled by another ready stock library 
called arctonyx (a bread of badger), which distributes BadgerDB using Hashicorp's implementation of 
the raft protocol. BadgerDB keeps track of all of the tenants, tables, columns, nodes and sequences 
in the cluster.

## Sonyflake
**Project:** https://github.com/sony/sonyflake

Copyright 2018 Sony Corporation

License (MIT) https://github.com/sony/sonyflake/blob/master/LICENSE

> Sonyflake is Sony's implementation of Twitter's "Snowflake ID" generation. Noah has two systems of
unique ID generation built into it, allowing for ID columns in Postgres to be used. The snowflake 
library will generate IDs slightly faster than noah's standard sequence generation, but will have 
larger IDs and have larger gaps between IDs when large amounts of time have passed between 
generations.

## raft
**Project:** https://github.com/hashicorp/raft

Copyright 2018 HashiCorp

License (MPL-2.0) https://github.com/hashicorp/raft/blob/master/LICENSE

> Ready stock forked the raft library and implemented a gRPC transport, allowing noah to serve it's 
internal gRPC methods along side the raft interface all over the same wire.

## pq
**Project:** github.com/lib/pq

Copyright 2018  'pq' Contributors Portions Copyright (C) 2018 Blake Mizerany

License https://github.com/lib/pq/blob/master/LICENSE.md

> pq is used in unit tests for noah to ensure that noah's wire protocol complies with common 
Postgres client libraries.

## go-linq
**Project:** github.com/ahmetb/go-linq

Copyright 2018 Ahmet Alp Balkan

License https://github.com/ahmetb/go-linq/blob/master/LICENSE

> go-linq is used in noah to break up large pieces of data from BadgerDB into smaller grouped pieces
of data. Like when a table needs to be altered, we need to get a list of all of the nodes in the
cluster, then only apply the query to nodes that are not replicas that are alive (but only if all of
the non-replica nodes are alive, if any are not then we want to reject the schema change). This is 
made simpler by querying the array of nodes with LINQ.