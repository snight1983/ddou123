package vpr

//https://github.com/liuxp0827/govpr.git

import (
	"fmt"
	"math"
)

type gMM struct {
	Frames          int         // number of total frames
	FeatureData     [][]float32 // feature buffer
	VectorSize      int         // Vector size											1
	Mixtures        int         // Mixtures of the GMM	 								1
	deterCovariance []float64   // determinant of the covariance matrix [mixture]		1
	MixtureWeight   []float64   // weight of each mixture[mixture]						1
	Mean            [][]float64 // mean vector [mixture,dimension]						1
	Covar           [][]float64 // covariance (diagonal) [mixture,dimension]			1
}

func newGMM() *gMM {
	gmm := &gMM{
		FeatureData:   make([][]float32, 0),
		MixtureWeight: make([]float64, 0),
		Mean:          make([][]float64, 0),
		Covar:         make([][]float64, 0),
	}
	return gmm
}

func (g *gMM) copy(gmm *gMM) {
	g.Frames = gmm.Frames
	g.VectorSize = gmm.VectorSize
	g.Mixtures = gmm.Mixtures

	g.FeatureData = make([][]float32, gmm.Frames, gmm.Frames)
	g.deterCovariance = make([]float64, g.Mixtures, g.Mixtures)
	g.MixtureWeight = make([]float64, g.Mixtures, g.Mixtures)
	g.Mean = make([][]float64, g.Mixtures, g.Mixtures)
	g.Covar = make([][]float64, g.Mixtures, g.Mixtures)

	for i := 0; i < g.Frames; i++ {
		g.FeatureData[i] = make([]float32, g.VectorSize, g.VectorSize)
	}

	for i := 0; i < g.Frames; i++ {
		for j := 0; j < g.VectorSize; j++ {
			g.FeatureData[i][j] = gmm.FeatureData[i][j]
		}
	}

	for i := 0; i < g.Mixtures; i++ {
		g.deterCovariance[i] = gmm.deterCovariance[i]
		g.MixtureWeight[i] = gmm.MixtureWeight[i]
		g.Mean[i] = make([]float64, g.VectorSize, g.VectorSize)
		g.Covar[i] = make([]float64, g.VectorSize, g.VectorSize)
		for j := 0; j < g.VectorSize; j++ {
			g.Mean[i][j] = gmm.Mean[i][j]
			g.Covar[i][j] = gmm.Covar[i][j]
		}
	}
}

func (g *gMM) dupModel(gmm *gMM) {
	g.Mixtures = gmm.Mixtures
	g.VectorSize = gmm.VectorSize
	g.deterCovariance = make([]float64, g.Mixtures, g.Mixtures)
	g.MixtureWeight = make([]float64, g.Mixtures, g.Mixtures)
	g.Mean = make([][]float64, g.Mixtures, g.Mixtures)
	g.Covar = make([][]float64, g.Mixtures, g.Mixtures)

	for i := 0; i < g.Mixtures; i++ {
		g.deterCovariance[i] = gmm.deterCovariance[i]
		g.MixtureWeight[i] = gmm.MixtureWeight[i]
		g.Mean[i] = make([]float64, g.VectorSize, g.VectorSize)
		g.Covar[i] = make([]float64, g.VectorSize, g.VectorSize)
		for j := 0; j < g.VectorSize; j++ {
			g.Mean[i][j] = gmm.Mean[i][j]
			g.Covar[i][j] = gmm.Covar[i][j]
		}
	}
}

func (g *gMM) loadModel(filename string) error {
	reader, err := newVPerile(filename)
	if err != nil {
		return err
	}

	g.Mixtures, err = reader.getInt()
	if err != nil {
		return err
	}

	g.VectorSize, err = reader.getInt()
	if err != nil {
		return err
	}

	g.deterCovariance = make([]float64, g.Mixtures, g.Mixtures)
	g.MixtureWeight = make([]float64, g.Mixtures, g.Mixtures)
	g.Mean = make([][]float64, g.Mixtures, g.Mixtures)
	g.Covar = make([][]float64, g.Mixtures, g.Mixtures)
	for i := 0; i < g.Mixtures; i++ {
		g.Mean[i] = make([]float64, g.VectorSize, g.VectorSize)
		g.Covar[i] = make([]float64, g.VectorSize, g.VectorSize)
	}

	for i := 0; i < g.Mixtures; i++ {
		g.deterCovariance[i] = 0.0
	}

	for i := 0; i < g.Mixtures; i++ {
		g.MixtureWeight[i], err = reader.getFloat64()
		if err != nil {
			//log.Error(err)
			return err
		}
	}

	for i := 0; i < g.Mixtures; i++ {
		_, err = reader.getFloat64() // not used
		if err != nil {
			//log.Error(err)
			return err
		}

		_, err = reader.getFloat64() // not used
		if err != nil {
			return err
		}

		_, err = reader.getByte() // not used
		if err != nil {
			return err
		}

		for j := 0; j < g.VectorSize; j++ {
			g.Covar[i][j], err = reader.getFloat64()
			if err != nil {
				return err
			}

			g.deterCovariance[i] += math.Log(g.Covar[i][j])
		}

		for j := 0; j < g.VectorSize; j++ {
			g.Mean[i][j], err = reader.getFloat64()
			if err != nil {
				return err
			}
		}
	}
	reader.close()
	return nil
}

