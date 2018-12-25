package middlewares

import (
	"github.com/authenticate/drivers"
	utils "github.com/authenticate/utilities"
)

func PingDatabaseAddress(database, address string, ch chan interface{}, counter *int) {
	if !utils.Contains(database, PermittedDatabaseTypes) {
		ch <- "Sorry the database has to be 'postgres' or 'mongodb' :("
	} else {
		if err := drivers.YieldDrivers(database)(address); err != nil {
			ch <- err.Error()
		}
	}
	*counter++
	ch <- true
}
