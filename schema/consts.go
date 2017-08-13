package schema

var (
	// Users is a bucket name for users
	UsersBucketName = []byte{85, 115, 101, 114, 115}
	// Locations is a bucket name for locations
	LocationsBucketName = []byte{76, 111, 99, 97, 116, 105, 111, 110, 115}
	// Visits is a bucket name for visits
	VisitsBucketName = []byte{86, 105, 115, 105, 116, 115}
)

var (
	// Buckets is a array of all buckets for bolt db
	Buckets = [][]byte{
		UsersBucketName,
		LocationsBucketName,
		VisitsBucketName,
	}

	// BucketsMap is a map of entity to bucket
	BucketsMap = map[Entity][]byte{
		EntityUsers:     UsersBucketName,
		EntityLocations: LocationsBucketName,
		EntityVisits:    VisitsBucketName,
	}
)

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
