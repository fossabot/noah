package conf

var (
 	//Configuration Config
)



type Config struct {
	AdminPort int            `json:"admin_port"`
	Database  DatabaseConfig `json:"database"`
	Cluster   ClusterConfig  `json:"cluster"`
	Nodes     []NodeConfig   `json:"nodes"`
}

type DatabaseConfig struct {
	AdvertiseAddress string `json:"advertise_address"`
	ReadBuffer       int    `json:"read_buffer"`
}

type ClusterConfig struct {
	DenyConnectionIfNoNodes    bool `json:"deny_connection_if_no_nodes"`
	StartingConnectionsPerNode int  `json:"starting_connections_per_node"`
	MaxConnectionsPerNode      int  `json:"max_connections_per_node"`
}

type NodeConfig struct {
	NodeID   int    `json:"node_id"`
	Address  string `json:"address"`
	Port     uint16 `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}
