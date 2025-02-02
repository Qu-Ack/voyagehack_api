package observers

import (
	"context"
	"fmt"
	"sync"

	"github.com/Qu-Ack/voyagehack_api/api/graph/model"
)

// chan *MailBoxSubscriptionResponse
type Subscriptions struct {
	mu       sync.RWMutex
	Channels map[string]interface{}
}

type ObserverService struct {
	mu        sync.RWMutex
	observers map[string]*Subscriptions
}

func NewObserversService() *ObserverService {
	return &ObserverService{
		observers: make(map[string]*Subscriptions),
	}
}

func (o *ObserverService) subscribetomail(ctx context.Context, resourceId string) <-chan *model.MailBoxSubscriptionResponse {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.observers[resourceId]; !exists {
		o.observers[resourceId] = &Subscriptions{
			Channels: make(map[string]interface{}),
		}
	}

	var ch chan *model.MailBoxSubscriptionResponse
	ch = make(chan *model.MailBoxSubscriptionResponse, 10)

	o.observers[resourceId].mu.Lock()
	o.observers[resourceId].Channels["MAIL"] = ch
	o.observers[resourceId].mu.Unlock()

	fmt.Println("We have subscribed To Mail", o.observers[resourceId])
	go func() {
		<-ctx.Done()
		o.UnSubscribe(resourceId, "MAIL")
	}()

	return ch

}

func (o *ObserverService) subscribetomessage(ctx context.Context, resourceId string) <-chan *model.MessageSubscriptionResponse {
	fmt.Println("in subscribe to message")
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.observers[resourceId]; !exists {
		o.observers[resourceId] = &Subscriptions{
			Channels: make(map[string]interface{}),
		}
	}

	var ch chan *model.MessageSubscriptionResponse
	ch = make(chan *model.MessageSubscriptionResponse, 10)

	o.observers[resourceId].mu.Lock()
	o.observers[resourceId].Channels["MESSAGE"] = ch
	o.observers[resourceId].mu.Unlock()
	fmt.Println(o.observers[resourceId])

	go func() {
		<-ctx.Done()
		o.UnSubscribe(resourceId, "MESSAGE")
	}()

	return ch

}

func (o *ObserverService) UnSubscribe(resourceId string, subscriptionType string) {
	o.mu.RLock()
	observer, exists := o.observers[resourceId]
	o.mu.RUnlock()
	if !exists {
		return
	}

	observer.mu.Lock()
	defer observer.mu.Unlock()

	switch subscriptionType {
	case "MAIL":
		close(observer.Channels[subscriptionType].(chan *model.MailBoxSubscriptionResponse))
		delete(observer.Channels, subscriptionType)
	case "MESSAGE":
		close(observer.Channels[subscriptionType].(chan *model.MessageSubscriptionResponse))
		delete(observer.Channels, subscriptionType)
	}

	if len(observer.Channels) == 0 {
		o.mu.Lock()
		delete(o.observers, resourceId)
		o.mu.Unlock()
	}
}

func (o *ObserverService) PublishToMail(message *model.MailBoxSubscriptionResponse, resourceId string) {
	o.mu.RLock()
	observer, ok := o.observers[resourceId]
	o.mu.RUnlock()

	fmt.Println("observer is ", observer)
	if !ok {
		return
	}

	observer.mu.Lock()
	chinterface, exists := observer.Channels["MAIL"]
	observer.mu.Unlock()

	if !exists {
		fmt.Println("NO MAIL CHANNEL")
		return
	}

	observer.mu.Lock()
	defer observer.mu.Unlock()

	chinterface.(chan *model.MailBoxSubscriptionResponse) <- message

}

func (o *ObserverService) PublishToMessage(message *model.MessageSubscriptionResponse, resourceId string) {
	o.mu.RLock()
	observer, ok := o.observers[resourceId]
	o.mu.RUnlock()
	fmt.Println(observer)

	if !ok {
		fmt.Println("returning in the ok")
		return
	}

	observer.mu.Lock()
	chinterface, exists := observer.Channels["MESSAGE"]
	observer.mu.Unlock()

	if !exists {
		fmt.Println("NO MESSAGE CHANNEL")
		return
	}

	observer.mu.Lock()
	defer observer.mu.Unlock()

	chinterface.(chan *model.MessageSubscriptionResponse) <- message

}

func (o *ObserverService) SubscribeToMail(ctx context.Context, mailId string) <-chan *model.MailBoxSubscriptionResponse {
	return o.subscribetomail(ctx, mailId)
}

func (o *ObserverService) PublishMail(receiver string, message *model.MailBoxSubscriptionResponse) {
	fmt.Println("in publish mail time")
	o.PublishToMail(message, receiver)
}

func (o *ObserverService) SubscribeToMessage(ctx context.Context, mailId string) <-chan *model.MessageSubscriptionResponse {
	return o.subscribetomessage(ctx, mailId)
}

func (o *ObserverService) PublishMessage(receiver string, message *model.MessageSubscriptionResponse) {
	fmt.Println(receiver)
	o.PublishToMessage(message, receiver)
}
