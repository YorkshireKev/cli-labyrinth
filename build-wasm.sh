#Builds the web version and places the wasm file into the docs folder so
#that it can be deployed as github pages.
GOOS=js GOARCH=wasm go build -o docs/cli-labyrinth.wasm
