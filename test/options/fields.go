package options

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type Fields []Field

// Map creates a map from the elements of the D.
func (f Fields) Map() Map {
	m := make(Map, len(f))
	for _, e := range f {
		m[e.Key] = e.Value
	}
	return m
}

// E represents a BSON element for a D. It is usually used inside a D.
type Field struct {
	Key   string
	Value interface{}
}

type Map map[string]interface{}

func (m Map) Has(key string) bool {
	_, exists := m[key]
	return exists
}

func (m Map) String(key string, defaultValue string) string {
	if v, exists := m[key]; exists {
		return v.(string)
	}
	return defaultValue
}

func (m Map) ObjectID(key string, defaultValue objectid.ObjectID) objectid.ObjectID {
	if v, exists := m[key]; exists {
		return v.(objectid.ObjectID)
	}
	return defaultValue
}

func (m Map) Time(key string, defaultValue time.Time) time.Time {
	if v, exists := m[key]; exists {
		return v.(time.Time)
	}
	return defaultValue
}
