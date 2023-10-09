package main

import (
	"os"
	"testing"
	"time"

	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	const testServerAddr = "localhost:8081"

	db, err := initDb("../discounts.test.db")
	if err != nil {
		t.Fatal("cannot init db: ", err.Error())
	}

	go startServer(testServerAddr, db)

	_ = os.Setenv("PACT_DO_NOT_TRACK", "true")
	v := provider.NewVerifier()

	err = v.VerifyProvider(t, provider.VerifyRequest{
		ProviderBaseURL: "http://" + testServerAddr,
		Provider:        "pact-poc-provider",
		ProviderVersion: "pact-proc-provider-v1.0",
		PactDirs:        []string{"../pacts"}, // if we're using a Pact broker such as pact flow, this is not needed
		StateHandlers: models.StateHandlers{
			"a discount exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup { // setup hook
					_, err := db.Exec("INSERT INTO discounts VALUES(NULL, 'title', 'description', 'amount', 5.5, ?, NULL)", time.Now())
					if err != nil {
						return nil, err
					}
				} else { // teardown hook
					_, err := db.Exec("DELETE FROM discounts")
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
			"no discount exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				_, err := db.Exec("DELETE FROM discounts")
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
			"two discounts of the same type exist": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup { // setup hook
					for i := 0; i < 2; i++ {
						_, err := db.Exec("INSERT INTO discounts VALUES(NULL, 'title', 'description', 'percentage', 5.5, ?, NULL)", time.Now())
						if err != nil {
							return nil, err
						}
					}
				} else { // teardown hook
					_, err := db.Exec("DELETE FROM discounts WHERE type = 'percentage'")
					if err != nil {
						return nil, err
					}
				}

				return models.ProviderStateResponse{
					"discount_type": "percentage",
				}, nil
			},
			"discounts limit is reached": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				if setup {
					for i := 0; i < 3; i++ {
						_, err := db.Exec("INSERT INTO discounts VALUES(NULL, 'title', 'description', 'percentage', 5.5, ?, NULL)", time.Now())
						if err != nil {
							return nil, err
						}
					}
				} else {
					_, err := db.Exec("DELETE FROM discounts")
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	})

	assert.NoError(t, err)
}
