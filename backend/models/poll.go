package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PollOption struct {
	ID   string `bson:"id" json:"id"`
	Text string `bson:"text" json:"text"`
}

type Poll struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question  string             `bson:"question" json:"question"`
	Options   []PollOption       `bson:"options" json:"options"`
	CreatedBy primitive.ObjectID `bson:"createdBy" json:"createdBy"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
	IsClosed  bool               `bson:"isClosed" json:"isClosed"`
}

type Vote struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PollID    primitive.ObjectID `bson:"pollId" json:"pollId"`
	OptionID  string             `bson:"optionId" json:"optionId"`
	VoterID   string             `bson:"voterId" json:"voterId"`
	CreatedAt int64              `bson:"createdAt" json:"createdAt"`
}
