package cassandra

import (
	"github.com/gocql/gocql"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type CassandraConn struct {
	Hosts       []string
	Port        string
	Keyspace    string
	Consistency string
	User        string
	Password    string
	CaPath      string
}

type Cassandra struct {
	Session *gocql.Session
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func (conn *CassandraConn) InitSession() *gocql.Session {
	port := func(p string) int {
		i, err := strconv.Atoi(p)
		if err != nil {
			return 9042
		}

		return i
	}

	consistancy := func(c string) gocql.Consistency {
		gc, err := gocql.ParseConsistencyWrapper(c)
		if err != nil {
			return gocql.All
		}

		return gc
	}

	cluster := gocql.NewCluster(strings.Join(conn.Hosts, ","))
	cluster.Port = port(conn.Port)
	cluster.Keyspace = conn.Keyspace
	cluster.Consistency = consistancy(conn.Consistency)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conn.User,
		Password: conn.Password,
	}
	cluster.SslOpts = &gocql.SslOptions{
		CaPath: conn.CaPath,
	}
	cluster.ConnectTimeout = time.Second * 60
	cluster.DisableInitialHostLookup = true
	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("ERROR: fail create cassandra session, %s", err.Error())
		os.Exit(1)
	}

	return session
}

func (cassandra *Cassandra) Clear() {
	cassandra.Session.Close()
}
