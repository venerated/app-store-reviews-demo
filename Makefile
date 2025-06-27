APP_DIR := app
API_DIR := api

.PHONY: run backend frontend deps

# Install dependencies and run backend and frontend together
all: deps
	make backend & \
	make frontend & \
	wait

# Get Go deps
deps:
	cd $(API_DIR) && go mod tidy

# Run backend
backend:
	cd $(API_DIR) && go run .

# Install and serve frontend
frontend:
	cd $(APP_DIR) && yarn install && yarn dev
