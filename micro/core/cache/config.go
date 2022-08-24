package cache

type Config struct {
	Default string
	Expire  int

	FreeCache struct {
		CacheSize  int
		Expiration int
	}
}
