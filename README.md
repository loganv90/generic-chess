# Generic Chess
<img src="https://github.com/loganv90/generic-chess/blob/main/art/rp_scaled.png" width="80rem"/>

## Description
Written in Go, this is a simple chess engine that can be used to play chess on any size board and with any number of players.

## Installation
1. Clone the repository
2. cd vite-app && npm install && npm run dev
3. cd go-app && go mod tidy && go run .

## Usage
1. Open a browser and navigate to localhost:3000
2. Click one of the buttons to start a game

## Running Tests
- go test ./...
- go test ./... -coverprofile cover.out
- go test -v -run \<test_name\> ./...
- go tool cover -func cover.out
- go clean -testcache

