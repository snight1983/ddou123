package vpr

//https://github.com/liuxp0827/govpr.git
import (
	"fmt"
)

type parameter struct {
	lowCutOff              int  // low cut-off
	highCutOff             int  // high cut-off
	filterBankSize         int  // num of filter-bank
	frameLength            int  // frame length
	frameShift             int  // frame shift
	mfccOrder              int  // mfcc order
	isStatic               bool // static mfcc
	isDynamic              bool // dynamic mfcc
	isAcce                 bool // acce mfcc
	cmsvn                  bool // cmsvn
	isZeroGlobalMean       bool // zero global mean
	isDBNorm               bool // decibel normalization
	isDiffPolish           bool // polish differential formula
	isDiffPowerSpectrum    bool // differentail power spectrum
	isPredDiffAmplSpectrum bool // predictive differential amplitude spectrum
	isEnergyNorm           bool
	silFloor               int16
	energyscale            int16
	isFeatWarping          bool
	featWarpWinSize        int16
	isRasta                bool
	rastaCoff              float64
}

func extract(data []int16, gmm *gMM) error {
	var p, para []float32
	var info wavInfo
	var cp *cParam = newCParam()
	var pm *parameter = new(parameter)
	var err error
	var icol, irow int
	var buflen int = len(data)

	p = make([]float32, buflen, buflen)
	for i := 0; i < buflen; i++ {
		p[i] = float32(data[i])
	}

	pm.lowCutOff = LOWCUTOFF
	pm.highCutOff = HIGHCUTOFF
	pm.filterBankSize = FILTERBANKSIZE
	pm.frameLength = FRAMELENGTH
	pm.frameShift = FRAMESHIFTt
	pm.mfccOrder = MFCCORDER

	pm.isStatic = BSTATIC
	pm.isDynamic = BDYNAMIC
	pm.isAcce = BACCE

	pm.cmsvn = CMSVN
	pm.isZeroGlobalMean = ZEROGLOBALMEAN
	pm.isDiffPolish = DIFPOL
	pm.isDiffPowerSpectrum = DPSCC
	pm.isPredDiffAmplSpectrum = PDASCC
	pm.isEnergyNorm = ENERGYNORM
	pm.silFloor = SILFLOOR

	pm.energyscale = ENERGYSCALE
	pm.isFeatWarping = FEATWARP
	pm.featWarpWinSize = FEATUREWARPINGWINSIZE
	pm.isDBNorm = DBNORM
	pm.isRasta = RASTA
	pm.rastaCoff = RASTACOFF

	info.SampleRate = SAMPLERATE
	info.Length = int64(buflen)
	info.BitSPSample = BITPERSAMPLE

	if pm.highCutOff > pm.lowCutOff {
		err = cp.initFBank2(info.SampleRate, pm.frameLength, pm.filterBankSize, pm.lowCutOff, pm.highCutOff)
	} else {
		err = cp.initFBank(info.SampleRate, pm.frameLength, pm.filterBankSize)
	}

	if err != nil {
		return err
	}

	err = cp.initMfcc(pm.mfccOrder, float32(pm.frameShift))
	if err != nil {
		return err
	}

	if pm.isStatic {
		cp.getMfcc().IsStatic = true
	} else {
		cp.getMfcc().IsStatic = false
	}

	if pm.isDynamic {
		cp.getMfcc().IsDynamic = true
	} else {
		cp.getMfcc().IsDynamic = false
	}

	if pm.isAcce {
		cp.getMfcc().IsAcce = true
	} else {
		cp.getMfcc().IsAcce = false
	}

	if pm.isZeroGlobalMean {
		cp.getMfcc().IsZeroGlobalMean = true
	} else {
		cp.getMfcc().IsZeroGlobalMean = false
	}

	if pm.isDBNorm {
		cp.getMfcc().IsDBNorm = true
	} else {
		cp.getMfcc().IsDBNorm = false
	}

	if pm.isDiffPolish {
		cp.getMfcc().IsPolishDiff = true
	} else {
		cp.getMfcc().IsPolishDiff = false
	}

	if pm.isDiffPowerSpectrum {
		cp.getMfcc().IsDiffPowerSpectrum = true
	} else {
		cp.getMfcc().IsDiffPowerSpectrum = false
	}

	if pm.isPredDiffAmplSpectrum {
		cp.getMfcc().IsPredDiffAmpSpetrum = true
	} else {
		cp.getMfcc().IsPredDiffAmpSpetrum = false
	}

	if pm.isEnergyNorm {
		cp.getMfcc().IsEnergyNorm = true
	} else {
		cp.getMfcc().IsEnergyNorm = false
	}

	if pm.isEnergyNorm {
		cp.getMfcc().SilFloor = pm.silFloor
	} else {
		cp.getMfcc().SilFloor = SILFLOOR
	}

	if pm.isEnergyNorm {
		cp.getMfcc().EnergyScale = pm.energyscale
	} else {
		cp.getMfcc().EnergyScale = ENERGYSCALE
	}

	if pm.isFeatWarping {
		cp.getMfcc().IsFeatWarping = true
	} else {
		cp.getMfcc().IsFeatWarping = false
	}

	if pm.isFeatWarping {
		cp.getMfcc().FeatWarpWinSize = pm.featWarpWinSize
	} else {
		cp.getMfcc().FeatWarpWinSize = FEATUREWARPINGWINSIZE
	}

	if pm.isRasta {
		cp.getMfcc().IsRasta = true
	} else {
		cp.getMfcc().IsRasta = false
	}

	cp.getMfcc().RastaCoff = pm.rastaCoff

	if nil != cp.wav2Mfcc(p, info, &para, &icol, &irow) && irow < MINFRAMES {
		return fmt.Errorf("Feature Extract error -2")
	}

	gmm.VectorSize = icol
	gmm.Frames = irow
	gmm.FeatureData = make([][]float32, gmm.Frames, gmm.Frames)
	for i := 0; i < gmm.Frames; i++ {
		gmm.FeatureData[i] = make([]float32, gmm.VectorSize, gmm.VectorSize)
	}

	for ii := 0; ii < irow; ii++ {
		for jj := 0; jj < icol; jj++ {
			gmm.FeatureData[ii][jj] = para[ii*icol+jj]
		}
	}

	// CMS & CVN
	if pm.cmsvn {
		if err = cp.featureNorm(gmm.FeatureData, icol, irow); err != nil {
			//log.Error(err)
			return fmt.Errorf("Feature Extract error -3")
		}
	}

	return nil
}
