package goanomaly

import (
	"math"
	"math/big"
	"sync"
)

var (
	pi   = big.NewFloat(math.Pi)
	zero = big.NewFloat(float64(0))
	one  = big.NewFloat(float64(1))
	two  = big.NewFloat(float64(2))
	e    = big.NewFloat(math.E)

	// NOT USED AT THE MOMENT
	delimiter = big.NewFloat(float64(100)) // this is the delimiter for between a small set and a large set (does not change much)

	// pre calculated constants
	doublePi         = new(big.Float).Mul(two, pi)
	doublePiValue, _ = doublePi.Float64()
	doublePiSqrt     = big.NewFloat(math.Sqrt(doublePiValue))
)

type AnomalyDetection struct {
	dataSet      []big.Float
	totalSamples big.Float
	totalSum     big.Float
	mean         big.Float // this is the average mean
	variance     big.Float // this is the average variance
	deviation    big.Float // this is the average deviation
}

type AnomalyDetectionVector []*AnomalyDetection

// Creates an anomaly detection object with multi dimension dataset (multivariate)
func NewAnomalyDetectionVector(vector ...[]big.Float) AnomalyDetectionVector {

	var adVector AnomalyDetectionVector

	// wait group
	var wg sync.WaitGroup

	// mutex used to append
	var initVectorMutex sync.Mutex

	for _, data := range vector {
		//adv = append(adv, data)
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(m sync.Mutex, anomalyVector *AnomalyDetectionVector, set ...big.Float) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			// init the anomaly detection object: this is the expensive call, based on the dataset
			anomalyDetection := NewAnomalyDetection(set...)

			// lock the mutex and append
			m.Lock()
			// de-reference pointer
			adv := *anomalyVector
			// append
			adv = append(adv, anomalyDetection)
			// set the new pointer
			anomalyVector = &adv
			// unlock the mutex
			m.Unlock()

		}(initVectorMutex, &adVector, data...)
	}

	// wait for thego routines to finish
	wg.Wait()

	return adVector
}

func (adVector AnomalyDetectionVector) EventIsAnomalous(eventX big.Float, threshold *big.Float) (bool, float64) {

	// wait group
	var wg sync.WaitGroup

	var singleProbabilities []*big.Float

	// mutex used to append
	var probabilityMutex sync.Mutex

	for _, ad := range adVector {
		//adv = append(adv, data)
		// Increment the WaitGroup counter.
		wg.Add(1)
		// Launch a goroutine to fetch the URL.
		go func(m sync.Mutex, anomaly *AnomalyDetection, prob *[]*big.Float, eX big.Float) {

			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			// calculate the probability
			probability := anomaly.calculateProbability(eX)

			// lock the mutex and append
			m.Lock()
			// de-reference pointer
			p := *prob
			// append
			p = append(p, probability)
			// set the new pointer
			prob = &p
			// unlock the mutex
			m.Unlock()

		}(probabilityMutex, ad, &singleProbabilities, eventX)
	}

	// multiply all the probabilities together
	totalProbability := big.NewFloat(1)
	for _, probability := range singleProbabilities {
		totalProbability.Mul(totalProbability, probability)
	}

	// get the float64 form the total probability
	r, _ := totalProbability.Float64()

	// if the total probability is lower than the threshold then the event is anomalous
	return totalProbability.Cmp(threshold) < 0, r

}

// Creates an anomaly detection object with a one dimension dataset
func NewAnomalyDetection(data ...big.Float) *AnomalyDetection {

	ad := AnomalyDetection{}
	ad.dataSet = data
	ad.totalSamples.SetFloat64(float64(len(data)))

	// estimate the mean already
	ad.estimateMean()

	// estimate variance
	ad.estimateVariance()

	return &ad
}

func (ad *AnomalyDetection) ExpandDataSet(data ...big.Float) {
	if ad.dataSet == nil {
		// Should return an error because the dataSet was cleared
		return;
	}

	ad.totalSum = *new(big.Float)

	// means totalSamples is smaller than delimiter
	totalSamples := big.NewFloat(float64(len(data)))
	// if totalSamples.Cmp(delimiter) < 0 {
	// 	totalSamples.Sub(totalSamples, one)
	// }
	ad.totalSamples = *totalSamples

	ad.estimateMean()

	ad.estimateVariance()
}

// ClearDataSet reset the dataSet to nil to release resources
func (ad *AnomalyDetection) ClearDataSet() {
	ad.dataSet = nil
}

