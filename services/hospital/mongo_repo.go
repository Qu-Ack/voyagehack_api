package hospital

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HospitalRepo struct {
	db *mongo.Database
}

func NewHospitalRepo(db *mongo.Database) *HospitalRepo {
	return &HospitalRepo{
		db: db,
	}
}

func (h *HospitalRepo) GetHospitalFromId(ctx context.Context, hospitalId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	var hospital Hospital
	err := hospitalCollection.FindOne(ctx, bson.M{"_id": hospitalId}).Decode(&hospital)

	if err != nil {
		return nil, err
	}

	return &hospital, nil
}

func (h *HospitalRepo) AddRootToHospital(ctx context.Context, hospitalId primitive.ObjectID, rootId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	var hospital Hospital

	err := hospitalCollection.FindOneAndUpdate(ctx, bson.M{"_id": hospitalId}, bson.D{{"$push", bson.D{{"participants", bson.D{{"roots", rootId}}}}}}).Decode(&hospital)

	if err != nil {
		return nil, err
	}

	return &hospital, nil
}

func (h *HospitalRepo) addStaffToHospital(ctx context.Context, hospitalId primitive.ObjectID, staffId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	var hospital Hospital

	err := hospitalCollection.FindOneAndUpdate(ctx, bson.M{"_id": hospitalId}, bson.D{{"$push", bson.D{{"participants", bson.D{{"staff", staffId}}}}}}).Decode(&hospital)

	if err != nil {
		return nil, err
	}

	return &hospital, nil

}

func (h *HospitalRepo) addDoctorToHospital(ctx context.Context, hospitalId primitive.ObjectID, doctorId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	var hospital Hospital

	err := hospitalCollection.FindOneAndUpdate(ctx, bson.M{"_id": hospitalId}, bson.D{{"$push", bson.D{{"participants", bson.D{{"doctors", doctorId}}}}}}).Decode(&hospital)

	if err != nil {
		return nil, err
	}

	return &hospital, nil

}

func (h *HospitalRepo) checkRootUser(ctx context.Context, participantId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	filter := bson.M{
		"participants.roots": participantId,
	}

	var hospital Hospital
	err := hospitalCollection.FindOne(ctx, filter).Decode(&hospital)
	if err != nil {
		return nil, err
	}

	return &hospital, nil
}

func (h *HospitalRepo) checkStaffUser(ctx context.Context, participantId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	filter := bson.M{
		"participants.staff": participantId,
	}

	var hospital Hospital
	err := hospitalCollection.FindOne(ctx, filter).Decode(&hospital)
	if err != nil {
		return nil, err
	}

	return &hospital, nil
}

func (h *HospitalRepo) checkDoctorUser(ctx context.Context, participantId primitive.ObjectID) (*Hospital, error) {
	hospitalCollection := h.db.Collection("hospitals")

	filter := bson.M{
		"participants.doctors": participantId,
	}

	var hospital Hospital
	err := hospitalCollection.FindOne(ctx, filter).Decode(&hospital)
	if err != nil {
		return nil, err
	}

	return &hospital, nil
}
