package vpr

//https://github.com/liuxp0827/govpr.git
import (
	"io/ioutil"
	"log"
)

type engine struct {
	eng *vprEngine
}

func newEngine(sampleRate, delSilRange int, ubmFile, userModelFile string) (*engine, error) {
	vprEngine, err := newVPREngine(sampleRate, delSilRange, false, ubmFile, userModelFile)
	if err != nil {
		return nil, err
	}
	return &engine{eng: vprEngine}, nil
}

func loadWaveData(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	data = data[44:]
	return data, nil
}

func (obj *engine) trainSpeech(buffers [][]byte) error {
	var err error
	count := len(buffers)
	for i := 0; i < count; i++ {
		err = obj.eng.addTrainBuffer(buffers[i])
		if err != nil {
			return err
		}
	}

	defer obj.eng.clearTrainBuffer()
	defer obj.eng.clearAllBuffer()

	err = obj.eng.trainModel()
	if err != nil {
		return err
	}

	return nil
}

func (obj *engine) recSpeech(buffer []byte) (float64, error) {
	err := obj.eng.addVerifyBuffer(buffer)
	defer obj.eng.clearVerifyBuffer()
	if err != nil {
		return -1.0, err
	}

	err = obj.eng.verifyModel()
	if err != nil {
		return -1.0, err
	}

	return obj.eng.getScore(), nil
}

// WavTrain 训练模型
func WavTrain(wavPath string, ubmPF string, modePF string) bool {
	vprEngine, err := newEngine(16000, 50, ubmPF, modePF)
	if nil != err {
		return false
	}
	trainBuffer := make([][]byte, 0)
	buf, err := loadWaveData(wavPath)
	if nil != err {
		return false
	}
	trainBuffer = append(trainBuffer, buf)
	err = vprEngine.trainSpeech(trainBuffer)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

// WavVerification 验证说话人身份
func WavVerification(wavPath string, ubmPF string, modePF string) float64 {

	vprEngine, err := newEngine(16000, 50, ubmPF, modePF)
	if nil != err {
		return -1
	}

	//var threshold float64 = 1.0
	buf, err := waveLoad(wavPath)
	if nil != err {
		return -1
	}

	score, err := vprEngine.recSpeech(buf)
	if nil != err {
		return -1
	}

	return score
}
