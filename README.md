# Fake Shop
Just a little hop demo which handles inventory and carts. Promotions can be applied.

## How To Use
The CI creates aa binary with a sample configuration at the default path and uploads
it to Github as an artifact. You can download that binary and run it on most Linux systems.

To compile the code, run

    go generate ./...
    go build -o fakeshop cmd/fakeshop.go

at the repository root.

To see the binary's options, run

    ./fakeshop -help

## Next Steps
- [ ] Make carts thread-safe
- [ ] Improve test coverage
- [ ] Refactor for better readability
- [ ] Add goroutine to expire carts and release stock back to inventory
- [ ] Notify user of expiry (preferably via websocket)