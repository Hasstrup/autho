package middlewares

func PingDatabaseAddress(database, address string, ch chan interface{}, counter *int) {
	if database != "postgres" || database != "mongodb" {
		ch <- "Sorry the database has to be 'postgres' or 'mongodb' :("
		*counter++
		ch <- true
		return
	}
}
