package database

import (
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
)

// CreateUser creates a new user
func (db *DB) CreateUser(data interface{}) error {
	session := db.Session.Copy()
	d := session.DB(db.Name)
	defer session.Close()

	err := d.C("users").Insert(data)

	if err != nil {
		if !mgo.IsDup(err) {
			e := fmt.Sprintf("Error creating user: %v", err)
			logrus.Error(e)
		}
		return err
	}

	return nil
}

// DeleteUser deletes a user from an id
func (db *DB) DeleteUser(id string) error {
	session := db.Session.Copy()
	d := session.DB(db.Name)
	defer session.Close()

	err := d.C("users").RemoveId(bson.ObjectIdHex(id))
	if err != nil && err != mgo.ErrNotFound {
		e := fmt.Sprintf("Error deleting user: %s", err)
		logrus.Error(e)
	} else if err == mgo.ErrNotFound {
		return err
	}

	return nil
}

func (db *DB) query() error {
	return nil
}

func (db *DB) findUser(m bson.M, data interface{}) error {
	session := db.Session.Copy()
	d := session.DB(db.Name)
	defer session.Close()

	err := d.C("users").Find(m).One(data)

	if err != nil && err != mgo.ErrNotFound {
		e := fmt.Sprintf("Error finding user: %s", err)
		logrus.Error(e)
	} else if err == mgo.ErrNotFound {
		return err
	}

	return nil
}

// FindUserByEmail finds a user based on an email address
func (db *DB) FindUserByEmail(email string, data interface{}) error {
	m := bson.M{"email": email}
	return db.findUser(m, data)
}

// FindUserById finds a user from an id
func (db *DB) FindUserById(id string, data interface{}) error {
	m := bson.M{"_id": bson.ObjectIdHex(id)}
	return db.findUser(m, data)
}

// GetAllUsers returns all the users
func (db *DB) GetAllUsers(data interface{}) error {
	session := db.Session.Copy()
	d := session.DB(db.Name)
	defer session.Close()

	err := d.C("users").Find(nil).All(data)
	if err != nil && err != mgo.ErrNotFound {
		e := fmt.Sprintf("Error finding users: %s", err)
		logrus.Error(e)
	} else if err == mgo.ErrNotFound {
		return err
	}

	return nil
}
