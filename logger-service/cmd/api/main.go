package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Application struct {
	Models data.Models
}

func main() {
	// connect to mongo

	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Application{
		Models: data.New(client),
	}

	err = rpc.Register(new(RPCServer))
	if err != nil {
		fmt.Println("RPC server register failed")
	}
	go app.rpcListne()

	// start web serve
	log.Println("Server started at port ", webPort)

	// go server()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic()
	}
}

func (app *Application) rpcListne() error {
	fmt.Println("Starting RPC server on port", rpcPort)
	listne, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listne.Close()

	for {
		rpcConn, err := listne.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "root",
		Password: "renhaoran001",
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connect mongo", err)
		return nil, err
	}
	log.Println("Mongo connected...")
	return c, nil
}
