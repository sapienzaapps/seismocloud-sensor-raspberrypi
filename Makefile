.PHONY: clean raspi raspi-prod all phi test

seismosensor: *.go
	go build -o seismosensor

raspi-prod:
	GOARCH=arm GOARM=6 go build -tags "rpi adxl345 prod" -o seismosensor

raspi:
	GOARCH=arm GOARM=6 go build -tags "rpi adxl345" -o seismosensor

phi:
	go build -tags "phidget" -o seismosensor

clean:
	rm -f seismosensor

test:
	go test -v ./...
	go vet ./...
	gosec ./...
	staticcheck ./...
	ineffassign .
	errcheck ./...
