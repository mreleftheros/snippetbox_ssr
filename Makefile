all:
	go run ./cmd/web
cert:
	go run ~/Documents/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
db:
	export DATABASE_URL=postgresql://postgres@localhost/main
