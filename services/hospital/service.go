package hospital

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HospitalService struct {
	repo *HospitalRepo
}

func NewHospitalService(repo *HospitalRepo) *HospitalService {
	return &HospitalService{
		repo: repo,
	}
}

func (h *HospitalService) GetHospital(ctx context.Context, hospitalId string) (*Hospital, error) {
	hospitalObjectId, err := primitive.ObjectIDFromHex(hospitalId)

	if err != nil {
		return nil, err
	}

	return h.repo.GetHospitalFromId(ctx, hospitalObjectId)
}

func (h *HospitalService) AddParticipant(ctx context.Context, participantType string, participantId string, hospitalId string) (*Hospital, error) {
	hospitalObjectId, err := primitive.ObjectIDFromHex(hospitalId)
	if err != nil {
		return nil, err
	}

	participantObjectId, err := primitive.ObjectIDFromHex(participantId)

	if err != nil {
		return nil, err
	}

	switch participantType {
	case "ROOT":
		return h.repo.AddRootToHospital(ctx, hospitalObjectId, participantObjectId)
	case "STAFF":
		return h.repo.addStaffToHospital(ctx, hospitalObjectId, participantObjectId)
	case "DOCTOR":
		return h.repo.addDoctorToHospital(ctx, hospitalObjectId, participantObjectId)
	default:
		return nil, errors.New("should be a valid participant type")
	}

}

func (h *HospitalService) CheckParticpant(ctx context.Context, participantType string, participantId string) (*Hospital, error) {

	participantObjectId, err := primitive.ObjectIDFromHex(participantId)

	if err != nil {
		return nil, err
	}

	switch participantType {
	case "ROOT":
		hospital, err := h.repo.checkRootUser(ctx, participantObjectId)
		if err != nil {
			return nil, err
		}
		return hospital, nil
	case "STAFF":
		hospital, err := h.repo.checkStaffUser(ctx, participantObjectId)
		if err != nil {
			return nil, err
		}
		return hospital, nil
	case "DOCTOR":
		hospital, err := h.repo.checkDoctorUser(ctx, participantObjectId)
		if err != nil {
			return nil, err
		}
		return hospital, nil
	default:
		return nil, errors.New("participant Type should be Valid")

	}
}
