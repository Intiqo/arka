package cache

type localCacheManager struct {
	store map[string]interface{}
}

func (r *localCacheManager) initialize() {
	r.store = make(map[string]interface{})
}

func (r *localCacheManager) GetStatus() string {
	return "UP"
}

func (r *localCacheManager) Set(key string, val interface{}) error {
	r.store[key] = val
	return nil
}

func (r *localCacheManager) Get(key string) (string, error) {
	if r.store[key] == nil {
		return "", nil
	}
	val := r.store[key]
	switch val.(type) {
	case string:
		return r.store[key].(string), nil
	case []uint8:
		return string(r.store[key].([]uint8)), nil
	}
	return r.store[key].(string), nil
}
