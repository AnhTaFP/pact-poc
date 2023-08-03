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

	assert.NoError(t, err)
}
