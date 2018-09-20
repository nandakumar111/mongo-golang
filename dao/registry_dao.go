package dao

import (
	"log"

	. "../models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	Server   string
	Database string
}

var db *mgo.Database

// Establish a connection to database
func (m *MongoDB) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Find list of channels
func (m *MongoDB) FindAll(COLLECTION string,query interface{},skip int,limit int,sortAttr string) ([]ChannelcodeModel, error) {
	var channelCodes []ChannelcodeModel
	err := db.C(COLLECTION).Find(query).Sort(sortAttr).Skip(skip).Limit(limit).Select(bson.M{"_id":0}).All(&channelCodes)
	return channelCodes, err
}

// Insert a channels into database
func (m *MongoDB) Insert(COLLECTION string,channelCodes ChannelcodeModel) error {
	err := db.C(COLLECTION).Insert(&channelCodes)
	return err
}
