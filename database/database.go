package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/jpillora/backoff"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DB represents the structure of our database
type DB struct {
	Name    string
	Session *mgo.Session
	Uri     string
}

// NewDB creates a DB instance
func NewDB(uri string) DB {
	return DB{
		Name: viper.GetString("app-name"),
		Uri:  uri,
	}
}

// Dial connects to the server
func (db *DB) Dial() {
	b := &backoff.Backoff{
		Jitter: true,
	}

	for {
		session, err := mgo.Dial(db.Uri)
		db.Session = session

		if err != nil {
			d := b.Duration()
			logrus.Errorf("%s, reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}

		b.Reset()

		logrus.Info("Successfully connected to MongoDB")

		db.Session.SetSocketTimeout(time.Second * 3)
		db.Session.SetSyncTimeout(time.Second * 3)
	}
}

// Init initializes the database and creates indexes
func (db *DB) Init() error {
	err := db.createIndexes()
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) createIndexes() error {
	session := db.Session.Copy()
	d := session.DB(db.Name)
	defer session.Close()

	c := d.C("users")
	if c == nil {
		e := fmt.Sprint("Error creating collection")
		logrus.Error(e)
		return errors.New(e)
	}

	usersIndex := mgo.Index{
		Key:      []string{"email"},
		Unique:   true,
		DropDups: true,
	}
	err := c.EnsureIndex(usersIndex)
	if err != nil {
		e := fmt.Sprint("Error creating index")
		logrus.Error(e)
		return errors.New(e)
	}

	return nil

}
