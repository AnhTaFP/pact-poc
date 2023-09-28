package consumer

import (
	"fmt"
	"testing"

	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
	mockProvider, err := consumer.NewV4Pact(consumer.MockHTTPProviderConfig{
		Consumer: "pact-poc-consumer",
		Provider: "pact-poc-provider",
		Port:     8089,
		PactDir:  "../pacts",
	})

	assert.NoError(t, err)

	err = mockProvider.
		AddInteraction().
		Given("discount #1 exists").
		UponReceiving("a request to get discount #1").
		WithRequest("GET", "/discounts/1").
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.JSONBody(matchers.Map{
				"id":          matchers.Integer(1),
				"title":       matchers.Like("5.8% off"),
				"description": matchers.Like("5.8% off for Singaporean 58th national day"),
				"type":        matchers.Like("percentage"),
				"value":       matchers.Decimal(5.8),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			d, err := c.GetDiscount(1)

			assert.NoError(t, err)
			assert.Equal(t, 1, d.ID)

			return nil
		})

	err = mockProvider.
		AddInteraction().
		Given("discount #2 does not exist").
		UponReceiving("a request to get discount #2").
		WithRequest("GET", "/discounts/2").
		WillRespondWith(404).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			_, err := c.GetDiscount(2)

			assert.ErrorIs(t, err, errNotFound)

			return nil
		})

	err = mockProvider.
		AddInteraction().
		Given("two percentage discounts exist").
		UponReceiving("a request to get all percentage discounts").
		WithRequest("GET", "/discounts", func(b *consumer.V4RequestBuilder) {
			b.Query("type", matchers.S("percentage"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.JSONBody(matchers.Map{
				"discounts": matchers.ArrayMinMaxLike(matchers.Map{
					"id":          matchers.Integer(1),
					"title":       matchers.Like("5.8% off"),
					"description": matchers.Like("5.8% off for Singaporean 58th national day"),
					"type":        matchers.Like("percentage"),
					"value":       matchers.Decimal(5.8),
				}, 2, 2),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			discounts, err := c.GetDiscounts(map[string]string{
				"type": "percentage",
			})

			assert.NoError(t, err)
			assert.Equal(t, 2, len(discounts))

			return nil
		})

	err = mockProvider.
		AddInteraction().
		Given("no percentage discounts exist").
		UponReceiving("a request to get all percentage discounts").
		WithRequest("GET", "/discounts", func(b *consumer.V4RequestBuilder) {
			b.Query("type", matchers.S("percentage"))
		}).
		WillRespondWith(200, func(b *consumer.V4ResponseBuilder) {
			b.JSONBody(matchers.Map{
				"discounts": matchers.ArrayContaining([]interface{}{}),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			discounts, err := c.GetDiscounts(map[string]string{
				"type": "percentage",
			})

			assert.NoError(t, err)
			assert.Equal(t, 0, len(discounts))

			return nil
		})

	err = mockProvider.
		AddInteraction().
		Given("discount #1 exists").
		UponReceiving("a request to modify discount #1").
		WithRequest("PUT", "/discounts/1", func(b *consumer.V4RequestBuilder) {
			b.JSONBody(matchers.Map{
				"title":       matchers.Like("5.8% off"),
				"description": matchers.Like("5.8% off for Singaporean 58th national day"),
				"type":        matchers.Like("percentage"),
				"value":       matchers.Decimal(5.8),
			})
			b.Header("Content-Type", matchers.S("application/json"))
		}).
		WillRespondWith(200).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			err := c.PutDiscount(Discount{
				ID:          1,
				Title:       "new title",
				Description: "new description",
				Type:        "amount",
				Value:       6.5,
			})

			assert.NoError(t, err)

			return nil
		})

	err = mockProvider.
		AddInteraction().
		Given("discount #1 does not exist").
		UponReceiving("a request to modify discount #1").
		WithRequest("PUT", "/discounts/1", func(b *consumer.V4RequestBuilder) {
			b.JSONBody(matchers.Map{
				"title":       matchers.Like("5.8% off"),
				"description": matchers.Like("5.8% off for Singaporean 58th national day"),
				"type":        matchers.Like("percentage"),
				"value":       matchers.Decimal(5.8),
			})
			b.Header("Content-Type", matchers.S("application/json"))
		}).
		WillRespondWith(404).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			c := NewClient(fmt.Sprintf("http://%s:%d", config.Host, config.Port))
			err := c.PutDiscount(Discount{
				ID:          1,
				Title:       "new title",
				Description: "new description",
				Type:        "amount",
				Value:       6.5,
			})

			assert.ErrorIs(t, err, errNotFound)

			return nil
		})

	assert.NoError(t, err)
}
