package main

import (
	"context"
	"log"
	"log-service/data"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogoEntry{
		ID:        primitive.NewObjectID(),
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}
	*resp = "Processed payload vis RPC" + payload.Name
	return nil
}
