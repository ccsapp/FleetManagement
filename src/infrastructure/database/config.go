package database

type Config interface {
	GetMongoDbHost() string
	GetMongoDbPort() int
	GetMongoDbDatabase() string
	GetMongoDbUser() string
	GetMongoDbPassword() string
	GetAppCollectionPrefix() string
}
