package providertest

import (
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	_ = os.Setenv("PACT_DO_NOT_TRACK", "true")
	v := provider.NewVerifier()

	err := v.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: "http://localhost:8080",
		Provider:        "pact-poc-provider",
		ProviderVersion: "pact-proc-provider-v1.0",
		PactDirs:        []string{"../../pacts"}, // if we're using a Pact broker such as pact flow, this is not needed
		StateHandlers: models.StateHandlers{
			"discount #1 exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				return models.ProviderStateResponse{}, nil
			},
		},
	})

	assert.NoError(t, err)
}
