package main

import (
	"crypto/rand"
	//"fmt"
	"math/big"
	"testing"
)

var (
	testRandReader  = rand.Reader
	deviationMax    = big.NewInt(20)
	higherDeviation = big.NewInt(100)

	deviationMaxFloat    = big.NewFloat(20)
	higherDeviationFloat = big.NewFloat(100)
)

// func TestEventIsAnomalous(t *testing.T) {
// 	dataSet := fakeFixedData()

// 	ad := NewAnomalyDetection(dataSet...)

// 	ad.EventIsAnomalous(*big.NewFloat(6020))

// }

func TestSmallDataSet(t *testing.T) {

	dataSet := fakeSmallData()

	ad := NewAnomalyDetection(dataSet...)

	totalSamples, _ := ad.TotalSamples.Uint64()
	if totalSamples != 20 {
		t.Fatal("Wrong number of element in the set. There are for sure 20, instead we got:", totalSamples)
	}

	totalSum, _ := ad.totalSum.Uint64()
	if totalSum != 83 {
		t.Fatal("Total sum does not add up, should be 83, we got:", totalSum)
	}

	mean, _ := ad.mean.Float64()
	if mean != 4.15 {
		t.Fatal("Mean is wrong, should be 4.15, we got:", mean)
	}

	deviation, _ := ad.deviation.Float64()
	if deviation != 1.3199999999999996 {
		t.Error("Deviation is wrong, should be 1.32 (or 1.3199999999999996), we got:", deviation)
	}

	variance, _ := ad.variance.Float64()
	if variance != 2.3325000000000005 {
		t.Error("Variance is wrong, should be 2.3325 (or 2.3325000000000005), we got:", variance)
	}

	anomaly, result := ad.EventIsAnomalous(*big.NewFloat(4.3), big.NewFloat(0.02))
	if anomaly {
		t.Errorf("5.3 should not be an anomaly !!! %f\n", result)
	}

	// stricter threshold
	anomaly, _ = ad.EventIsAnomalous(*big.NewFloat(7.3), big.NewFloat(0.1))
	if !anomaly {
		t.Error("7.3 should be an anomaly !!!")
	}

	ad.ExpandDataSet(fakeSmallData()...)
	ad.ExpandDataSet(fakeSmallData()...)

	anomaly, _ = ad.EventIsAnomalous(*big.NewFloat(5.3), big.NewFloat(0.01))
	if anomaly {
		t.Error("5.3 should not be an anomaly !!!")
	}

	// stricter threshold
	anomaly, _ = ad.EventIsAnomalous(*big.NewFloat(7.3), big.NewFloat(0.1))
	if !anomaly {
		t.Error("7.3 should be an anomaly !!!")
	}

}

func TestFixedEstimateMean(t *testing.T) {

	dataSet := fakeFixedData()

	ad := NewAnomalyDetection(dataSet...)

	mean, _ := ad.mean.Float64()
	if mean != 5052.0320320320325 {
		t.Fatal("With the current code the mean should be 5052.0320320320325, instead we got", mean)
	}

	variance, _ := ad.variance.Float64()
	//fmt.Println("Variance:", variance)

	if variance != 1536.5114864614086 {
		t.Fatal("With the current code the variance should be 1536.5114864614086, instead we got", variance)
	}

	//  threshold: 0.1%
	anomaly, result := ad.EventIsAnomalous(*big.NewFloat(5050), big.NewFloat(0.001))
	if anomaly {
		t.Errorf("5050 should NOT be an anomaly !!! %f\n", result)
	}

}

// // =======================================================================================

func TestRandomEstimateMean(t *testing.T) {

	dataSet := fakeRandomData()

	ad := NewAnomalyDetection(dataSet...)
	ad.ExpandDataSet(dataSet...)

	ad.estimateMean()

	ad.estimateVariance()

	//
	anomaly, result := ad.EventIsAnomalous(*big.NewFloat(4982), big.NewFloat(0.001))
	if anomaly {
		t.Errorf("4982 should NOT be an anomaly !!! %f\n", result)
	}

	anomaly, result = ad.EventIsAnomalous(*big.NewFloat(5117), big.NewFloat(0.001))
	if anomaly {
		t.Errorf("5117 should NOT be an anomaly !!! %f\n", result)
	}

}

func fakeSmallData() []big.Float {

	// total number of elements: 20
	// sum: 30 + 53 = 83
	// mean: 4.15
	// deviation: 1.32
	// variance: 2.3325

	// 3.15 + | 9.9225 +
	// 2.15 + | 4.6225 +
	// 1.15 + | 1.3225 +
	// 0.15 + | 0.0225 +
	// 0.85 + | 0.7225 +
	// -----  | ------
	// 7.45 *2 = 14.9  |  16.6125 * 2 = 33.225

	// 0.95 + | 0.9025 +
	// 1.05 + | 1.1025 +
	// 1.15 + | 1.3225 +
	// 1.25 + | 1.5625 +
	// 1.35 + | 1.8225 +
	// -------| ------
	// 5.75 * 2 = 11.5 |  6.7125 * 2 = 13.425
	// total = 14.9 + 11.5 = 26.4  |  33.225 + 13.425 = 46.65

	//

	dataSet := []big.Float{
		*big.NewFloat(1),
		*big.NewFloat(2),
		*big.NewFloat(3),
		*big.NewFloat(4),
		*big.NewFloat(5),
		*big.NewFloat(5.1),
		*big.NewFloat(5.2),
		*big.NewFloat(5.3),
		*big.NewFloat(5.4),
		*big.NewFloat(5.5),
		*big.NewFloat(5.5),
		*big.NewFloat(5.4),
		*big.NewFloat(5.3),
		*big.NewFloat(5.2),
		*big.NewFloat(5.1),
		*big.NewFloat(5),
		*big.NewFloat(4),
		*big.NewFloat(3),
		*big.NewFloat(2),
		*big.NewFloat(1),
	}

	return dataSet

}

// helper function to fake fixed data
func fakeFixedData() []big.Float {

	var dataSet []big.Float

	baseInt := big.NewFloat(5000)
	//defaultInt := big.NewInt(5000)
	// fake sample data
	var increment big.Float
	for i := 0; i < 999; i++ {
		baseInt = big.NewFloat(5000)
		if i > 200 && i < 800 {
			increment = *baseInt.Add(baseInt, deviationMaxFloat)
			dataSet = append(dataSet, increment)
			//fmt.Println("_", baseInt.Uint64())

		} else {
			increment = *baseInt.Add(baseInt, higherDeviationFloat)
			dataSet = append(dataSet, increment)
			//fmt.Println(".", baseInt.Uint64())
		}
	}
	return dataSet

}

// ======== ******************* HELPER FUNCTIONS FOR TESTING

// helper function to fake fixed data
func fakeRandomData() []big.Float {

	var dataSet []big.Float

	baseInt := big.NewFloat(5000)
	// fake sample data
	for i := 0; i < 100000; i++ {
		baseInt = big.NewFloat(5000)
		if i > 200 && i < 800 {
			dataSet = append(dataSet, *baseInt.Add(baseInt, getDeviation(deviationMax)))
		} else {
			dataSet = append(dataSet, *baseInt.Add(baseInt, getDeviation(higherDeviation)))
		}
	}
	return dataSet

}

// Helper function to randomize data
func getDeviation(max *big.Int) *big.Float {
	deviation, _ := rand.Int(testRandReader, max)
	return big.NewFloat(float64(deviation.Int64()))
}
