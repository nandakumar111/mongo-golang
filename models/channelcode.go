package models

import "gopkg.in/mgo.v2/bson"

// Represents a user, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type ChannelcodeModel struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
	Channelcode string `bson:"channelcode" json:"channelcode"`
	Url string `bson:"url" json:"url"`
}
