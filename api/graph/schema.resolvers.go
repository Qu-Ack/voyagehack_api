package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.63

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Qu-Ack/voyagehack_api/api/graph/model"
	"github.com/Qu-Ack/voyagehack_api/services/mail"
	"github.com/Qu-Ack/voyagehack_api/services/messaging"
	"github.com/Qu-Ack/voyagehack_api/services/observers"
	"github.com/Qu-Ack/voyagehack_api/services/payment"
	"github.com/Qu-Ack/voyagehack_api/services/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.LoginInput) (*model.AuthPayload, error) {
	authPayload, err := r.Resolver.UserService.Login(ctx, user.LoginInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		return nil, err
	}

	return &model.AuthPayload{
		Token: authPayload.Token,
		User: &model.User{
			ID:         authPayload.User.ID,
			Name:       authPayload.User.Name,
			Email:      authPayload.User.Email,
			ProfilePic: authPayload.User.ProfilePic,
			Role:       model.Role(authPayload.User.Role),
		},
	}, nil
}

// CreateDoctor is the resolver for the createDoctor field.
func (r *mutationResolver) CreateDoctor(ctx context.Context, input model.DoctorInput) (*model.Doctor, error) {
	panic(fmt.Errorf("not implemented: CreateDoctor - createDoctor"))
}

// CreateRoot is the resolver for the createRoot field.
func (r *mutationResolver) CreateRoot(ctx context.Context, input model.UserInput) (*model.User, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	publicUser, err := r.Resolver.UserService.CreateRoot(ctx, user.UserInput{
		Name:       input.Name,
		Email:      input.Email,
		ProfilePic: input.ProfilePic,
		Password:   input.Password,
	}, user.PublicUser{
		ID:    authedUser.ID,
		Email: authedUser.Email,
		Role:  user.Role(authedUser.Role),
	})

	if err != nil {
		return nil, err
	}

	_, err = r.HospitalService.AddParticipant(ctx, "ROOT", publicUser.ID, input.HospitalID)

	if err != nil {
		return nil, err
	}

	_, err = r.MailService.InitializeMailBox(ctx, publicUser.Email)
	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:         publicUser.ID,
		Email:      publicUser.Email,
		Role:       model.Role(publicUser.Role),
		ProfilePic: publicUser.ProfilePic,
		Name:       publicUser.Name,
	}, nil
}

// CreateStaff is the resolver for the createStaff field.
func (r *mutationResolver) CreateStaff(ctx context.Context, input model.UserInput) (*model.User, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	publicUser, err := r.Resolver.UserService.CreateStaff(ctx, user.UserInput{
		Name:       input.Name,
		Email:      input.Email,
		ProfilePic: input.ProfilePic,
		Password:   input.Password,
	}, user.PublicUser{
		ID:    authedUser.ID,
		Role:  user.Role(authedUser.Role),
		Email: authedUser.Email,
	})
	if err != nil {
		return nil, err
	}

	_, err = r.HospitalService.AddParticipant(ctx, "STAFF", publicUser.ID, input.HospitalID)

	if err != nil {
		return nil, err
	}

	_, err = r.MailService.InitializeMailBox(ctx, publicUser.Email)

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:         publicUser.ID,
		Name:       publicUser.Name,
		Email:      publicUser.Email,
		ProfilePic: publicUser.ProfilePic,
		Role:       model.Role(publicUser.Role),
	}, nil
}

// CreateTestUser is the resolver for the createTestUser field.
func (r *mutationResolver) CreateTestUser(ctx context.Context, input model.TestUserInput) (*model.User, error) {
	_, ok := ctx.Value(UserContextKey).(string)

	if !ok {
		return nil, errors.New("not a test user")
	}

	publicUser, err := r.Resolver.UserService.CreateTestUser(ctx, user.UserInput{
		Name:       input.Name,
		Email:      input.Email,
		ProfilePic: input.ProfilePic,
		Password:   input.Password,
		Role:       user.Role(input.Role),
	})

	if err != nil {
		return nil, err
	}

	_, err = r.MailService.InitializeMailBox(ctx, publicUser.Email)

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:         publicUser.ID,
		Name:       publicUser.Name,
		Email:      publicUser.Email,
		Role:       model.Role(publicUser.Role),
		ProfilePic: publicUser.ProfilePic,
	}, nil
}

