package vpr

//https://github.com/liuxp0827/govpr.git
import (
	"fmt"
	"os"
	"path"
)

type vprEngine struct {
	trainBuf      []int16
	verifyBuf     []int16
	score         float64
	ubmFile       string
	userModelFile string
	deleteSil     bool
	delSilRange   int
	ubm           *gMM
	_minTrainLen  int64
	_minVerLen    int64
}

func newVPREngine(sampleRate, delSilRange int, deleteSil bool, ubmFile, userModelFile string) (*vprEngine, error) {
	engine := vprEngine{
		ubmFile:       ubmFile,
		userModelFile: userModelFile,
		verifyBuf:     make([]int16, 0),
		trainBuf:      make([]int16, 0),
		deleteSil:     deleteSil,
		delSilRange:   delSilRange,
		ubm:           newGMM(),
		_minTrainLen:  int64(sampleRate * 2),
		_minVerLen:    int64(float64(sampleRate) * 0.25),
	}

	err := engine.init()
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (obj *vprEngine) init() error {
	if obj.ubm == nil {
		return fmt.Errorf("ubm model is nil")
	}

	if err := obj.ubm.loadModel(obj.ubmFile); err != nil {
		return fmt.Errorf("model load failed")
	}
	return nil
}

func (obj *vprEngine) trainModel() error {
	if obj.trainBuf == nil || int64(len(obj.trainBuf)) < obj._minTrainLen {
		return fmt.Errorf("no available data")
	}

	tmpubm := newGMM()
	tmpubm.copy(obj.ubm)

	client := newGMM()
	client.dupModel(obj.ubm)
	if err := extract(obj.trainBuf, tmpubm); err != nil {
		return fmt.Errorf("memory insufficient")
	}

	for k := 0; k < MAXLOP; k++ {
		if ret, err := tmpubm.eM(tmpubm.Mixtures); ret == 0 || err != nil {
			return fmt.Errorf("train failed")
		}

		for i := 0; i < tmpubm.Mixtures; i++ {
			for j := 0; j < tmpubm.VectorSize; j++ {
				client.Mean[i][j] = (float64(tmpubm.Frames)*tmpubm.MixtureWeight[i])*
					tmpubm.Mean[i][j] + RELFACTOR*client.Mean[i][j]

				client.Mean[i][j] /= (float64(tmpubm.Frames)*tmpubm.MixtureWeight[i] + RELFACTOR)
			}
		}
	}

	userModelPath := path.Dir(obj.userModelFile)
	err := os.MkdirAll(userModelPath, 0755)
	if err != nil {
		return fmt.Errorf("train failed")
	}

	if err = client.SaveModel(obj.userModelFile); err != nil {
		return fmt.Errorf("train failed")
	}
	return nil
}

func (obj *vprEngine) verifyModel() error {
	if obj.verifyBuf == nil || len(obj.verifyBuf) <= 0 {
		return fmt.Errorf("no available data")
	}

	var buf []int16 = obj.verifyBuf
	var length int64

	//buf = DelSilence(obj.verifyBuf, obj.delSilRange)

	length = int64(len(buf))
	if length < obj._minVerLen {
		return fmt.Errorf("need more sample ")
	}

	var client *gMM = newGMM()
	err := client.loadModel(obj.userModelFile)
	if err != nil {
		return fmt.Errorf("model load failed")
	}

	tmpubm := newGMM()
	tmpubm.copy(obj.ubm)

	err = extract(buf, client)
	if err != nil {
		return fmt.Errorf("memory insufficient")
	}

	err = tmpubm.copyFeatureData(client)
	if err != nil {
		return fmt.Errorf("memory insufficient")
	}

	var logClient, logWorld float64
	logClient = client.lProb(client.FeatureData, 0, int64(client.Frames))
	logWorld = tmpubm.lProb(tmpubm.FeatureData, 0, int64(tmpubm.Frames))
	obj.score = (logClient - logWorld) / float64(client.Frames)
	return nil
}

func (obj *vprEngine) addTrainBuffer(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return fmt.Errorf("no available data")
	}

	sBuff := make([]int16, 0)
	length := len(buf)
	for ii := 0; ii < length-1; ii += 2 {
		cBuff16 := int16(buf[ii])
		cBuff16 |= int16(buf[ii+1]) << 8
		sBuff = append(sBuff, cBuff16)
	}

	if obj.deleteSil {
		sBuff = delSilence(sBuff, obj.delSilRange)
	}

	obj.trainBuf = append(obj.trainBuf, sBuff...)
	return nil
}

func (obj *vprEngine) addVerifyBuffer(buf []byte) error {
	if buf == nil || len(buf) == 0 {
		return fmt.Errorf("no available data")
	}

	sBuff := make([]int16, 0)
	length := len(buf)
	for ii := 0; ii < length; ii += 2 {
		cBuff16 := int16(buf[ii])
		cBuff16 |= int16(buf[ii+1]) << 8
		sBuff = append(sBuff, cBuff16)
	}

	if obj.deleteSil {
		sBuff = delSilence(sBuff, obj.delSilRange)
	}

	obj.verifyBuf = sBuff
	return nil
}

func (obj *vprEngine) clearTrainBuffer() {
	obj.trainBuf = obj.trainBuf[:0]
}

func (obj *vprEngine) clearVerifyBuffer() {
	obj.verifyBuf = obj.verifyBuf[:0]
}

func (obj *vprEngine) clearAllBuffer() {
	obj.clearTrainBuffer()
	obj.clearVerifyBuffer()
}

func (obj *vprEngine) getScore() float64 {
	return obj.score
}
