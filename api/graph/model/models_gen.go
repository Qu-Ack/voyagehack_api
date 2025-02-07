// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AddressInformation struct {
	StreetAddress *string `json:"streetAddress,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	PinCode       *string `json:"pinCode,omitempty"`
}

type Amenities struct {
	BedCapacity    *BedCapacity  `json:"bedCapacity,omitempty"`
	MedicalStaff   *MedicalStaff `json:"medicalStaff,omitempty"`
	Facilities     []*string     `json:"facilities,omitempty"`
	Specialization []*string     `json:"specialization,omitempty"`
}

type Application struct {
	ID             string    `json:"id"`
	Content        string    `json:"content"`
	PatientName    string    `json:"patientName"`
	Sender         string    `json:"sender"`
	Receiver       string    `json:"receiver"`
	PatientAge     string    `json:"patientAge"`
	Documents      []*string `json:"documents,omitempty"`
	PatientGender  string    `json:"patientGender"`
	Passport       string    `json:"passport"`
	PhoneNumber    string    `json:"phoneNumber"`
	Allergies      string    `json:"allergies"`
	Type           EmailType `json:"type"`
	ForwardedChain []string  `json:"ForwardedChain"`
	CreatedAt      string    `json:"createdAt"`
}

type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type BasicInfo struct {
	HospitalName       *string             `json:"hospitalName,omitempty"`
	RegistrationNumber *string             `json:"registrationNumber,omitempty"`
	ContactInformation *ContactInformation `json:"contactInformation,omitempty"`
	AddressInformation *AddressInformation `json:"addressInformation,omitempty"`
	OperatingHours     *OperatingHours     `json:"operatingHours,omitempty"`
}

type BedCapacity struct {
	GeneralWardBeds *int32 `json:"generalWardBeds,omitempty"`
	PrivateRoomBeds *int32 `json:"privateRoomBeds,omitempty"`
	EmergencyBeds   *int32 `json:"emergencyBeds,omitempty"`
	IcuBeds         *int32 `json:"icuBeds,omitempty"`
}

type ContactInformation struct {
	ContactPersonName *string `json:"contactPersonName,omitempty"`
	ContactNumber     *string `json:"contactNumber,omitempty"`
	ContactEmail      *string `json:"contactEmail,omitempty"`
	Website           *string `json:"website,omitempty"`
}

type Doctor struct {
	User      *User    `json:"user"`
	Specialty string   `json:"specialty"`
	Documents []string `json:"documents"`
}

type DoctorInput struct {
	Specialty  string   `json:"specialty"`
	Documents  []string `json:"documents"`
	Name       string   `json:"name"`
	Password   string   `json:"password"`
	Email      string   `json:"email"`
	HospitalID string   `json:"hospitalId"`
	ProfilePic string   `json:"profilePic"`
}

type Hospital struct {
	ID                 string              `json:"id"`
	BasicInfo          *BasicInfo          `json:"basicInfo,omitempty"`
	Media              *Media              `json:"media,omitempty"`
	Amenities          *Amenities          `json:"amenities,omitempty"`
	OnSiteVerification *OnSiteVerification `json:"onSiteVerification,omitempty"`
	OnsiteRating       *int32              `json:"onsiteRating,omitempty"`
	Reviews            []*Review           `json:"reviews"`
	PatientRating      *int32              `json:"patientRating,omitempty"`
	Ratings            []int32             `json:"ratings"`
	ConsultationFee    *int32              `json:"consultationFee,omitempty"`
	Participants       *Participants       `json:"participants,omitempty"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Mail struct {
	ID             string    `json:"id"`
	Sender         string    `json:"sender"`
	Receiver       string    `json:"receiver"`
	Content        string    `json:"content"`
	Documents      []*string `json:"documents"`
	Type           EmailType `json:"type"`
	ForwardedChain []string  `json:"ForwardedChain"`
	CreatedAt      string    `json:"createdAt"`
}

type MailBox struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	Sent           []*string `json:"sent"`
	Received       []*string `json:"received"`
	Sentmails      []*Mail   `json:"sentmails"`
	ReceivedEmails []*Mail   `json:"receivedEmails"`
}