// SendApplication is the resolver for the sendApplication field.
func (r *mutationResolver) SendApplication(ctx context.Context, input model.SendMailInput) (*model.Mail, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	err := r.PaymentService.ValidatePayment(&payment.ValidatePaymentRequest{
		RazorpayPaymentId: input.RazorpayPaymentID,
		RazorpayOrderId:   input.RazorpayOrderID,
		RazorpaySignature: input.RazorpaySignature,
	}, user.PublicUser{
		ID:    authedUser.ID,
		Email: authedUser.Email,
		Role:  user.Role(authedUser.Role),
	})

	if err != nil {
		return nil, err
	}

	hospital, err := r.HospitalService.GetHospital(ctx, input.HospitalID)

	if err != nil {
		return nil, err
	}

	receiverUserId := selectRandomParticipant(hospital.Participants.Roots)

	user, err := r.UserService.Me(ctx, receiverUserId.Hex())

	if err != nil {
		return nil, err
	}

	mail, err := r.MailService.SendApplication(ctx, &mail.Mail{
		Content:   input.Content,
		Sender:    authedUser.Email,
		Receiver:  user.Email,
		Documents: input.Documents,
		Type:      mail.Application,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	})

	if err != nil {
		return nil, err
	}

	r.ObserverService.PublishMail(user.Email, &observers.MailBoxSubscriptionResponse{
		Received: observers.Mail{
			ID:        mail.ID,
			Sender:    mail.Sender,
			Receiver:  mail.Receiver,
			Content:   mail.Content,
			Documents: mail.Documents,
			CreatedAt: mail.CreatedAt,
		},
	})
	r.ObserverService.PublishMail(authedUser.Email, &observers.MailBoxSubscriptionResponse{
		Sent: observers.Mail{
			ID:        mail.ID,
			Sender:    mail.Sender,
			Receiver:  mail.Receiver,
			Content:   mail.Content,
			Documents: mail.Documents,
			CreatedAt: mail.CreatedAt,
		},
	})

	return &model.Mail{
		ID:        mail.ID.Hex(),
		Sender:    mail.Sender,
		Receiver:  mail.Receiver,
		Content:   mail.Content,
		Documents: mail.Documents,
		CreatedAt: mail.CreatedAt.Time().String(),
	}, nil
}

// SendNormalMail is the resolver for the sendNormalMail field.
func (r *mutationResolver) SendNormalMail(ctx context.Context, input model.SendNormalMailInput) (*model.Mail, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	mail, err := r.MailService.SendNormalMail(ctx, &mail.Mail{
		Content:   input.Content,
		Sender:    authedUser.Email,
		Receiver:  input.Receiver,
		Documents: input.Documents,
		Type:      mail.Application,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}, user.PublicUser{
		ID:    authedUser.ID,
		Role:  user.Role(authedUser.Role),
		Email: authedUser.Email,
	})

	if err != nil {
		return nil, err
	}

	r.ObserverService.PublishMail(input.Receiver, &observers.MailBoxSubscriptionResponse{
		Received: observers.Mail{
			ID:        mail.ID,
			Sender:    mail.Sender,
			Receiver:  mail.Receiver,
			Content:   mail.Content,
			Documents: mail.Documents,
			CreatedAt: mail.CreatedAt,
		},
	})
	r.ObserverService.PublishMail(authedUser.Email, &observers.MailBoxSubscriptionResponse{
		Sent: observers.Mail{
			ID:        mail.ID,
			Sender:    mail.Sender,
			Receiver:  mail.Receiver,
			Content:   mail.Content,
			Documents: mail.Documents,
			CreatedAt: mail.CreatedAt,
		},
	})

	return &model.Mail{
		ID:        mail.ID.Hex(),
		Sender:    mail.Sender,
		Receiver:  mail.Receiver,
		Content:   mail.Content,
		Type:      model.EmailType(mail.Type),
		Documents: mail.Documents,
		CreatedAt: mail.CreatedAt.Time().String(),
	}, nil
}

