package repository

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	log "github.com/sirupsen/logrus"

	"github.com/Businge931/hexagonal-architecture/domain"
)

type mongoRepository struct {
	client  *mongo.Client
	db      string
	timeout time.Duration
}

// Create a Mongo Client
func newMongoClient(mongoServerURL string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoServerURL))
	if err != nil {
		return nil, err
	}
	// We could ping the server to test connectivity if we want

	return client, nil
}

func NewMongoRepository(serverURL, dB string, timeout int) (domain.Repository, error) {
	mongoClient, err := newMongoClient(serverURL, timeout)
	repo := &mongoRepository{
		client: mongoClient, db: dB, timeout: time.Duration(timeout) * time.Second,
	}
	if err != nil {
		return nil, errors.Wrap(err, "client error")
	}

	return repo, nil
}

func (r *mongoRepository) Store(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.db).Collection("items")

	_, err := collection.InsertOne(ctx, bson.M{"code": product.Code, "name": product.Name, "price": product.Price})
	if err != nil {
		return errors.Wrap(err, "Error writing to repository")
	}
	return nil
}

func (r *mongoRepository) Update(product *domain.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	collection := r.client.Database(r.db).Collection("items")
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"code": product.Code},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "name", Value: product.Name}, {Key: "price", Value: product.Price}}},
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *mongoRepository) Find(code string) (*domain.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	product := &domain.Product{}
	collection := r.client.Database(r.db).Collection("items")
	filter := bson.M{"code": code}

	err := collection.FindOne(ctx, filter).Decode(product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("Error Finding a catalogue item")
		}
		return nil, errors.Wrap(err, "repository research")
	}

	return product, nil

}

func (r *mongoRepository) FindAll() ([]*domain.Product, error) {

	var items []*domain.Product

	collection := r.client.Database(r.db).Collection("items")
	cur, err := collection.Find(context.Background(), bson.D{})

	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.TODO()) {

		var item domain.Product
		if err := cur.Decode(&item); err != nil {
			log.Fatal(err)
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil

}

func (r *mongoRepository) Delete(code string) error {

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	filter := bson.M{"code": code}

	collection := r.client.Database(r.db).Collection("items")
	_, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	return nil

}