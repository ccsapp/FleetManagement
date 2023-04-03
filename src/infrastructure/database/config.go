package database

type Config interface {
	GetMongoDbConnectionString() string
	GetMongoDbDatabase() string
	GetAppCollectionPrefix() string
}
