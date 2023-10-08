## Pact PoC
Proof of concept for service and service contract testing.

### Pre-requisites
- Install pact-go version 2 by this command `go install github.com/pact-foundation/pact-go/v2@2.x.x`
- Download and install all the required libraries by this command `pact-go -l DEBUG install`

### Folder Structure
- `provider` is an HTTP server that provides a CRUD API, storing data in sqlite.
- `consumer` provides an HTTP client to call the CRUD API.
- `pact` is where the contract between the Consumer and the Provider is stored.

### Contract Test
- To run consumer test, navigate to the root directory of the repository, then run `go test -v ./consumer/...`. This will produce a pact file under `/pacts`
- To run provider test, navigate to the root directory of the repository, then run `go test -v ./provider/providertest/...`