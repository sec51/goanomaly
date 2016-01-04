[![Build Status](https://travis-ci.org/sec51/goanomaly.svg)](https://travis-ci.org/sec51/goanomaly) [![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://godoc.org/github.com/sec51/goanomaly/)

### Anomaly detection library in golang, implemented via Gaussian distribution

1. Choose feature `x(i)` that might be indicative of anomalous examples
2. Estimate parameters by calculating the `mean` and `standard deviation`
3. Given a new example `x`, compute `p(x)` via the Gaussian normal distribution formaula
4. Anomaly if `p(x) < k`

### Example

````
import (
	"goanomaly"
	"log"
	"math/big"
)

// ==============================================================

// get a set of data
dataSet := fakeFixedData()

// init the AnomalyDetection object
anomalyDetection := goanomaly.NewAnomalyDetection(dataSet...)

// Call EventIsAnomalous with a new data point, which you can test against, and a threshold 
anomaly, result := anomalyDetection.EventIsAnomalous(*big.NewFloat(5050), big.NewFloat(0.001))

if anomaly {
	log.Println("Data point", 5050, "is anomalous")
}

// ==============================================================

// helper function to fake fixed data
func fakeFixedData() []big.Float {

	var dataSet []big.Float

	baseInt := big.NewFloat(5000)
	
	// fake sample data
	var increment big.Float
	for i := 0; i < 999; i++ {
		baseInt = big.NewFloat(5000)
		if i > 200 && i < 800 {
			increment = *baseInt.Add(baseInt, deviationMaxFloat)
			dataSet = append(dataSet, increment)		

		} else {
			increment = *baseInt.Add(baseInt, higherDeviationFloat)
			dataSet = append(dataSet, increment)			
		}
	}
	return dataSet

}

````

### Important

The Gaussian distribution anomaly detection usually does not keep into consideration the relation between the features.
For instance if you want to have a relation between say CPU load(x1) and Memory usage(x2) then you should create a 3rd feature in this way: (x1)/(x2)
and add it to the dataset

===

### LICENSE

Copyright (c) 2015-2016 Sec51.com <info@sec51.com>

Permission to use, copy, modify, and distribute this software for any
purpose with or without fee is hereby granted, provided that the above
copyright notice and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
