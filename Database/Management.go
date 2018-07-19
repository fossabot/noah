package Database

import (
	"github.com/Ready-Stock/Noah/Configuration"
	"github.com/kataras/iris/core/errors"
	"net"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)


var (
	healthyNodes = map[int]bool{}
)



func SetupNodes() error {
	if len(Conf.Configuration.Nodes) == 0 {
		return errors.New("Error, no nodes found in config!")
	}

	for _, node := range Conf.Configuration.Nodes {
		if _, err := net.ResolveTCPAddr("tcp", node.Address); err != nil {
			return err
		} else {
			n := GetNodeFromConfig(node)
			Nodes[n.NodeID] = n
		}
	}
	go nodeHealthChecker()
	return nil
}

func GetNodeIsHealthy(nodeId int) bool {
	if val, ok := healthyNodes[nodeId]; ok && val {
		return true
	} else {
		return false
	}
}


func nodeHealthChecker() {
	for {
		healthy := len(Nodes)
		for _, node := range Nodes {
			if db, err := gorm.Open("postgres", node.ConnectionString); err != nil {
				fmt.Println("Error, node", node.NodeID, "is not healthy (", err, ")")
				healthyNodes[node.NodeID] = false
				healthy--
			} else {
				if v := db.Exec("SELECT 1").RowsAffected; v != 1 {
					healthyNodes[node.NodeID] = false
					healthy--
				} else {
					healthyNodes[node.NodeID] = true
				}
				db.Close()
			}
		}
		fmt.Println(healthy, "/", len(Nodes), "Node(s) Healthy")
		time.Sleep(5 * time.Second)
	}
}

func GetNodeFromConfig(n Conf.NodeConfig) Node {
	return Node {
		NodeID: n.NodeID,
		ConnectionString:"postgresql://" + n.User + ":" + n.Password + "@" + n.Address + "/" + n.Database + "?sslmode=disable",
		Shards:n.Shards,
	}
}

