package system

type NNode struct {
	NodeID    int
	Region    string
	IPAddress string
	Port      uint16
	Database  string
	User      string
	Password  string
	ReplicaOf *int
}
