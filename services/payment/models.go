package payment

type OrderRequest struct {
	Amount int `json:"amount"`
}
type Order struct {
	ID         string            `json:"id"`
	Entity     string            `json:"entity"`
	Amount     int               `json:"amount"`
	AmountPaid int               `json:"amount_paid"`
	AmountDue  int               `json:"amount_due"`
	Currency   string            `json:"currency"`
	Receipt    string            `json:"receipt"`
	Status     string            `json:"status"`
	Attempts   int               `json:"attempts"`
	Notes      map[string]string `json:"notes"`
	CreatedAt  int64             `json:"created_at"`
}

type ValidatePaymentRequest struct {
	RazorpayPaymentId string `json:"razorpay_payment_id"`
	RazorpayOrderId   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
}
