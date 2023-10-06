package main

import (
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvider(t *testing.T) {
	const testServerAddr = "localhost:8081"

	db, err := initDb("discounts.test.db")
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
			"discount #1 exists": func(setup bool, state models.ProviderState) (models.ProviderStateResponse, error) {
				r, err := db.Exec("INSERT INTO discounts VALUES(NULL, ?,?,?,?)", "title", "description", "amount", 5.5)
				if err != nil {
					return models.ProviderStateResponse{}, err
				}

				lastInsertID, err := r.LastInsertId()
				if err != nil {
					return models.ProviderStateResponse{}, err
				}

				return models.ProviderStateResponse{
					"id": int(lastInsertID),
				}, nil
			},
		},
	})

	assert.NoError(t, err)
}
