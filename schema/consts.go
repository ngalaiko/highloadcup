package schema

// GetIEntity returns IEntity by Entity
func GetIEntity(entity Entity) IEntity {
	switch entity {
	case EntityUsers:
		return &User{}
	case EntityVisits:
		return &Visit{}
	case EntityLocations:
		return &Location{}
	default:
		return nil
	}
}
