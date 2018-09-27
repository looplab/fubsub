# fubsub

Silly GCP pubsub-as-file example. Used as a demonstration of Go 1.11 modules in a Go West meetup in Sweden.

https://meetup.com/sweden-go-west

# Usage

Run `make test` to run the tests, it will start the GCP pubsub emulator in the background using Docker Compose. Run `make stop-emulator` to stop the emaulator.
