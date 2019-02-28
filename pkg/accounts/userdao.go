package accounts

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"bitbucket.org/imlabs/go_auth_service/pkg/db"
)

// UserModel is
type UserModel struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `bson:"Username" json:"Username"`
	Password string        `bson:"Password" json:"Password"`
}

func userModelIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"Username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// COLLECTION is
const (
	COLLECTION = "user"
)

// Insert is used to insert a User record in the database
func Insert(user UserModel) error {
	sess := db.Session.Clone()
	defer sess.Close()
	
	dataBase := sess.DB(db.ConnectionParams.Database)
	// Now we use sess to execute the query:
	
	err := dataBase.C(COLLECTION).Insert(&user)
	return err
}
