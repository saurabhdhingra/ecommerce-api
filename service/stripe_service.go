package service

import (
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/client"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

// StripeService defines the contract for payment operations.
type StripeService interface {
	CreatePaymentIntent(amountCents int64, currency string, description string) (*stripe.PaymentIntent, error)
	// Other methods: HandleWebhook, CapturePayment, etc.
}

// stripeService is the concrete implementation using stripe-go.
type stripeService struct {
	sc *client.API
}

// NewStripeService initializes and returns the Stripe service.
func NewStripeService(key string) StripeService {
	sc := &client.API{}
	sc.Init(key, nil)
	return &stripeService{sc: sc}
}

// CreatePaymentIntent creates a new Payment Intent with Stripe.
func (s *stripeService) CreatePaymentIntent(amountCents int64, currency string, description string) (*stripe.PaymentIntent, error) {
	// Stripe requires amount in smallest currency unit (e.g., cents for USD)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Description: stripe.String(description),
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, err
	}
	
	return pi, nil
}