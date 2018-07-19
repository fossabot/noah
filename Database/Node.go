package Database

var (
	Nodes = map[int]Node{}
)


type Node struct {
	NodeID int
	ConnectionString string
	Shards int64
}
