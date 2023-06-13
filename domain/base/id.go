package base

import "github.com/google/uuid"

type UUID uuid.UUID

var Nil = UUID(uuid.Nil)

func NewUUID() UUID {
	return UUID(uuid.New())
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}

func UUIDFromString(id string) (UUID, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return Nil, err
	}
	return UUID(parsedUUID), nil
}

func UUIDFromBytes(id []byte) (UUID, error) {
	parsedUUID, err := uuid.FromBytes(id)
	if err != nil {
		return Nil, err
	}
	return UUID(parsedUUID), nil
}
