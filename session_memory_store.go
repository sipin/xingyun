package xingyun

// Simple session storage using memory, handy for development
// **NEVER** use it in production!!!
type memoryStore struct {
	data map[string][]byte
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{
		data: make(map[string][]byte),
	}
}

func (ms *memoryStore) SetSession(sessionID string, key string, data []byte) {
	ms.data[sessionID+key] = data
}

func (ms *memoryStore) GetSession(sessionID string, key string) []byte {
	data, _ := ms.data[sessionID+key]
	return data
}

func (ms *memoryStore) ClearSession(sessionID string, key string) {
	delete(ms.data, sessionID+key)
}
