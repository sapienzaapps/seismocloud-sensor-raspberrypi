.PHONY: clean raspi all

all:
	$(info Targets: raspi)

raspi:
	GOARCH=arm GOARM=6 go build -tags rpi -o seismosensor

clean:
	rm -f seismocloud-sensor-raspberrypi