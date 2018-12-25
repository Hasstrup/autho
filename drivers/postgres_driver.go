package drivers

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func PingPostgres(address string) error {
	client, _ := sql.Open("postgres", address)
	err := client.Ping()
	defer Recover()
	if err != nil {
		// Uncomment this line to view the errors connecting to your database
		// log.Println(err.Error())
		return errors.New(`Sorry we could not connect to the postgres url provided.
			The url should match this format 'user=pqgotest dbname=pqgotest sslmode=verify-full' 
			or some long ass string that looks like this 'postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full'
		`)
	}
	return nil
}

func YieldDrivers(database string) func(string) error {
	var fields = map[string]func(string) error{
		"postgres": PingPostgres,
		"mongodb":  PingMongo,
	}
	if function, ok := fields[database]; ok {
		return function
	}
	return func(string) error {
		return errors.New("Sorry we could not ping the address provided")
	}
}
