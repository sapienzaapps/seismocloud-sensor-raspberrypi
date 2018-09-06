.PHONY: clean raspi raspi-debug all

all:
	$(info Targets: raspi, raspi-debug, clean)

raspi:
	GOARCH=arm GOARM=6 go build -tags rpi -o seismosensor

raspi-debug:
	GOARCH=arm GOARM=6 go build -tags "rpi debug" -o seismosensor

clean:
	rm -f seismocloud-sensor-raspberrypi