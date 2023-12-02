package dbcassandra

import (
	"fmt"
	"github.com/gocql/gocql"
	"time"
)

func CreateConnectionWithAuth(Servers ...string) *gocql.Session {
	cluster := gocql.NewCluster(Servers...)
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	cluster.Compressor = &gocql.SnappyCompressor{}
	cluster.RetryPolicy = &gocql.ExponentialBackoffRetryPolicy{NumRetries: 3}
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 10 * time.Second
	cluster.Timeout = 10 * time.Second
	cluster.ProtoVersion = 4
	cluster.DisableInitialHostLookup = true

	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Println("Cassandra cluster.CreateSession err : ", err)
	}

	return session

}
