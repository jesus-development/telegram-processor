build-demo:
	@cp -n .env.example .env
	@go mod download
	@go build -tags demo -o telegram-processor