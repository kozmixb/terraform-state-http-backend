package drivers

type Driver interface {
	Exists(group string, key string) bool
	Show(group string, key string) string
	Update(group string, key string, payload string) string
	// Lock(group string, key string) []byte
	// lock(group string, key string) []byte
}
