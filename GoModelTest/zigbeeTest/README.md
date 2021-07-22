
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build zigbee.go

./zigbee -read true

./zigbee -write true
