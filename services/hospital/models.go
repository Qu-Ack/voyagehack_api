package hospital

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hospital struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"ID"`
	Participants Participants       `bson:"participants,omitempty" json:"participants"`
}

type Participants struct {
	Roots   []primitive.ObjectID `bson:"roots,omitempty" json:"roots"`
	Doctors []primitive.ObjectID `bson:"doctors,omitempty" json:"doctors"`
	Staff   []primitive.ObjectID `bson:"staff,omitempty" json:"staff"`
}
