package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogoEntry: LogoEntry{},
	}
}

type Models struct {
	LogoEntry LogoEntry
}

type LogoEntry struct {
	ID        string    `bson:"_id,moitepmty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogoEntry) Insert(entry LogoEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogoEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}
	return nil
}

func (l *LogoEntry) All() ([]*LogoEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()

	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Finding all docs err", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var Logs []*LogoEntry

	for cursor.Next(ctx) {
		var item LogoEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice", err)
			return nil, err
		} else {
			Logs = append(Logs, &item)
		}
	}
	return Logs, nil
}

func (l *LogoEntry) GetOne(id string) (*LogoEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")
	docID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var entry LogoEntry
	err = collection.FindOne(ctx, bson.M{"_id": docID}).Decode(&entry)

	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (l *LogoEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}

func (l *LogoEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docID, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": docID}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