// StartChat is the resolver for the startChat field.
func (r *mutationResolver) StartChat(ctx context.Context, participantID string) (*model.Room, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	room, err := r.MessagingService.CreateRoom(ctx, participantID, user.PublicUser{
		ID:    authedUser.ID,
		Role:  user.Role(authedUser.Role),
		Email: authedUser.Email,
	})

	if err != nil {
		return nil, err
	}

	return &model.Room{
		ID:           room.ID.Hex(),
		Messages:     convertMessages(room.Messages),
		Participants: convertParticipants(room.Participants),
	}, nil
}

// SendMessage is the resolver for the sendMessage field.
func (r *mutationResolver) SendMessage(ctx context.Context, input model.SendMessageInput) (*model.Room, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	room, err := r.MessagingService.SendMessage(ctx, input.RoomID, user.PublicUser{
		ID:    authedUser.ID,
		Email: authedUser.Email,
		Role:  user.Role(authedUser.Role),
	}, &messaging.Message{
		Content: input.Content,
		Author:  authedUser.Email,
	})

	if err != nil {
		return nil, err
	}

	return &model.Room{
		ID:           room.ID.Hex(),
		Messages:     convertMessages(room.Messages),
		Participants: convertParticipants(room.Participants),
	}, nil
}

// CloseRoom is the resolver for the closeRoom field.
func (r *mutationResolver) CloseRoom(ctx context.Context, roomid string) (*model.Room, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	room, err := r.MessagingService.ChangeRoomState(ctx, "CLOSE", user.PublicUser{
		ID:    authedUser.ID,
		Role:  user.Role(authedUser.Role),
		Email: authedUser.Email,
	}, roomid)

	if err != nil {
		return nil, err
	}

	return &model.Room{
		ID:           room.ID.Hex(),
		Messages:     convertMessages(room.Messages),
		Participants: convertParticipants(room.Participants),
	}, nil
}

// OpenRoom is the resolver for the openRoom field.
func (r *mutationResolver) OpenRoom(ctx context.Context, roomid string) (*model.Room, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	room, err := r.MessagingService.ChangeRoomState(ctx, "OPEN", user.PublicUser{
		ID:    authedUser.ID,
		Email: authedUser.Email,
		Role:  user.Role(authedUser.Role),
	}, roomid)

	if err != nil {
		return nil, err
	}

	return &model.Room{
		ID:           room.ID.Hex(),
		Messages:     convertMessages(room.Messages),
		Participants: convertParticipants(room.Participants),
	}, nil
}

// AddToHospital is the resolver for the addToHospital field.
func (r *mutationResolver) AddToHospital(ctx context.Context, userMail string, hospitalID string) (*model.Hospital, error) {
	panic("not implemented")
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	fmt.Println(authedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	publicUser, err := r.Resolver.UserService.Me(ctx, authedUser.ID)

	if err != nil {
		return nil, err
	}

	return &model.User{
		ID:    publicUser.ID,
		Email: publicUser.Email,
		Role:  model.Role(publicUser.Role),
		Name:  publicUser.Name,
	}, nil
}

// MyDoctorProfile is the resolver for the myDoctorProfile field.
func (r *queryResolver) MyDoctorProfile(ctx context.Context) (*model.Doctor, error) {
	panic(fmt.Errorf("not implemented: MyDoctorProfile - myDoctorProfile"))
}

// GetEmailByID is the resolver for the getEmailByID field.
func (r *queryResolver) GetEmailByID(ctx context.Context, id string) (*model.Mail, error) {
	_, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	mail, err := r.MailService.GetMail(ctx, id)

	if err != nil {
		return nil, err
	}

	return &model.Mail{
		ID:        mail.ID.Hex(),
		Content:   mail.Content,
		Sender:    mail.Sender,
		Receiver:  mail.Receiver,
		Documents: mail.Documents,
		Type:      model.EmailType(mail.Type),
		CreatedAt: mail.CreatedAt.Time().String(),
	}, nil
}

// GetMailBox is the resolver for the getMailBox field.
func (r *queryResolver) GetMailBox(ctx context.Context) (*model.MailBox, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	mailbox, err := r.MailService.GetMailBox(ctx, authedUser.Email)

	if err != nil {
		return nil, err
	}

	sentmails := convertMailSlice(mailbox.SentEMails)
	receivedmails := convertMailSlice(mailbox.ReceivedEmails)
	sentObjectids := convertObjectIDToStringSlice(mailbox.MailBox.Sent)
	receivedObjectids := convertObjectIDToStringSlice(mailbox.MailBox.Received)

	return &model.MailBox{
		Sentmails:      sentmails,
		ReceivedEmails: receivedmails,
		Sent:           sentObjectids,
		Received:       receivedObjectids,
		Email:          mailbox.MailBox.Email,
	}, nil
}

// GetRoom is the resolver for the getRoom field.
func (r *queryResolver) GetRoom(ctx context.Context, roomID string) (*model.Room, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	room, err := r.MessagingService.GetRoom(ctx, roomID, user.PublicUser{
		ID:    authedUser.ID,
		Email: authedUser.Email,
		Role:  user.Role(authedUser.Role),
	})
	if err != nil {
		return nil, err
	}
	return &model.Room{
		ID:           room.ID.Hex(),
		Messages:     convertMessages(room.Messages),
		Participants: convertParticipants(room.Participants),
	}, nil
}

// GetS3Url is the resolver for the getS3Url field.
func (r *queryResolver) GetS3Url(ctx context.Context) (string, error) {
	_, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return "", fmt.Errorf("unauthorized: user not found in context")
	}

	return r.UploadService.GetPresignedURL()
}

