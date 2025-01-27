package observers

import (
	"context"
	"sync"
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

func (o *ObserverService) Subscribe(ctx context.Context, resourceId string, subscriptionType string) interface{} {
	o.mu.Lock()
	defer o.mu.Unlock()

	if _, exists := o.observers[resourceId]; !exists {
		o.observers[resourceId] = &Subscriptions{}
	}

	var ch interface{}
	switch subscriptionType {
	case "MAIL":
		ch = make(chan *MailBoxSubscriptionResopnse, 10)
		break
	case "MESSAGE":
		ch = make(chan *MessageSubscriptionResponse, 10)

	}

	o.observers[resourceId].mu.Lock()
	o.observers[resourceId].Channels[subscriptionType] = ch
	o.observers[resourceId].mu.Unlock()

	go func() {
		<-ctx.Done()
		o.UnSubscribe(resourceId, subscriptionType)
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
		close(observer.Channels[subscriptionType].(chan *MailBoxSubscriptionResopnse))
		delete(observer.Channels, subscriptionType)
	case "MESSAGE":
		close(observer.Channels[subscriptionType].(chan *MessageSubscriptionResponse))
		delete(observer.Channels, subscriptionType)
	}

	if len(observer.Channels) == 0 {
		o.mu.Lock()
		delete(o.observers, resourceId)
		o.mu.Unlock()
	}
}

func (o *ObserverService) Publish(message interface{}, resourceId string, subscriptionType string) {
	o.mu.RLock()
	observer, ok := o.observers[resourceId]
	o.mu.RUnlock()

	if !ok {
		return
	}

	observer.mu.Lock()
	defer observer.mu.Unlock()

	switch subscriptionType {
	case "MAIL":
		observer.Channels[subscriptionType].(chan *MailBoxSubscriptionResopnse) <- message.(*MailBoxSubscriptionResopnse)
	case "MESSAGE":
		observer.Channels[subscriptionType].(chan *MessageSubscriptionResponse) <- message.(*MessageSubscriptionResponse)
	}

}

func (o *ObserverService) SubscribeToMail(ctx context.Context, mailId string) interface{} {
	return o.Subscribe(ctx, mailId, "MAIL")
}

func (o *ObserverService) PublishMail(receiver string, message *MailBoxSubscriptionResopnse) {
	o.Publish(message, receiver, "MAIL")
}

func (o *ObserverService) SubscribeToMessage(ctx context.Context, mailId string) interface{} {
	return o.Subscribe(ctx, mailId, "MESSAGE")
}

func (o *ObserverService) PublishMessage(receiver string, message *MessageSubscriptionResponse) {
	o.Publish(message, receiver, "MESSAGE")
}
