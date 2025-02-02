package hospital

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hospital struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	BasicInfo          BasicInfo          `bson:"BasicInfo,omitempty" json:"BasicInfo"`
	Media              Media              `bson:"Media,omitempty" json:"Media"`
	Amenities          Amenities          `bson:"Amenities,omitempty" json:"Amenities"`
	OnSiteVerification OnSiteVerification `bson:"OnSiteVerification,omitempty" json:"OnSiteVerification"`
	OnsiteRating       int32              `bson:"OnsiteRating,omitempty" json:"OnsiteRating"`
	Reviews            []Review           `bson:"Reviews,omitempty" json:"Reviews"`
	PatientRating      int32              `bson:"PatientRating,omitempty" json:"PatientRating"`
	Ratings            []int32            `bson:"Ratings,omitempty" json:"Ratings"`
	ConsultationFee    int32              `bson:"ConsultationFee,omitempty" json:"ConsultationFee"`
	Participants       Participants       `bson:"Participants,omitempty" json:"Participants"`
	Version            int32              `bson:"__v,omitempty" json:"__v"`
}

type Review struct {
	Content string             `bson:"Content,omitempty" json:"content"`
	Author  primitive.ObjectID `bson:"Author,omitempty" json:"author"`
}

type BasicInfo struct {
	HospitalName       string             `bson:"HospitalName,omitempty" json:"HospitalName"`
	RegistrationNumber string             `bson:"RegistrationNumber,omitempty" json:"RegistrationNumber"`
	ContactInformation ContactInformation `bson:"ContactInformation,omitempty" json:"ContactInformation"`
	AddressInformation AddressInformation `bson:"AddressInformation,omitempty" json:"AddressInformation"`
	OperatingHours     OperatingHours     `bson:"OperatingHours,omitempty" json:"OperatingHours"`
}

type ContactInformation struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ContactPersonName string             `bson:"ContactPersonName,omitempty" json:"ContactPersonName"`
	ContactNumber     string             `bson:"ContactNumber,omitempty" json:"ContactNumber"`
	ContactEmail      string             `bson:"ContactEmail,omitempty" json:"ContactEmail"`
	Website           string             `bson:"Website,omitempty" json:"Website"`
}

type AddressInformation struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	StreetAddress string             `bson:"StreetAddress,omitempty" json:"StreetAddress"`
	City          string             `bson:"City,omitempty" json:"City"`
	State         string             `bson:"State,omitempty" json:"State"`
	PinCode       string             `bson:"PinCode,omitempty" json:"PinCode"`
}

type OperatingHours struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	OpeningTime string             `bson:"OpeningTime,omitempty" json:"OpeningTime"`
	ClosingTime string             `bson:"ClosingTime,omitempty" json:"ClosingTime"`
}

type Media struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	FrontUrl     string             `bson:"FrontUrl,omitempty" json:"FrontUrl"`
	ReceptionUrl string             `bson:"ReceptionUrl,omitempty" json:"ReceptionUrl"`
	OperationUrl string             `bson:"OperationUrl,omitempty" json:"OperationUrl"`
}

type Amenities struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	BedCapacity    BedCapacity        `bson:"BedCapacity,omitempty" json:"BedCapacity"`
	MedicalStaff   MedicalStaff       `bson:"MedicalStaff,omitempty" json:"MedicalStaff"`
	Facilities     []string           `bson:"Facilities,omitempty" json:"Facilities"`
	Specialization []string           `bson:"Specialization,omitempty" json:"Specialization"`
}

type BedCapacity struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	GeneralWardBeds int32              `bson:"GeneralWardBeds,omitempty" json:"GeneralWardBeds"`
	PrivateRoomBeds int32              `bson:"PrivateRoomBeds,omitempty" json:"PrivateRoomBeds"`
	EmergencyBeds   int32              `bson:"EmergencyBeds,omitempty" json:"EmergencyBeds"`
	IcuBeds         int32              `bson:"IcuBeds,omitempty" json:"IcuBeds"`
}

type MedicalStaff struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	PermenantDoctors    int32              `bson:"PermenantDoctors,omitempty" json:"PermenantDoctors"`
	VisitingConsultants int32              `bson:"VisitingConsultants,omitempty" json:"VisitingConsultants"`
	Nurses              int32              `bson:"Nurses,omitempty" json:"Nurses"`
	SupportStaff        int32              `bson:"SupportStaff,omitempty" json:"SupportStaff"`
}

type OnSiteVerification struct {
	ID                  primitive.ObjectID  `bson:"_id,omitempty" json:"_id"`
	PreferredDate       primitive.DateTime  `bson:"PreferredDate,omitempty" json:"PreferredDate"`
	PreferredTime       string              `bson:"PreferredTime,omitempty" json:"PreferredTime"`
	VerificationContact VerificationContact `bson:"VerificationContact,omitempty" json:"VerificationContact"`
}

type VerificationContact struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name           string             `bson:"Name,omitempty" json:"Name"`
	Position       string             `bson:"Position,omitempty" json:"Position"`
	PhoneNumber    string             `bson:"PhoneNumber,omitempty" json:"PhoneNumber"`
	AlternatePhone string             `bson:"AlternatePhone,omitempty" json:"AlternatePhone"`
}

type Participants struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"_id"`
	RootUsers []primitive.ObjectID `bson:"RootUsers,omitempty" json:"RootUsers"`
	Staff     []primitive.ObjectID `bson:"Staff,omitempty" json:"Staff"`
	Doctors   []primitive.ObjectID `bson:"Doctors,omitempty" json:"Doctors"`
}