// This method calculates the probability with probability density formula
// TODO: CREATE THE SQRT and EXP methods for bignum
func (ad *AnomalyDetection) calculateProbability(eventX big.Float) *big.Float {
	// Right term
	rightTerm := new(big.Float).Sub(&eventX, &ad.mean) // no need to take the Absolute value because we square on the next step
	rightTerm.Mul(rightTerm, rightTerm)

	// take the variance and double it
	tempA := new(big.Float).Mul(two, &ad.variance)

	// divide eventXDeviationSquared with doubleVariance
	rightTerm.Quo(rightTerm, tempA)

	// get its value
	rightFloat, _ := rightTerm.Float64()

	// do e^(-right term value)
	rightFloat = math.Exp(-rightFloat)

	// ======================================
	// Left term

	// Init the holder of the final value of the first term
	// Multiply the Square root of the the 2*pi for the deviation
	tempA.Mul(doublePiSqrt, &ad.deviation)

	// multiply the two terms
	return rightTerm.Quo(rightTerm.SetFloat64(rightFloat), tempA)
}

// Verifies whether a specific event X is anomalous or not
func (ad *AnomalyDetection) EventIsAnomalous(eventX big.Float, threshold *big.Float) (bool, float64) {

	probability := ad.calculateProbability(eventX)
	r, _ := probability.Float64()
	return probability.Cmp(threshold) < 0, r
}


// Verifies whether a specific event X is anomalous or not
// This method calculates the probability with probability density formula
// TODO: CREATE THE SQRT and EXP methods for bignum
func (ad *AnomalyDetection) EventXIsAnomalous(eventX, threshold *big.Float) (bool, *big.Float) {
	// Right term
	rightTerm := new(big.Float).Sub(eventX, &ad.mean) // no need to take the Absolute value because we square on the next step
	rightTerm.Mul(rightTerm, rightTerm)

	// take the variance and double it
	tempA := new(big.Float).Mul(two, &ad.variance)

	if tempA.Cmp(zero) == 0 {
		return false, zero
	}

	// divide eventXDeviationSquared with doubleVariance
	rightTerm.Quo(rightTerm, tempA)

	// get its value
	rightFloat, _ := rightTerm.Float64()

	// do e^(-right term value)
	rightFloat = math.Exp(-rightFloat)

	// ======================================
	// Left term

	// Init the holder of the final value of the first term
	// Multiply the Square root of the the 2*pi for the deviation
	tempA.Mul(doublePiSqrt, &ad.deviation)

	// multiply the two terms
	rightTerm.Quo(rightTerm.SetFloat64(rightFloat), tempA)

	return rightTerm.Cmp(threshold) < 0, rightTerm
}

// Estimates the Mean based on the data set
// If the data set is relatively small (< 1000 examples), then remove 1 from the total
func (ad *AnomalyDetection) estimateMean() *big.Float {

	// initialize the total to zero
	totalSum := new(big.Float)

	// Loop thorugh the data set
	for _, element := range ad.dataSet {

		// sum up its elements
		totalSum.Add(totalSum, &element)
		//e, _ := element.Float64()

	}

	// make a copy of the total sum and assign it to the anomaly detection object
	ad.totalSum.Copy(totalSum)

	// calculate the mean and return
	return ad.mean.Quo(totalSum, &ad.totalSamples)
}

// Estimates the Variance based on the data set
// If the data set is relatively small (< 1000 examples), then remove 1 from the total
func (ad *AnomalyDetection) estimateVariance() *big.Float {

	// this means that the mean was never calculated before, therefore do it now
	// the means is needed for the cimputation of the deviation
	if ad.mean.Cmp(zero) == 0 {
		ad.estimateMean()
	}

	// initialize the total to zero
	totalVariance := new(big.Float)
	totalDeviation := new(big.Float)

	deviation := new(big.Float)
	var singleVariance *big.Float

	// Loop while a is smaller than 1e100.
	for _, element := range ad.dataSet {
		// first calculate the deviation for each element, by subtracting the mean, take the absolute value
		deviation.Sub(&element, &ad.mean).Abs(deviation)

		// add it to the total
		totalDeviation.Add(totalDeviation, deviation)

		// calculate the variance by squaring it
		singleVariance = deviation.Mul(deviation, deviation) // ^2

		// the calculate the variance
		totalVariance.Add(totalVariance, singleVariance)
	}

	// calculate the variance
	// assign the variance to the anomaly detection object
	ad.variance = *totalVariance.Quo(totalVariance, &ad.totalSamples)

	// calculate the deviation
	ad.deviation = *totalDeviation.Quo(totalDeviation, &ad.totalSamples)

	return &ad.variance
}
