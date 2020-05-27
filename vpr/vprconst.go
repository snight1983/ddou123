package vpr

//https://github.com/liuxp0827/govpr.git

//Speaker Recognition ---------------
const (
	RELFACTOR             = 16 // RELFACTOR 自适应因子
	MAXLOP                = 1  // MAXLOP 自适应次数
	BITPERSAMPLE          = 16
	SAMPLERATE            = 16000 // 采样率
	MINFRAMES             = 300
	DB                    = -3.0      // 归一化分贝量
	LOGZERO               = (-1.0E10) /* ~log(0) */
	LSMALL                = (-0.5E10) /* log values < LSMALL are set to LOGZERO */
	PI                    = 3.14159265358979
	VARFLOOR              = 0.005 // Variance floor, make sure the variance to be large enough	(old:0.005)
	VARCEILING            = 10.0  // Variance ceiling
	MAXLOOP               = 10
	LOWCUTOFF             = 250
	HIGHCUTOFF            = 3800
	FILTERBANKSIZE        = 24
	FRAMELENGTH           = 20
	FRAMESHIFTt           = 10
	MFCCORDER             = 16
	BSTATIC               = true
	BDYNAMIC              = true
	BACCE                 = false
	CMSVN                 = true
	DBNORM                = true
	ZEROGLOBALMEAN        = true
	FEATWARP              = false
	DIFPOL                = false
	DPSCC                 = false
	PDASCC                = false
	ENERGYNORM            = false
	RASTA                 = false
	SILFLOOR              = 50
	ENERGYSCALE           = 19
	FEATUREWARPINGWINSIZE = 300
	RASTACOFF             = 0.94
	VOCBLOCKLEN           = 1600
)

// WAV ---------------
var (
	MINVOCENG int     = 500
	SHRTMAX   int     = 35535
	DLOG2PAI  float64 = 1.837877066
)
