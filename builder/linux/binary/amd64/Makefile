shelly-client:
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -trimpath -o shelly-client ../../../../shelly-client.go

clean:
	$(RM) -r shelly-client
