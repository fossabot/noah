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

