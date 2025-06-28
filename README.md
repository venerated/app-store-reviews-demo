# app-store-reviews-demo

Demo app that fetches and caches reviews from the App Store, using a Go backend and a React frontend.

## Features

- Fetches latest App Store reviews for a given app
- Caches results locally to avoid redundant API calls
- Gracefully handles missing or malformed cache data

## Tech Stack

- Backend: Go, Gin
- Frontend: React, TypeScript, Vite
- Build Tools: Make, Yarn

## Requirements

- Node v23.x
- yarn v4.9.2

## Installation

- Copy `app/.env.example` to `app/.env` or `app/.env.local`

### Using Makefile

- Run `make` from root folder
- This also runs the backend and frontend, so you can now open the app at http://localhost:5173/

### Manually

#### Backend

- From the root folder, run `cd api && go mod tidy`

#### Frontend

- From the root folder, run `cd app && yarn install`

## Running the Site

Using `make` from the root folder to install and run is recommended, but if you want to run manually:

1. From the root folder, run `cd api && go run .`
2. From the root folder, run `cd app && yarn dev`

View the site at http://localhost:5173/

## Tests

### Backend

1. From thr oot folder, run `cd api && go test -v ./...`

### Frontend

1. From the root folder, run `cd app && yarn test`
