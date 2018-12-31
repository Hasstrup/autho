package drivers

import (
	"context"
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
)

func PingMongo(address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	_, err := mongo.Connect(ctx, address)
	if err != nil {
		return errors.New("We tried connecting to the mongostring provided, did not work out :(. The mongo string provided should match 'mongodb://therestofyourmongostring' ")
	}
	defer cancel()
	defer Recover()
	return nil
}

func Recover() error {
	if r := recover(); r != nil {
		return errors.New("Sorry I could not process this request, check the input sent and try again")
	}
	return nil
}