type MailBoxSubscriptionResponse struct {
	Sent     *Application `json:"sent,omitempty"`
	Received *Application `json:"received,omitempty"`
}

type Media struct {
	FrontURL     *string `json:"frontUrl,omitempty"`
	ReceptionURL *string `json:"receptionUrl,omitempty"`
	OperationURL *string `json:"operationUrl,omitempty"`
}

type MedicalStaff struct {
	PermenantDoctors    *int32 `json:"permenantDoctors,omitempty"`
	VisitingConsultants *int32 `json:"visitingConsultants,omitempty"`
	Nurses              *int32 `json:"nurses,omitempty"`
	SupportStaff        *int32 `json:"supportStaff,omitempty"`
}

type Message struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type MessageSubscriptionResponse struct {
	RoomID  string   `json:"roomId"`
	Message *Message `json:"message"`
}

type Mutation struct {
}

type OnSiteVerification struct {
	PreferredDate       *string              `json:"preferredDate,omitempty"`
	PreferredTime       *string              `json:"preferredTime,omitempty"`
	VerificationContact *VerificationContact `json:"verificationContact,omitempty"`
}

type OperatingHours struct {
	OpeningTime *string `json:"openingTime,omitempty"`
	ClosingTime *string `json:"closingTime,omitempty"`
}

type Participants struct {
	Roots   []string `json:"roots"`
	Staff   []string `json:"staff"`
	Doctors []string `json:"doctors"`
}

type PatientInput struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Password   string `json:"password"`
}

type Query struct {
}

type Review struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Room struct {
	ID           string     `json:"id"`
	Messages     []*Message `json:"messages"`
	Participants []string   `json:"participants"`
	State        string     `json:"state"`
}

type Subscription struct {
}

type TestUserInput struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Password   string `json:"password"`
	Role       string `json:"role"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Role       Role   `json:"role"`
}

type UserInput struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
	Password   string `json:"password"`
	HospitalID string `json:"hospitalId"`
}

type VerificationContact struct {
	Name           *string `json:"name,omitempty"`
	Position       *string `json:"position,omitempty"`
	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	AlternatePhone *string `json:"alternatePhone,omitempty"`
}

type SendInvitationInput struct {
	Receiver  string   `json:"receiver"`
	Content   string   `json:"content"`
	Documents []string `json:"documents"`
}

type SendMailInput struct {
	HospitalID        string    `json:"hospitalId"`
	Content           string    `json:"content"`
	PatientName       string    `json:"patientName"`
	PatientAge        string    `json:"patientAge"`
	PatientGender     string    `json:"patientGender"`
	PhoneNumber       string    `json:"phoneNumber"`
	Passport          string    `json:"passport"`
	Allergies         string    `json:"allergies"`
	Documents         []*string `json:"documents,omitempty"`
	RazorpayPaymentID string    `json:"razorpay_payment_id"`
	RazorpayOrderID   string    `json:"razorpay_order_id"`
	RazorpaySignature string    `json:"razorpay_signature"`
}

type SendMessageInput struct {
	Content string `json:"content"`
	RoomID  string `json:"roomId"`
}

type SendNormalMailInput struct {
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	Documents []*string `json:"documents,omitempty"`
}

type EmailType string

const (
	EmailTypeApplication EmailType = "APPLICATION"
	EmailTypeNormal      EmailType = "NORMAL"
)

var AllEmailType = []EmailType{
	EmailTypeApplication,
	EmailTypeNormal,
}

func (e EmailType) IsValid() bool {
	switch e {
	case EmailTypeApplication, EmailTypeNormal:
		return true
	}
	return false
}

func (e EmailType) String() string {
	return string(e)
}

func (e *EmailType) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EmailType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EmailType", str)
	}
	return nil
}

func (e EmailType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Role string

const (
	RoleRoot   Role = "ROOT"
	RoleDoctor Role = "DOCTOR"
	RoleStaff  Role = "STAFF"
)

var AllRole = []Role{
	RoleRoot,
	RoleDoctor,
	RoleStaff,
}

func (e Role) IsValid() bool {
	switch e {
	case RoleRoot, RoleDoctor, RoleStaff:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalGQL(v any) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
