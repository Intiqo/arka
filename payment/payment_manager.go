package payment

import (
	"time"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
)

const DependencyPaymentManager = "payment_manager"

const ProviderStripe = "STRIPE"

type CustomerParams struct {
	Email string
	Name  string
	Phone string
}

type CardParams struct {
	CardHolder  string
	Number      string
	ExpiryMonth string
	ExpiryYear  string
	CVC         string
	Brand       string
	PaymentId   string
	Creator     string
}

type PriceParams struct {
	Currency     string
	ProductId    string
	Amount       float64
	Name         string
	IntervalType string
	Interval     int64
	TrialPeriod  int64
}

type SubscriptionParams struct {
	CustomerId string
	PriceId    string
}

type SubscriptionResponse struct {
	Id                 string
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	Status             string
}

type IntentParams struct {
	CustomerId      string
	Currency        string
	Amount          float64
	Description     string
	PaymentMethodId string
}

type Manager interface {
	// All methods related to customer

	// CreateCustomer ... Creates a new customer on the payment gateway and returns the ID of the customer
	CreateCustomer(params CustomerParams) (string, error)
	// FindDefaultPaymentMethodForCustomer ... Finds the default payment method for a customer
	FindDefaultPaymentMethodForCustomer(customerId string) (string, error)
	// SetDefaultPaymentMethodForCustomer ... Sets the payment method as the default for a customer
	SetDefaultPaymentMethodForCustomer(customerId string, paymentId string) error

	// All methods related to payment method

	// FindAllPaymentMethodsForCustomer ... Finds all the payment methods associated with a customer
	FindAllPaymentMethodsForCustomer(customerId string) ([]CardParams, error)
	// CreatePaymentMethod ... Creates a new payment method for a customer
	CreatePaymentMethod(customerId string, c *CardParams) error
	// UpdatePaymentMethod ... Update a payment method
	UpdatePaymentMethod(paymentId string, c *CardParams) error
	// DeletePaymentMethod ... Deletes a payment method
	DeletePaymentMethod(paymentId string) error

	// All methods related to product

	// CreateProduct ... Creates a new product
	CreateProduct(name string) (string, error)
	// UpdateProduct ... Updates an existing product
	UpdateProduct(id string, name string) error
	// DeleteProduct ... Deletes an existing product
	DeleteProduct(id string) error

	// All methods related to price

	// CreatePrice ... Creates a new price for an existing product
	CreatePrice(params PriceParams) (string, error)
	// UpdatePrice ... Updates an existing price for an existing product
	UpdatePrice(priceId string, params PriceParams) error

	// All methods related to subscription

	// CreateSubscription ... Creates a new subscription for a customer and a plan
	CreateSubscription(params SubscriptionParams) (SubscriptionResponse, error)
	// CancelSubscription ... Cancels an existing subscription for a customer and a plan
	CancelSubscription(id string) (SubscriptionResponse, error)

	// All methods related to independent charges

	// CreatePaymentIntent ... Creates a new payment intent against a customer
	CreatePaymentIntent(params IntentParams) error
}

func Bootstrap(provider string) {
	dm := dependency.GetManager()
	var pm interface{}
	switch provider {
	case ProviderStripe:
		pm = &stripePaymentManager{
			cm: dm.Get(config.DependencyConfigManager).(config.Manager),
		}
		pm.(*stripePaymentManager).initialize()
	}
	dm.Register(DependencyPaymentManager, pm)
}
