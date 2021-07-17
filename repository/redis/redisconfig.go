package redis

type Config struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
	Timeout  int
}
