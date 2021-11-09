package cache

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func TestGroupKindBoolMap_Reload(t *testing.T) {
	m := &GroupKindBoolMap{}
	m.Store(schema.GroupKind{Group: "group", Kind: "k1"}, true)
	m.Store(schema.GroupKind{Group: "group", Kind: "k2"}, true)
	m.Reload(map[schema.GroupKind]bool{
        {Group: "group", Kind: "k3"}: true,
    })

	assert.Equal(t, 1, m.Len())
}

func TestNamespaceResourcesMap_LoadAndDelete(t *testing.T) {
	m := NamespaceResourcesMap{}
	m.Store("a", &ResourceMap{})
	m.Store("b", &ResourceMap{})
	existed, _ := m.LoadAndDelete("a")
	assert.NotNil(t, existed)
	assert.NotNil(t, existed.All())
}
