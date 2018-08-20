main:
	go build ./internal
	go build -o ./cmd/wowexchange/wowexchange ./cmd/wowexchange
