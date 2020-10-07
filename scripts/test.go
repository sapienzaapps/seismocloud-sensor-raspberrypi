package main

import (
	"io/ioutil"
	"os"
	"fmt"
	"strings"
	"strconv"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-sensor-raspberrypi/utils"
)

func main() {
	fcontent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	} else {
		runAvg := utils.NewRunningAvgFloat64()
		for _, row := range strings.Split(string(fcontent), "\n") {
			values := strings.Split(row, "\t")
			if len(values) != 4 {
				continue
			}
			vect, err := strconv.ParseFloat(values[3], 64)
			if err != nil {
				panic(err)
			}
			runAvg.AddValue(vect)
		}

		fmt.Println(runAvg.GetAverage())
		fmt.Println(runAvg.GetStandardDeviation())
	}
}