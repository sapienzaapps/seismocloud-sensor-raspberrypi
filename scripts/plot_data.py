import matplotlib.pyplot as plt
import numpy as np
import sys
import progressbar
import math


def runningavg(serie):
	n = 0
	avg = 0
	variance = 0
	for x in serie:
		n += 1
		delta = x - avg
		avg += delta / float(n)
		variance = variance + delta*(x-avg)
	if n == 0:
		return 0, 0
	elif n == 1:
		return avg, 0
	return avg, math.sqrt(variance / float(n-1))


print("Ingesting data")
X = []
Y = []
Z = []
Vect = []
OldThreshold = []
NewThresholdValue = 0
NewThresholdVar = 0
NewThreshold = []
fp = open(sys.argv[1], "r")
rows = fp.read().split("\n")
with progressbar.ProgressBar(max_value=len(rows)) as pbar:
	for r in rows:
		if r.count("\t") == 0:
			pbar.update(pbar.value + 1)
			continue
		cells = r.split("\t")
		X.append(float(cells[0]))
		Y.append(float(cells[1]))
		Z.append(float(cells[2]))
		Vect.append(float(cells[3]))

		avg, variance = runningavg(Vect)
		OldThreshold.append(avg + variance*6)

		if len(Vect) > 2000 and NewThresholdValue == 0:
			NewThresholdValue, NewThresholdVar = runningavg(Vect)
			NewThresholdValue += NewThresholdVar*6
			NewThreshold.append(NewThresholdValue)
		elif NewThresholdValue != 0:
			NewThreshold.append(NewThresholdValue)
		else:
			NewThreshold.append(0)

		pbar.update(pbar.value + 1)

timespan = range(0, len(X))

print(runningavg(Vect))

print("Drawing")
#plt.plot(timespan, X, label="X")
#plt.plot(timespan, Y, label="Y")
#plt.plot(timespan, Z, label="X")
plt.plot(timespan, Vect, label="Vect")
plt.plot(timespan, OldThreshold, label="OldThreshold")
plt.plot(timespan, NewThreshold, label="NewThreshold")
plt.show()
