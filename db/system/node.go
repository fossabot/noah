package system

type NNode struct {
	NodeID    uint64
	Region    string
	IPAddress string
	Port      uint16
	Database  string
	User      string
	Password  string
	ReplicaOf *int
}
