
package config

import (
	"os"

	mgo "gopkg.in/mgo.v2"
)

func GetMongoDB() (*mgo.Database, error) {

	host := os.Getenv("MONGO_HOST")
	dbName := os.Getenv("MONGO_DB_NAME")

	//get host
	session, err := mgo.Dial(host)

	if err != nil {
		return nil, err
	}

	//obtain session
	db := session.DB(dbName)

	return db, nil
}
