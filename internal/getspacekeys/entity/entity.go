package entity

import "fmt"

type Space struct {
	id   string
	name string
	key  string
}

func NewSpace(id, name, key string) *Space {
	return &Space{id, name, key}
}

func (s *Space) AsString() string {
	return fmt.Sprintf("ID: %6v | Key: %.10v... | Name: %s", s.id, s.key, s.name)
}
