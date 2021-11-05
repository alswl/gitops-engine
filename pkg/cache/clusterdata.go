package cache

import (
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sync"
)

// ApisMetaMap is thread-safe map of apiMeta
type ApisMetaMap struct {
	log     logr.Logger
	syncMap sync.Map
}

func (m *ApisMetaMap) Load(gk schema.GroupKind) (*apiMeta, bool) {
	val, ok := m.syncMap.Load(gk)
	typedVal, typeOk := val.(*apiMeta)
	if !ok || !typeOk {
		return nil, false
	}
	return typedVal, true
}

func (m *ApisMetaMap) Store(gk schema.GroupKind, meta *apiMeta) {
	m.syncMap.Store(gk, meta)
}

func (m *ApisMetaMap) Delete(gk schema.GroupKind) {
	m.syncMap.Delete(gk)
}

//Range loops the map, and Range ensures every item will be load, but not guarantee missing(phantom read)
func (m *ApisMetaMap) Range(fn func(key schema.GroupKind, value *apiMeta) bool) {
	m.syncMap.Range(func(key, value interface{}) bool {
		typedKey, keyTypeOk := key.(schema.GroupKind)
		typedValue, valueTypeOk := value.(*apiMeta)
		if !keyTypeOk || !valueTypeOk {
			m.log.Info("Failed to cast key and value to GroupKind and *apiMeta")
			return false
		}
		return fn(typedKey, typedValue)
	})
}

//Len return ApisMetaMap length, roughly, it depends on the time point of Range each loop
func (m *ApisMetaMap) Len() int {
	length := 0
	m.syncMap.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}
