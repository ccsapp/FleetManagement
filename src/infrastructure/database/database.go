package database

import (
	"PFleetManagement/logic/errors"
	"PFleetManagement/logic/model"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

const fleetCollectionName = "fleets"

type connection struct {
	database *mongo.Database
	client   *mongo.Client
}

type fleet struct {
	FleetId model.FleetID `bson:"_id"`
	Vins    []model.Vin   `bson:"vins"`
}

func OpenDatabase() (FleetDB, error) {
	m := connection{}
	return &m, m.SetUpDatabase()
}

func (m *connection) SetUpDatabase() error {
	opts := options.Client()
	opts.ApplyURI("mongodb://" + os.Getenv("MONGODB_DATABASE_USER") + ":" + os.Getenv("MONGODB_DATABASE_PASSWORD") + "@" + os.Getenv("MONGODB_DATABASE_HOST") + ":" + "27017" + "/" + os.Getenv("MONGODB_DATABASE_NAME"))

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	m.client, err = mongo.Connect(ctx, opts)
	if err != nil {
		return err
	}

	m.database = m.client.Database(os.Getenv("MONGODB_DATABASE_NAME"), options.Database())
	return nil
}

func (m *connection) Disconnect() error {
	return m.client.Disconnect(context.Background())
}

func (m *connection) AddFleet(ctx context.Context, fleetId model.FleetID) error {
	_, err := m.database.Collection(fleetCollectionName).InsertOne(ctx, fleet{FleetId: fleetId, Vins: []model.Vin{}})
	if mongo.IsDuplicateKeyError(err) {
		return errors.ErrFleetAlreadyExists
	}
	return err
}

func (m *connection) AddCarToFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error {
	filter := bson.D{{"_id", fleetId}}
	update := bson.D{{"$addToSet", bson.D{{"vins", vin}}}}
	result, err := m.database.Collection(fleetCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.ErrFleetNotFound
	}
	if result.ModifiedCount == 0 {
		return errors.ErrCarAlreadyInFleet
	}
	return nil
}

func (m *connection) GetCarsForFleet(ctx context.Context, fleetId model.FleetID) ([]model.Vin, error) {
	var fleet fleet
	err := m.database.Collection(fleetCollectionName).FindOne(ctx, bson.D{{"_id", fleetId}}).Decode(&fleet)
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrFleetNotFound
	}
	return fleet.Vins, err
}

func (m *connection) IsCarInFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) (bool, error) {
	vins, err := m.GetCarsForFleet(ctx, fleetId)
	if err != nil {
		return false, err
	}
	for _, foundVin := range vins {
		if foundVin == vin {
			return true, nil
		}
	}
	return false, nil
}

func (m *connection) RemoveCarFromFleet(ctx context.Context, fleetId model.FleetID, vin model.Vin) error {
	filter := bson.D{{"_id", fleetId}}
	update := bson.D{{"$pullAll", bson.D{{"vins", bson.A{vin}}}}}
	result, err := m.database.Collection(fleetCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.ErrFleetNotFound
	}
	if result.ModifiedCount == 0 {
		return errors.ErrCarNotInFleet
	}
	return nil
}