// GetOrderID is the resolver for the getOrderId field.
func (r *queryResolver) GetOrderID(ctx context.Context) (string, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return "", fmt.Errorf("unauthorized: user not found in context")
	}

	return r.PaymentService.NewOrder(&payment.OrderRequest{
		Amount: 100,
	}, user.PublicUser{
		ID:    authedUser.ID,
		Role:  user.Role(authedUser.Role),
		Email: authedUser.Email,
	})
}

// GetHospital is the resolver for the getHospital field.
func (r *queryResolver) GetHospital(ctx context.Context) (*model.Hospital, error) {
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	switch user.Role(authedUser.Role) {
	case user.RoleRoot:
		hospital, err := r.HospitalService.CheckParticpant(ctx, "ROOT", authedUser.ID)
		if err != nil {
			return nil, err
		}

		return &model.Hospital{
			Participants: &model.Participants{
				Roots:   convertParticipants(hospital.Participants.Roots),
				Staff:   convertParticipants(hospital.Participants.Staff),
				Doctors: convertParticipants(hospital.Participants.Doctors),
			},
		}, nil
	case user.RoleDoctor:
		hospital, err := r.HospitalService.CheckParticpant(ctx, "DOCTOR", authedUser.ID)
		if err != nil {
			return nil, err
		}

		return &model.Hospital{
			Participants: &model.Participants{
				Roots:   convertParticipants(hospital.Participants.Roots),
				Staff:   convertParticipants(hospital.Participants.Staff),
				Doctors: convertParticipants(hospital.Participants.Doctors),
			},
		}, nil
	case user.RoleStaff:
		hospital, err := r.HospitalService.CheckParticpant(ctx, "STAFF", authedUser.ID)
		if err != nil {
			return nil, err
		}

		return &model.Hospital{
			Participants: &model.Participants{
				Roots:   convertParticipants(hospital.Participants.Roots),
				Staff:   convertParticipants(hospital.Participants.Staff),
				Doctors: convertParticipants(hospital.Participants.Doctors),
			},
		}, nil
	default:
		return nil, errors.New("invalid role")

	}
}

// MailBoxSubscription is the resolver for the MailBoxSubscription field.
func (r *subscriptionResolver) MailBoxSubscription(ctx context.Context) (<-chan *model.MailBoxSubscriptionResponse, error) {
	fmt.Println("called")
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}

	return r.ObserverService.SubscribeToMail(ctx, authedUser.Email).(<-chan *model.MailBoxSubscriptionResponse), nil
}

// MessageBoxSubscription is the resolver for the MessageBoxSubscription field.
func (r *subscriptionResolver) MessageBoxSubscription(ctx context.Context) (<-chan *model.MessageSubscriptionResponse, error) {
	fmt.Println("called")
	authedUser, ok := ctx.Value(UserContextKey).(AuthenticatedUser)
	if !ok {
		return nil, fmt.Errorf("unauthorized: user not found in context")
	}
	return r.ObserverService.SubscribeToMessage(ctx, authedUser.Email).(<-chan *model.MessageSubscriptionResponse), nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
