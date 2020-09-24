# Raspberry Pi SeismoCloud client

**Important**: this code is incomplete.

## To Do list:

* [ ] Check for updates at startup
* [ ] Test update procedure
* [ ] Add support for new local discovery protocol
* [ ] Add support for sensor data stream
* [ ] Add support for set probe speed command
* [ ] Add support for statistics publish (WiFI SSID, RSSI, etc)
* [X] Implement new threshold algorithm (see NodeMCU code)
* [ ] Handle errors
* [X] Make the sensor/seismometer (platform-agnostic) similar to the NodeMCU/Android one (e.g. rename variables,
refactor control flow, etc)

## Phidget instructions

Requirements:
```sh
apt-get install libusb-1.0-0-dev
```

Then download the [libphidget22](https://www.phidgets.com/downloads/phidget22/libraries/linux/libphidget22.tar.gz) file. Compile it.