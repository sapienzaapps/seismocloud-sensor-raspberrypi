.PHONY: clean raspi raspi-debug all

seismosensor: *.go
	go build -o seismosensor

raspi:
	GOARCH=arm GOARM=6 go build -tags rpi -o seismosensor

raspi-debug:
	GOARCH=arm GOARM=6 go build -tags "rpi debug" -o seismosensor

clean:
	rm -f seismosensor

test:
	go test -v ./...
	go vet ./...
	gosec ./...
	staticcheck ./...
	ineffassign .
	errcheck ./...
