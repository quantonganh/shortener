package cassandra

import (
	"github.com/gocql/gocql"
)

type DB struct {
	hosts []string
	session *gocql.Session
}

func NewDB(hosts ...string) *DB {
	return &DB{
		hosts: hosts,
	}
}

func (db *DB) Open() error {
	cluster := gocql.NewCluster(db.hosts...)
	cluster.Keyspace = "shortener"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "cassandra",
		Password: "cassandra",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return err
	}
	db.session = session
	return nil
}

func (db *DB) Close() error {
	if db.session != nil {
		db.session.Close()
	}
	return nil
}