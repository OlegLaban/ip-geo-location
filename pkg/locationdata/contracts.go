package locationdata

type CacheInterface interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	GetWithCallback(key string, callback func () ([]byte, error)) ([]byte, error)
	Delete(key string) error
	Has(key string) bool
}