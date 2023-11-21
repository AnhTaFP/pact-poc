## Pact PoC
Proof of concept for service and service contract testing.

In this PoC, we explore how contract testing with Pact can be done for a basic CRUD app. `provider` is the application that provides a set of CRUD API end-points to manage discounts, leveraging sqlite as the data store. `consumer` is the application that invokes API call to create/read/update/delete discounts. In a real world scenario, these 2 applications may belong to different code bases, but for the sake of simplicity, we have both application source codes in the same code base.

Because Pact is consumer-driven, we can start exploring [consumer_test.go](https://github.com/AnhTaFP/pact-poc/blob/master/consumer/consumer_test.go), which asserts the consumer's expectation of the API and produces the expectation result in a [special formatted file](https://github.com/AnhTaFP/pact-poc/blob/master/pacts/pact-poc-consumer-pact-poc-provider.json) that can be understood by Pact verifier. `provider` is where the verification process works in which all assertions are applied against the API to ensure it meets what's expected by the `consumer`, check out [provider_test.go](https://github.com/AnhTaFP/pact-poc/blob/master/provider/provider_test.go) for more info

### Pre-requisites
- Install pact-go version 2 by this command `go install github.com/pact-foundation/pact-go/v2@2.x.x`
- Download and install all the required libraries by this command `pact-go -l DEBUG install`
- Install [sqlite](https://formulae.brew.sh/formula/sqlite)

### Folder Structure
- `provider` is an HTTP server that provides a CRUD API, storing data in sqlite.
- `consumer` provides an HTTP client to call the CRUD API.
- `pact` is where the contract between the Consumer and the Provider is stored.

### Contract Test
- To run consumer test, navigate to the root directory of the repository, then run `go test -v ./consumer/...`. This will produce a pact file under `/pacts`
- To run provider test, navigate to the root directory of the repository, then run `go test -v ./provider/providertest/...`
