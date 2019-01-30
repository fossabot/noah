# Key Value Store
Rather than using the underlying PostgreSQL database (Citus), or a separate distributed 
key-value store like etcd or zookeeper (Vitess); noah uses an embedded key value store
that is distributed using the raft consensus protocol. BadgerDB is used for the underlying
key-value store but it sits behind an abstraction layer called arctonyx. Arctonyx 
acts as a distribution layer for the key-value store using hasicorps raft library.

## Data
All of the data noah needs to plan a query and manage the cluster is stored in the 
key-value store. 

- Table Definitions
    - Name
    - Primary Key
    - Shard Key
    - Columns
- Table Sequences
- Tenant Definitions
    - Shards
- Nodes
    - Connection Info
    - Health

It helps to store all of this data in noah's own key value store because it allows noah
to quickly notify all of the coordinators in the cluster about changes in the cluster
state.

## Structure
Data nodes currently in the cluster.

`/noah/nodes/{active|inactive}/meta/{nodeId}/` 

#
Shards currently being managed by the cluster.

`/noah/shards/{active|inactive}/meta/{shardId}/`

#
A map of what nodes have what shards on them. A single shard is a single database in postgres.
There is also a reverse map of the same data to improve lookup time for data.

`/noah/nodes/{active|inactive}/shards/{nodeId}/{shardId}/`

`/noah/shards/{active|inactive}/nodes/{shardId}/{nodeId}/`

#
Accounts currently being managed by the cluster.

`/noah/accounts/{active|inactive}/meta/{accountId}/`

#
Accounts and what shards they exist on.

`/noah/accounts/{active|inactive}/shards/{accountId}/{shardId}/`



#
Given an `accountId` noah will perform the following key-value operations.

- Prefix scan (key only) of path `/noah/accounts/active/meta/{accountId}/`.
    - If the key exists then the operation continues, if it does not then the operation fails.
- Prefix scan (key only) of path `/noah/accounts/active/shards/{accountId}/`
    - If no shards are found then the operation fails.
- For each of the `shardIds` returned, do a prefix scan (key only) of the following path: `/noah/shards/active/nodes/{shardId}/`
    - If no nodes are found the the operation fails.
- For each of the distinct `nodeIds` returned, lookup the metadata for that node `/noah/nodes/active/meta/{nodeId}/`
    - If all of the nodes are not healthy then the operation fails.

With all of these operations we now have all of the information we need to connect to any of the account's shards
and perform the query.