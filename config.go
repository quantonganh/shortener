package shortener

type Config struct {
	Redis struct {
		Addr string
	}
	Cassandra struct {
		Hosts []string
	}
}
