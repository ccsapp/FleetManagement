package database

import (
	"PFleetManagement/logic/fleetErrors"
	"PFleetManagement/logic/model"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const fleetCollectionBaseName = "fleets"

type connection struct {
	database   *mongo.Database
	client     *mongo.Client
	collection string
}

type fleet struct {
	FleetId model.FleetID `bson:"_id"`
	Vins    []model.Vin   `bson:"vins"`
}

func OpenDatabase(config Config) (FleetDB, error) {
	m := connection{}
	return &m, m.setUpDatabase(config) // return the error (if) encountered in setup
}

func toConnectionUri(config Config) string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/%s",
		config.GetMongoDbUser(),
		config.GetMongoDbPassword(),
		config.GetMongoDbHost(),
		config.GetMongoDbPort(),
		config.GetMongoDbDatabase(),
	)
}

func (m *connection) setUpDatabase(config Config) error {
	// create the client options and construct the MongoDB connection URI from environment variables
	opts := options.Client()
	opts.ApplyURI(toConnectionUri(config))

	var err error

	// create a context with 5s timeout for connecting with MongoDB (see below)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// connect to the MongoDB server
	m.client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	// ensure connection was successful
	if err = m.client.Ping(ctx, nil); err != nil {
		return err
	}

	// store an additional pointer to the database of which the name is given by the environment
	m.database = m.client.Database(config.GetMongoDbDatabase(), options.Database())

	// save the collection name
	m.collection = config.GetAppCollectionPrefix() + fleetCollectionBaseName
	return nil
}

func (m *connection) CleanUpDatabase() error {
	return m.client.Disconnect(context.Background())
}

func (m *connection) AddFleet(ctx context.Context, fleetId model.FleetID) error {
	// create a new object with the given fleet ID and an empty car/VIN list
	_, err := m.database.Collection(m.collection).
		InsertOne(ctx, fleet{FleetId: fleetId, Vins: []model.Vin{}})

	// MongoDB detects duplicate _id (in BSON, field FleetId in struct)
	if mongo.IsDuplicateKeyError(err) {
		return fleetErrors.ErrFleetAlreadyExists
	}

	return err
}

func (m *connection) AddCarToFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error {
	// update only the fleet given by ID
	filter := bson.D{{"_id", fleetId}}
	// the $addToSet operator guarantees that the VIN will not occur multiple times in the resulting array/set
	update := bson.D{{"$addToSet", bson.D{{"vins", vin}}}}
	// perform the addition as atomic MongoDB document update (updateOne does not fail, if none found, see below)
	result, err := m.database.Collection(m.collection).UpdateOne(ctx, filter, update)

	if err != nil {
		// return database error
		return err
	}
	if result.MatchedCount == 0 {
		// this case occurs if the filter (by fleet ID) did not match any document -> no fleet with that ID exists
		return fleetErrors.ErrFleetNotFound
	}
	if result.ModifiedCount == 0 {
		//  -> the VIN was already in the set
		// if a document matched but none was modified, the $addToSet did not perform a change
		return fleetErrors.ErrCarAlreadyInFleet
	}

	// no error nor invalid post conditions -> success
	return nil
}

func (m *connection) GetCarsForFleet(ctx context.Context, fleetId model.FleetID) ([]model.Vin, error) {
	var fleet fleet

	// create a query with filter by _id (aka FleetId) and decode the document to the struct, fails if fleet not found
	err := m.database.Collection(m.collection).
		FindOne(ctx, bson.D{{"_id", fleetId}}).
		Decode(&fleet)

	if err == mongo.ErrNoDocuments {
		// this error is returned if no fleet matched the ID filter
		return nil, fleetErrors.ErrFleetNotFound
	}

	// the fleet.Vins might be invalid data if Decode returned an error but errors should be checked first
	return fleet.Vins, err
}

func (m *connection) IsCarInFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) (bool, error) {
	// just delegates to read list of VINs of all cars assigned to fleet
	vins, err := m.GetCarsForFleet(ctx, fleetId)

	// handle downstream errors
	if err != nil {
		return false, err
	}

	// iterate through VINs and return true (car is in fleet), as soon as the VIN to check for is found
	for _, foundVin := range vins {
		if foundVin == vin {
			return true, nil
		}
	}

	// no return from inside the loop -> no VIN matched -> car not in fleet
	return false, nil
}

func (m *connection) RemoveCarFromFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error {
	// update only the fleet given by ID
	filter := bson.D{{"_id", fleetId}}
	// the $pullAll operator deletes (an array of) values from an array atomically
	update := bson.D{{
		"$pullAll", bson.D{{
			"vins", bson.A{vin},
		}},
	}}
	// perform the atomic update
	result, err := m.database.Collection(m.collection).UpdateOne(ctx, filter, update)

	if err != nil {
		// return database error
		return err
	}
	if result.MatchedCount == 0 {
		// this case occurs if the filter (by fleet ID) did not match any document -> no fleet with that ID exists
		return fleetErrors.ErrFleetNotFound
	}
	if result.ModifiedCount == 0 {
		// if a document matched but none was modified, the $pullAll did not perform a change
		//  -> the VIN was not contained in the array
		return fleetErrors.ErrCarNotInFleet
	}

	// no error nor invalid post conditions -> success
	return nil
}

func (m *connection) DropCollection(ctx context.Context) error {
	return m.database.Collection(m.collection).Drop(ctx)
}
