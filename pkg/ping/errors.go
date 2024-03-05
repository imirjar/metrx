package ping

import "errors"

var (
	errConnectionString = errors.New("connection string is incorrect")
	errConnectionParams = errors.New("couldn't connect to db with matched params")
)
