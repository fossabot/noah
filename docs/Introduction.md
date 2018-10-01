# Introduction.

Noah is a distributed/sharded SQL database built on top of PostgreSQL. 
It is designed specifically for multi-tenant applications due to the method of mapping a `shard key` to a specific node.
Noah does not hash the shard key (referred to as the `account_id`) or use a range to distribute data across the nodes in the cluster. 
The shard key is stored in a distributed key-value store. This allows for more finite control over the distribution of data in the cluster.
You can move any data associated with an `account_id` to any server manually, allowing you to isolate accounts with higher demand or traffic to a private server.
When a new account is being created within Noah, the creation must be confirmed on the master node of the cluster. 
Accounts are the only record that demand this level of consensus from the cluster; this is because all of the active coordinators in the cluster need to be aware of all account_ids so that they can appropriately direct queries that they might receive.
