package db

import (
	"fmt"

	"bitbucket.org/imlabs/go_auth_service/pkg/logger"

	mgo "github.com/globalsign/mgo"
)

// ConnectionParameters is
type ConnectionParameters struct {
	Server   string
	Database string
}

// ConnectionParams is
var ConnectionParams = ConnectionParameters{
	Server:   "localhost:27017",
	Database: "AuthService",
}

// Username Password
const (
	Username = "YOUR_USERNAME"
	Password = "YOUR_PASS"
)

var log = logger.Logger

// Session is
var Session *mgo.Session

func init() {
	fmt.Println("Inside mgoconnection init() *************")
	ConnectionParams.Connect()
}

// Connect is
func (connection *ConnectionParameters) Connect() {
	var err error
	Session, err = mgo.Dial(connection.Server)
	if err != nil {
		log.Panic(err)
	}
}
