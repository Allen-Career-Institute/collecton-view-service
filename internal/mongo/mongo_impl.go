package mongo

import (
	"context"
	//"fmt"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"time"

	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	START_TIME_LOG = "Start Time to Fetch Doc in mongo = %v"
)

type mongoCollectionImpl struct {
	client *mongo.Client
	log    *log.Helper
}

func (q mongoCollectionImpl) Get(ctx context.Context, filter interface{}, dbName, collectionName string) (*mongo.SingleResult, error) {
	// Start a new session for the transaction
	q.log.WithContext(ctx).Debugf(START_TIME_LOG, time.Now())
	session, err := q.client.StartSession()
	defer session.EndSession(ctx)
	// Start the transaction
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	// Get a handle to the database and collection
	db := q.client.Database(dbName)
	collection := db.Collection(collectionName)

	// Query the collection
	result := collection.FindOne(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("No document found")
			return nil, nil
		}
		// Abort the transaction in case of an error
		err = session.AbortTransaction(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	return result, nil
}

func (q mongoCollectionImpl) List(ctx context.Context, filter interface{}, dbName, collectionName string, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	// Start a new session for the transaction
	q.log.WithContext(ctx).Debugf(START_TIME_LOG, time.Now())
	session, err := q.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	// Start the transaction
	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}
	// Get a handle to the database and collection
	db := q.client.Database(dbName)
	collection := db.Collection(collectionName)

	// Query the collection
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		// Abort the transaction in case of an error
		err = session.AbortTransaction(ctx)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	if err = cursor.Err(); err != nil {
		q.log.WithContext(ctx).Errorf("Error during cursor iteration = %v", err)
	}
	return cursor, err
}

func NewMongoCollectionImpl(data *Data, logger *log.Helper) MongoCollectionInterface {
	return &mongoCollectionImpl{data.client, logger}
}