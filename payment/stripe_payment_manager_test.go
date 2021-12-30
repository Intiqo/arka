package payment

import (
	"os"
	"testing"

	"github.com/adwitiyaio/arka/config"
	"github.com/adwitiyaio/arka/dependency"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/product"
)

type StripePaymentGatewayTestSuite struct {
	suite.Suite

	pg Manager

	customerIds []string
	paymentIds  []string
	productIds  []string
}

func TestStripePaymentGateway(t *testing.T) {
	suite.Run(t, new(StripePaymentGatewayTestSuite))
}

func (ts *StripePaymentGatewayTestSuite) SetupSuite() {
	config.Bootstrap(config.ProviderEnvironment, "../test.env")
	err := os.Setenv("CI", "true")
	require.NoError(ts.T(), err)
	Bootstrap(ProviderStripe)
	ts.pg = dependency.GetManager().Get(DependencyPaymentManager).(Manager)

	ts.customerIds = make([]string, 0)
	ts.paymentIds = make([]string, 0)

	ts.createCustomer()
}

func (ts *StripePaymentGatewayTestSuite) TearDownSuite() {
	for _, id := range ts.customerIds {
		_, err := customer.Del(id, nil)
		require.NoError(ts.T(), err)
	}
	for _, id := range ts.productIds {
		_, err := product.Del(id, nil)
		require.NoError(ts.T(), err)
	}
}

func (ts *StripePaymentGatewayTestSuite) createCustomer() {
	c := CustomerParams{
		Name:  gofakeit.Name(),
		Email: gofakeit.Email(),
		Phone: gofakeit.Phone(),
	}
	customerId, err := ts.pg.CreateCustomer(c)
	require.NoError(ts.T(), err)
	ts.customerIds = append(ts.customerIds, customerId)

	pm := &CardParams{
		Number:      "4242424242424242",
		ExpiryMonth: "02",
		ExpiryYear:  "2025",
		CVC:         "123",
		Creator:     gofakeit.UUID(),
	}
	err = ts.pg.CreatePaymentMethod(customerId, pm)
	require.NoError(ts.T(), err)
	ts.paymentIds = append(ts.paymentIds, pm.PaymentId)

	err = ts.pg.SetDefaultPaymentMethodForCustomer(customerId, pm.PaymentId)
	require.NoError(ts.T(), err)
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_FindAllPaymentMethodsForCustomer() {
	ts.Run("unknown customer", func() {
		pms, err := ts.pg.FindAllPaymentMethodsForCustomer(gofakeit.UUID())
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), 0, len(pms))
	})

	ts.Run("success", func() {
		pms, err := ts.pg.FindAllPaymentMethodsForCustomer(ts.customerIds[0])
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), "4242", pms[0].Number)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_FindDefaultPaymentMethodForCustomer() {
	ts.Run("unknown customer", func() {
		_, err := ts.pg.FindDefaultPaymentMethodForCustomer(gofakeit.UUID())
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		paymentId, err := ts.pg.FindDefaultPaymentMethodForCustomer(ts.customerIds[0])
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), ts.paymentIds[0], paymentId)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_UpdatePaymentMethod() {
	ts.Run("unknown payment id", func() {
		pm := &CardParams{
			ExpiryMonth: "05",
			ExpiryYear:  "2026",
		}
		err := ts.pg.UpdatePaymentMethod(gofakeit.UUID(), pm)
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		pm := &CardParams{
			ExpiryMonth: "05",
			ExpiryYear:  "2026",
		}
		err := ts.pg.UpdatePaymentMethod(ts.paymentIds[0], pm)
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), "05", pm.ExpiryMonth)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_DeletePaymentMethod() {
	ts.Run("unknown payment id", func() {
		err := ts.pg.DeletePaymentMethod(gofakeit.UUID())
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		pm := &CardParams{
			Number:      "4242424242424242",
			ExpiryMonth: "03",
			ExpiryYear:  "2025",
			CVC:         "123",
			Creator:     gofakeit.UUID(),
		}
		err := ts.pg.CreatePaymentMethod(ts.customerIds[0], pm)
		require.NoError(ts.T(), err)
		ts.paymentIds = append(ts.paymentIds, pm.PaymentId)

		err = ts.pg.DeletePaymentMethod(pm.PaymentId)
		require.NoError(ts.T(), err)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_SetDefaultPaymentMethod() {
	ts.Run("unknown payment id", func() {
		err := ts.pg.SetDefaultPaymentMethodForCustomer(ts.customerIds[0], gofakeit.UUID())
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		pm := &CardParams{
			Number:      "4242424242424242",
			ExpiryMonth: "04",
			ExpiryYear:  "2025",
			CVC:         "123",
			Creator:     gofakeit.UUID(),
		}
		err := ts.pg.CreatePaymentMethod(ts.customerIds[0], pm)
		require.NoError(ts.T(), err)
		ts.paymentIds = append(ts.paymentIds, pm.PaymentId)

		err = ts.pg.SetDefaultPaymentMethodForCustomer(ts.customerIds[0], pm.PaymentId)
		require.NoError(ts.T(), err)

		defaultPaymentId, err := ts.pg.FindDefaultPaymentMethodForCustomer(ts.customerIds[0])
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), pm.PaymentId, defaultPaymentId)

		err = ts.pg.SetDefaultPaymentMethodForCustomer(ts.customerIds[0], ts.paymentIds[0])
		require.NoError(ts.T(), err)

		defaultPaymentId, err = ts.pg.FindDefaultPaymentMethodForCustomer(ts.customerIds[0])
		require.NoError(ts.T(), err)
		assert.Equal(ts.T(), ts.paymentIds[0], defaultPaymentId)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_CreateProduct() {
	ts.Run("success", func() {
		prId, err := ts.pg.CreateProduct(gofakeit.Name())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), prId)
		ts.productIds = append(ts.productIds, prId)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_UpdateProduct() {
	ts.Run("unknown product", func() {
		err := ts.pg.UpdateProduct(gofakeit.UUID(), gofakeit.Name())
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		prId, err := ts.pg.CreateProduct(gofakeit.Name())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), prId)
		ts.productIds = append(ts.productIds, prId)

		err = ts.pg.UpdateProduct(prId, gofakeit.Name())
		assert.NoError(ts.T(), err)
	})
}

func (ts *StripePaymentGatewayTestSuite) Test_stripePaymentGateway_DeleteProduct() {
	ts.Run("unknown product", func() {
		err := ts.pg.DeleteProduct(gofakeit.UUID())
		assert.Error(ts.T(), err)
	})

	ts.Run("success", func() {
		prId, err := ts.pg.CreateProduct(gofakeit.Name())
		assert.NoError(ts.T(), err)
		assert.NotNil(ts.T(), prId)

		err = ts.pg.DeleteProduct(prId)
		assert.NoError(ts.T(), err)
	})
}

// Notes: We cannot currently test price & subscription APIs since it's not possible to clear these records on Stripe
