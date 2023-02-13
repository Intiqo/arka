package payment

import (
	"strconv"
	"time"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/product"
	"github.com/stripe/stripe-go/v72/sub"

	"github.com/adwitiyaio/arka/secrets"
)

const stripeTokenKey = "STRIPE_TOKEN"

type stripePaymentManager struct {
	sm secrets.Manager
}

func (s stripePaymentManager) initialize() {
	stripe.Key = s.sm.GetValueForKey(stripeTokenKey)
}

// CreateCustomer ... Creates a new customer on the payment gateway and returns the ID of the customer
func (s stripePaymentManager) CreateCustomer(params CustomerParams) (string, error) {
	stripeParams := &stripe.CustomerParams{
		Name:  stripe.String(params.Name),
		Phone: stripe.String(params.Phone),
		Email: stripe.String(params.Email),
	}

	result, err := customer.New(stripeParams)

	if err != nil {
		return "", err
	}

	return result.ID, err
}

// FindDefaultPaymentMethodForCustomer ... Finds the default payment method for a customer
func (s stripePaymentManager) FindDefaultPaymentMethodForCustomer(customerId string) (string, error) {
	result, err := customer.Get(customerId, nil)
	if err != nil {
		return "", err
	}
	if result.InvoiceSettings.DefaultPaymentMethod == nil {
		return "", err
	}
	return result.InvoiceSettings.DefaultPaymentMethod.ID, err
}

// SetDefaultPaymentMethodForCustomer ... Sets the payment method as the default for a customer
func (s stripePaymentManager) SetDefaultPaymentMethodForCustomer(customerId string, paymentId string) error {
	stripeParams := &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentId),
		},
	}
	_, err := customer.Update(customerId, stripeParams)
	if err != nil {
		return err
	}
	return nil
}

// FindAllPaymentMethodsForCustomer ... Finds all the payment methods associated with a customer
func (s stripePaymentManager) FindAllPaymentMethodsForCustomer(customerId string) ([]CardParams, error) {
	stripeParams := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerId),
		Type:     stripe.String("card"),
	}

	result := make([]CardParams, 0)

	i := paymentmethod.List(stripeParams)
	for i.Next() {
		pm := i.PaymentMethod()
		card := CardParams{
			CardHolder:  pm.Metadata["cardholder"],
			Number:      pm.Card.Last4,
			ExpiryMonth: strconv.Itoa(int(pm.Card.ExpMonth)),
			ExpiryYear:  strconv.Itoa(int(pm.Card.ExpYear)),
			Brand:       string(pm.Card.Brand),
			PaymentId:   pm.ID,
		}
		result = append(result, card)
	}

	return result, nil
}

// CreatePaymentMethod ... Creates a new payment method for a customer
func (s stripePaymentManager) CreatePaymentMethod(customerId string, params *CardParams) error {
	stripeParams := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(params.Number),
			ExpMonth: stripe.String(params.ExpiryMonth),
			ExpYear:  stripe.String(params.ExpiryYear),
			CVC:      stripe.String(params.CVC),
		},
		Type: stripe.String("card"),
	}
	stripeParams.AddMetadata("cardholder", params.CardHolder)
	stripeParams.AddMetadata("creator", params.Creator)

	pm, err := paymentmethod.New(stripeParams)

	if err != nil {
		return err
	}

	pmAttachParams := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerId),
	}

	_, err = paymentmethod.Attach(pm.ID, pmAttachParams)
	if err != nil {
		return err
	}

	params.Brand = string(pm.Card.Brand)
	params.PaymentId = pm.ID

	return err
}

// UpdatePaymentMethod ... Update a payment method
func (s stripePaymentManager) UpdatePaymentMethod(paymentId string, params *CardParams) error {
	stripeParams := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			ExpMonth: stripe.String(params.ExpiryMonth),
			ExpYear:  stripe.String(params.ExpiryYear),
		},
	}
	stripeParams.AddMetadata("cardholder", params.CardHolder)
	stripeParams.AddMetadata("updater", params.Creator)

	pm, err := paymentmethod.Update(paymentId, stripeParams)
	if err != nil {
		return err
	}

	params.Number = pm.Card.Last4
	params.Brand = string(pm.Card.Brand)
	params.PaymentId = pm.ID

	return err
}

