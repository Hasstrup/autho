package drivers

import (
	"context"
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

func PingMongo(address string) error {
	client, err := mongo.NewClient(address)
	if err != nil {
		return errors.New("Please check the mongo string provided should match 'mongodb://therestofyourmongostring' ")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	defer Recover()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		// TODO: You need to extract all these error messages somewhere
		return errors.New("We tried pinging the mongo address you provided, didn't work out")
	}
	return nil
}

func Recover() error {
	if r := recover(); r != nil {
		return errors.New("Sorry I could not process this request, check the input sent and try again")
	}
	return nil
}
