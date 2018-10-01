# Distributed ID Generation
This document describes the process Noah uses to generate unique numeric identifiers.

## Sequences
Noah supports sequences for inserting data into the database. 
Despite the name these are not precisely sequential. 

When a coordinator joins a cluster, it requests an allocation of sequential IDs for each sequence that currently exists. 
As that coordinator receives queries that require an ID, it will take from the array of IDs it has allocated for that specific sequence.
The next ID used will always be greater than the previous ID (on that coordinator), but it is not guaranteed to be consistently incremented. 
The amount the ID is incremented depends on how many chunks the master coordinator created when it generated the range.


Each node will receive a `Starting` and `Ending` value, as well as a `Node_Index` and `Number_Of_Nodes`. 
The node index is essentially that coordinator's sequence offset for ID generation. 
The number of nodes is how many coordinators received that specific range.

`Starting + (Node_Index) + (Number_Of_Nodes * Current_ID) - (Number_Of_Nodes - 1) = Next_ID`


For Example:
`1000 + (0) + (4 * 1) - (4 - 1)` Would return the next ID for a range of IDs starting at `1000`, with the range being distributed to 4 coordinators. 
This calculation is being performed by node `0` in the array of nodes that received that range; and the `Current_ID` for that node is 1.

The `Current_ID` is a value that is tracked on each coordinator and is never shared with others.
Once a coordinator reaches a given percentage of it's ID range, it will send a request to the `master` coordinator for more IDs.
It will continue to use it's current range until no more IDs can be generated. Then the `Current_ID` will be reset to `1` and the next range will be used.

If a coordinator node does not use all of the IDs it has allocated then those IDs are lost forever. 
A coordinator node doesn't ever report what IDs have been used back to the master node.


If a new coordinator node joins the cluster, it will either allocate the next available chunk of the current range, or the master will generate a new range with a new number of chunks to accommodate for the new node.

This process of generating semi-sequential but cluster-wide unique IDs is demonstrated in `distributed_id_test.go`. 
I have not yet implemented master ID generation either. 
In the test example the master node does not generate IDs.

## Snowflake ID

Noah also uses another form of unique ID generation based off of Twitter's "Snowflake". 
IDs are generated based off of a timestamp and the coordinator ID within the cluster.

These IDs will be primarily used for two-phase transaction IDs. 
Since they are sequential and can be generated without collaborating with other nodes in the cluster beyond getting an initial coordinator ID.
They also are very good for logging all write transactions that are sent to database nodes. 
Because they can be decomposed to reveal which coordinator processed the transaction and when (see more in `Transaction Replay Log`).