// DeletePaymentMethod ... Deletes a payment method
func (s stripePaymentManager) DeletePaymentMethod(paymentId string) error {
	_, err := paymentmethod.Detach(paymentId, nil)
	return err
}

// CreateProduct ... Creates a new product on the payment gateway
func (s stripePaymentManager) CreateProduct(name string) (string, error) {
	stripeParams := &stripe.ProductParams{
		Name: stripe.String(name),
	}

	p, err := product.New(stripeParams)
	if err != nil {
		return "", err
	}
	return p.ID, err
}

// UpdateProduct ... Updates an existing product on the payment gateway
func (s stripePaymentManager) UpdateProduct(id string, name string) error {
	stripeParams := &stripe.ProductParams{Name: stripe.String(name)}
	_, err := product.Update(id, stripeParams)
	return err
}

// DeleteProduct ... Deletes an existing product
func (s stripePaymentManager) DeleteProduct(id string) error {
	_, err := product.Del(id, nil)
	return err
}

// CreatePrice ... Creates a new price for an existing product
func (s stripePaymentManager) CreatePrice(params PriceParams) (string, error) {
	stripeParams := &stripe.PriceParams{
		Currency: stripe.String(params.Currency),
		Nickname: stripe.String(params.Name),
		Product:  stripe.String(params.ProductId),
		Recurring: &stripe.PriceRecurringParams{
			Interval:        stripe.String(params.IntervalType),
			IntervalCount:   stripe.Int64(params.Interval),
			TrialPeriodDays: stripe.Int64(params.TrialPeriod),
		},
		UnitAmount: stripe.Int64(int64(params.Amount * 100)),
	}

	result, err := price.New(stripeParams)
	if err != nil {
		return "", err
	}
	return result.ID, err
}

// UpdatePrice ... Updates an existing price for an existing product
func (s stripePaymentManager) UpdatePrice(priceId string, params PriceParams) error {
	stripeParams := &stripe.PriceParams{
		Currency: stripe.String(params.Currency),
		Nickname: stripe.String(params.Name),
		Product:  stripe.String(params.ProductId),
		Recurring: &stripe.PriceRecurringParams{
			Interval:        stripe.String(params.IntervalType),
			IntervalCount:   stripe.Int64(params.Interval),
			TrialPeriodDays: stripe.Int64(params.TrialPeriod),
		},
		TaxBehavior:       stripe.String(string(stripe.PriceTaxBehaviorExclusive)),
		UnitAmountDecimal: stripe.Float64(params.Amount * 100),
	}

	_, err := price.Update(priceId, stripeParams)
	if err != nil {
		return err
	}
	return err
}

// CreateSubscription ... Creates a new subscription for a customer and a plan
func (s stripePaymentManager) CreateSubscription(params SubscriptionParams) (SubscriptionResponse, error) {
	var response SubscriptionResponse

	stripeParams := &stripe.SubscriptionParams{
		Customer: stripe.String(params.CustomerId),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(params.PriceId),
			},
		},
	}

	result, err := sub.New(stripeParams)
	if err != nil {
		return response, err
	}

	response = SubscriptionResponse{
		Id:                 result.ID,
		CurrentPeriodStart: time.Unix(result.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(result.CurrentPeriodEnd, 0),
		Status:             string(result.Status),
	}
	return response, nil
}

// CancelSubscription ... Cancels an existing subscription for a customer and a plan
func (s stripePaymentManager) CancelSubscription(id string) (SubscriptionResponse, error) {
	var response SubscriptionResponse
	result, err := sub.Cancel(id, nil)
	if err != nil {
		return response, err
	}
	response = SubscriptionResponse{
		Id:                 result.ID,
		CurrentPeriodStart: time.Unix(result.CurrentPeriodStart, 0),
		CurrentPeriodEnd:   time.Unix(result.CurrentPeriodEnd, 0),
		Status:             string(result.Status),
	}
	return response, err
}

func (s stripePaymentManager) CreatePaymentIntent(params IntentParams) error {
	stripeParams := &stripe.PaymentIntentParams{
		Amount:        stripe.Int64(int64(params.Amount * 100)),
		Currency:      stripe.String(params.Currency),
		Customer:      stripe.String(params.CustomerId),
		Description:   stripe.String(params.Description),
		Confirm:       stripe.Bool(true),
		PaymentMethod: stripe.String(params.PaymentMethodId),
	}

	_, err := paymentintent.New(stripeParams)
	return err
}
