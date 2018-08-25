package cluster

type NTableType int64

const (
	NTableType_Standard  NTableType = 1 >> iota // Standard tables are tables that exist on every node in the cluster and are sharded by a tenant_id
	NTableType_Tenant                           // A tenant table is a table that keeps track of tenants. Only one of these types of tables are allowed. And they are replicated on all servers.
	NTableType_Reference                        // A reference table is not sharded, it has the same data in it on every server in the cluster.
)

type NTable struct {
	Schema string
	Name   string
	Type   NTableType
}
