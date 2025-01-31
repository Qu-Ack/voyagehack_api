package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Qu-Ack/voyagehack_api/services/user"
	razorpay "github.com/razorpay/razorpay-go"
	"github.com/razorpay/razorpay-go/utils"
	"golang.org/x/exp/rand"
)

type PaymentService struct {
	ongoingPayments map[string]string
	client          *razorpay.Client
}

func NewPaymentService() *PaymentService {
	razorpay_key := os.Getenv("RAZORPAY_KEY")
	razorpay_secret := os.Getenv("RAZORPAY_SECRET")
	client := razorpay.NewClient(razorpay_key, razorpay_secret)
	return &PaymentService{
		ongoingPayments: make(map[string]string, 0),
		client:          client,
	}
}

func (p *PaymentService) NewOrder(orderRequest *OrderRequest, requester user.PublicUser) (string, error) {
	orderId := generateReceiptID()
	data := map[string]interface{}{
		"amount":   orderRequest.Amount,
		"currency": "INR",
		"receipt":  orderId,
	}
	body, err := p.client.Order.Create(data, nil)

	if err != nil {
		return "", err
	}

	id, ok := body["id"].(string)

	if !ok {
		return "", errors.New("order was not created successfully")
	}

	p.ongoingPayments[requester.Email] = id

	return id, nil
}

func (p *PaymentService) ValidatePayment(validatePayment *ValidatePaymentRequest, requester user.PublicUser) error {
	orderId := p.ongoingPayments[requester.Email]
	razorpaySecret := os.Getenv("RAZORPAY_SECRET")

	fmt.Println(validatePayment.RazorpayPaymentId)
	params := map[string]interface{}{
		"razorpay_order_id":   orderId,
		"razorpay_payment_id": validatePayment.RazorpayPaymentId,
	}
	if utils.VerifyPaymentSignature(params, validatePayment.RazorpaySignature, razorpaySecret) {
		return nil
	}

	return errors.New("signatures don't match")
}

func GenerateSignature(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func generateRandomAlphanumeric(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}

func generateReceiptID() string {
	randomPart := generateRandomAlphanumeric(10)
	return "receipt_" + randomPart
}
