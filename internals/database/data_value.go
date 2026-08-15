package database

type DataValue struct {
	Value any
	Ttl   uint64
}

func NewDataValueObject(value any, ttl uint64) *DataValue {
	return &DataValue{
		Value: value,
		Ttl:   ttl,
	}
}

