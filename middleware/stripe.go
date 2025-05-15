package middleware

import (
	"os"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
}

func CreateCheckoutSession(userEmail, plan string) (string, error) {
	var priceID string
	if plan == "monthly" {
		priceID = "price_1234_month" // Replace with your real price ID
	} else {
		priceID = "price_5678_year"
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		CustomerEmail:      stripe.String(userEmail),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String("https://yourdomain.com/success"),
		CancelURL:  stripe.String("https://yourdomain.com/cancel"),
	}
	s, err := session.New(params)
	if err != nil {
		return "", err
	}
	return s.URL, nil
}
