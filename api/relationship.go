package api

type Relationship interface {
	Name() string
	Description() string
}

type RelationshipImpl struct {
	name        string
	description string
}

func NewRelationship(name string, description string) Relationship {
	return &RelationshipImpl{name, description}
}

func (relationship *RelationshipImpl) Name() string {
	return relationship.name
}

func (relationship *RelationshipImpl) Description() string {
	return relationship.description
}