func (g *gMM) SaveModel(filename string) error {
	writer, err := newVPerile(filename)
	if err != nil {
		return err
	}
	defer writer.close()

	_, err = writer.putInt(g.Mixtures)
	if err != nil {
		return err
	}

	_, err = writer.putInt(g.VectorSize)
	if err != nil {
		return err
	}

	for i := 0; i < g.Mixtures; i++ {
		_, err = writer.putFloat64(g.MixtureWeight[i])
		if err != nil {
			return err
		}
	}

	for i := 0; i < g.Mixtures; i++ {
		_, err = writer.putFloat64(0.0) // not used
		if err != nil {
			return err
		}

		_, err = writer.putFloat64(0.0) // not used
		if err != nil {
			return err
		}

		err = writer.putByte(byte(0)) // not used
		if err != nil {
			return err
		}

		for j := 0; j < g.VectorSize; j++ {
			_, err = writer.putFloat64(g.Covar[i][j])
			if err != nil {
				return err
			}
		}

		for j := 0; j < g.VectorSize; j++ {
			_, err = writer.putFloat64(g.Mean[i][j])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *gMM) copyFeatureData(gmm *gMM) error {

	g.Frames = gmm.Frames
	g.VectorSize = gmm.VectorSize
	g.FeatureData = make([][]float32, gmm.Frames, gmm.Frames)
	for i := 0; i < g.Frames; i++ {
		g.FeatureData[i] = make([]float32, g.VectorSize, g.VectorSize)
	}

	for i := 0; i < g.Frames; i++ {
		for j := 0; j < g.VectorSize; j++ {
			g.FeatureData[i][j] = gmm.FeatureData[i][j]
		}
	}

	return nil
}

func (g *gMM) eM(mixtures int) (int, error) {
	var dlogfrmprob, rubbish, lastrubbish float64
	var dsumgama, dlogmixw, dgama, dmixw []float64
	var threshold float64 = 1e-5
	var mean, covar [][]float64
	var loop int = 0

	mean = make([][]float64, mixtures, mixtures)
	covar = make([][]float64, mixtures, mixtures)
	for i := 0; i < mixtures; i++ {
		mean[i] = make([]float64, g.VectorSize, g.VectorSize)
		covar[i] = make([]float64, g.VectorSize, g.VectorSize)
	}

	dmixw = make([]float64, mixtures, mixtures)
	dlogmixw = make([]float64, mixtures, mixtures)
	dgama = make([]float64, mixtures, mixtures)
	dsumgama = make([]float64, mixtures, mixtures)
	rubbish = .0

	doing := func() (int, error) {
		lastrubbish = rubbish
		rubbish = .0
		for i := 0; i < mixtures; i++ {

			// speed up
			if g.MixtureWeight[i] <= 0 {
				dlogmixw[i] = LOGZERO
			} else {
				dlogmixw[i] = math.Log(g.MixtureWeight[i])
			}

			// clean up temporary values
			dmixw[i] = .0
			dsumgama[i] = .0
			for j := 0; j < g.VectorSize; j++ {
				mean[i][j] = .0
				covar[i][j] = .0
			}
		}

		for i := 0; i < g.Frames; i++ {
			dlogfrmprob = LOGZERO
			for j := 0; j < mixtures; j++ {
				dgama[j] = g.lMixProb(g.FeatureData[i], j)
				dgama[j] += dlogmixw[j]
				dlogfrmprob = g.logAdd(dgama[j], dlogfrmprob)
			}

			rubbish += dlogfrmprob

			for j := 0; j < mixtures; j++ {
				dgama[j] -= dlogfrmprob
				dgama[j] = math.Exp(dgama[j])
				dsumgama[j] += dgama[j]

				// update weights
				dmixw[j] += dgama[j]
				for k := 0; k < g.VectorSize; k++ {
					mean[j][k] += dgama[j] * float64(g.FeatureData[i][k])
					covar[j][k] += dgama[j] * float64(g.FeatureData[i][k]) * float64(g.FeatureData[i][k])
				}
			}
		}

		rubbish /= float64(g.Frames)

		for i := 0; i < mixtures; i++ {
			if dsumgama[i] == .0 {
				return -1, nil
			}

			g.MixtureWeight[i] = dmixw[i] / float64(g.Frames)

			for j := 0; j < g.VectorSize; j++ {
				g.Mean[i][j] = mean[i][j] / dsumgama[i]
				g.Covar[i][j] = covar[i][j] / dsumgama[i]
				g.Covar[i][j] -= g.Mean[i][j] * g.Mean[i][j]

				if g.Covar[i][j] < VARFLOOR {
					g.Covar[i][j] = VARFLOOR
				}

				if g.Covar[i][j] > VARCEILING {
					g.Covar[i][j] = VARCEILING
				}
			}
		}

		for i := 0; i < mixtures; i++ {
			g.deterCovariance[i] = .0
			for j := 0; j < g.VectorSize; j++ {
				g.deterCovariance[i] += math.Log(g.Covar[i][j])
			}
		}
		loop++
		return 0, nil
	}

DO:
	ret, err := doing()
	if err != nil {
		return 0, err
	} else if ret == -1 {
		return 0, fmt.Errorf("error train loop")
	}

	for loop < MAXLOOP && math.Abs((rubbish-lastrubbish)/(lastrubbish+0.01)) > threshold {
		goto DO
	}

	return loop, nil
}

func (g *gMM) lProb(featureBuf [][]float32, start, length int64) float64 {
	var dgama, dlogfrmprob, sum float64 = .0, .0, .0
	dlogmixw := make([]float64, g.Mixtures, g.Mixtures)
	for ij := 0; ij < g.Mixtures; ij++ {
		if g.MixtureWeight[ij] <= 0 {
			dlogmixw[ij] = LOGZERO
		} else {
			dlogmixw[ij] = math.Log(g.MixtureWeight[ij])
		}
	}

	for ii := int64(start); ii < (start + length); ii++ {
		dlogfrmprob = LOGZERO
		for jj := 0; jj < g.Mixtures; jj++ {
			dgama = g.lMixProb(featureBuf[ii], jj) + dlogmixw[jj]
			dlogfrmprob = g.logAdd(dgama, dlogfrmprob)
		}
		sum += dlogfrmprob
	}

	dlogmixw = nil
	return sum
}

// LogAdd Routine for adding two log-values in a linear scale, return the log
func (g *gMM) logAdd(lvar1 float64, lvar2 float64) float64 {
	var diff, z float64
	var minLogExp = -math.Log(-(LOGZERO))
	if lvar1 < lvar2 {
		lvar1, lvar2 = lvar2, lvar1
	}

	diff = lvar2 - lvar1
	if diff < minLogExp {
		if lvar1 < LSMALL {
			return LOGZERO
		}
		return lvar1

	}
	z = math.Exp(diff)
	return lvar1 + math.Log(1.0+z)

}

// LMixProb Return the log likelihood of the given vector to the given
// mixture (the score should be multiplied with the mixture weight).
// FeatureData :  	pointer to the feature vector in float
// MixIndex :  		Mixture index
// Return Value: 	The log-likelihood in float64
func (g *gMM) lMixProb(buffer []float32, mixIndex int) float64 {
	if g.Mean == nil || g.deterCovariance == nil {
		panic("Model not loaded")
	}

	var dsum, dtmp float64
	dsum = .0
	for ii := 0; ii < g.VectorSize; ii++ {
		dtmp = float64(buffer[ii]) - g.Mean[mixIndex][ii]
		dsum += dtmp * dtmp / g.Covar[mixIndex][ii]
	}

	dsum = -(float64(g.VectorSize) * DLOG2PAI) - g.deterCovariance[mixIndex] - dsum
	dsum /= 2

	if dsum < LOGZERO {
		return LOGZERO
	}
	return dsum

}