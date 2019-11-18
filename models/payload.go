package models

type RootPayload struct {
	Data PayloadBody `json:"payloads"`
}

type Histogram struct {
	NumBuckets int64                  `json:"histogramNumBuckets"`
	Value      map[string]interface{} `json:"histogramValue"`
}

type ApplicationLaunchMetrics struct {
	HistogrammedResumeTime         Histogram `json:"histogrammedResumeTime"`
	HistogrammedTimeToFirstDrawKey Histogram `json:"histogrammedTimeToFirstDrawKey"`
}

type ApplicationResponsivenessMetrics struct {
	HistogrammedAppHangTime Histogram `json:"histogrammedAppHangTime"`
}

type ApplicationTimeMetrics struct {
	CumulativeBackgroundAudioTime    string `json:"cumulativeBackgroundAudioTime"`
	CumulativeBackgroundLocationTime string `json:"cumulativeBackgroundLocationTime"`
	CumulativeBackgroundTime         string `json:"cumulativeBackgroundTime"`
	CumulativeForegroundTime         string `json:"cumulativeForegroundTime"`
}

type CellularConditionMetrics struct {
	CellConditionTime Histogram `json:"cellConditionTime"`
}

type CPUMetrics struct {
	CumulativeCPUTime string `json:"cumulativeCPUTime"`
}

type DiskIOMetrics struct {
	CumulativeLogicalWrites string `json:"cumulativeLogicalWrites"`
}

type AveragePixelLuminance struct {
	AverageValue      string `json:"averageValue"`
	SampleCount       int64  `json:"sampleCount"`
	StandardDeviation int64  `json:"standardDeviation"`
}

type DisplayMetrics struct {
	AveragePixelLuminance AveragePixelLuminance `json:"averagePixelLuminance"`
}

type GpuMetrics struct {
	CumulativeGPUTime string `json:"cumulativeGPUTime"`
}

type LocationActivityMetrics struct {
	CumulativeBestAccuracyForNavigationTime string `json:"cumulativeBestAccuracyForNavigationTime"`
	CumulativeBestAccuracyTime              string `json:"cumulativeBestAccuracyTime"`
	CumulativeHundredMetersAccuracyTime     string `json:"cumulativeHundredMetersAccuracyTime"`
	CumulativeKilometerAccuracyTime         string `json:"cumulativeKilometerAccuracyTime"`
	CumulativeNearestTenMetersAccuracyTime  string `json:"cumulativeNearestTenMetersAccuracyTime"`
	CumulativeThreeKilometersAccuracyTime   string `json:"cumulativeThreeKilometersAccuracyTime"`
}

type AverageSuspendedMemory struct {
	AverageValue      string `json:"averageValue"`
	SampleCount       int64  `json:"sampleCount"`
	StandardDeviation int64  `json:"standardDeviation"`
}

type MemoryMetrics struct {
	AverageSuspendedMemory AverageSuspendedMemory `json:"averageSuspendedMemory"`
	PeakMemoryUsage        string                 `json:"peakMemoryUsage"`
}

type MetaData struct {
	AppBuildVersion string `json:"appBuildVersion"`
	DeviceType      string `json:"deviceType"`
	OsVersion       string `json:"osVersion"`
	RegionFormat    string `json:"regionFormat"`
}

type NetworkTransferMetrics struct {
	CumulativeCellularDownload string `json:"cumulativeCellularDownload"`
	CumulativeCellularUpload   string `json:"cumulativeCellularUpload"`
	CumulativeWifiDownload     string `json:"cumulativeWifiDownload"`
	CumulativeWifiUpload       string `json:"cumulativeWifiUpload"`
}

type SignpostIntervalData struct {
	HistogrammedSignpostDurations   Histogram `json:"histogrammedSignpostDurations"`
	SignpostAverageMemory           string    `json:"signpostAverageMemory"`
	SignpostCumulativeCPUTime       string    `json:"signpostCumulativeCPUTime"`
	SignpostCumulativeLogicalWrites string    `json:"signpostCumulativeLogicalWrites"`
}

type SignpostMetrics struct {
	SignpostCategory     string               `json:"signpostCategory"`
	SignpostIntervalData SignpostIntervalData `json:"signpostIntervalData"`
	SignpostName         string               `json:"signpostName"`
	TotalSignpostCount   int64                `json:"totalSignpostCount"`
}

type PayloadBody struct {
	ApplicationLaunchMetrics         ApplicationLaunchMetrics         `json:"applicationLaunchMetrics"`
	ApplicationResponsivenessMetrics ApplicationResponsivenessMetrics `json:"applicationResponsivenessMetrics"`
	ApplicationTimeMetrics           ApplicationTimeMetrics           `json:"applicationTimeMetrics"`
	AppVersion                       string                           `json:"appVersion"`
	CellularConditionMetrics         CellularConditionMetrics         `json:"cellularConditionMetrics"`
	CPUMetrics                       CPUMetrics                       `json:"cpuMetrics"`
	DiskIOMetrics                    DiskIOMetrics                    `json:"diskIOMetrics"`
	DisplayMetrics                   DisplayMetrics                   `json:"displayMetrics"`
	GpuMetrics                       GpuMetrics                       `json:"gpuMetrics"`
	LocationActivityMetrics          LocationActivityMetrics          `json:"locationActivityMetrics"`
	MemoryMetrics                    MemoryMetrics                    `json:"memoryMetrics"`
	MetaData                         MetaData                         `json:"metaData"`
	NetworkTransferMetrics           NetworkTransferMetrics           `json:"networkTransferMetrics"`
	SignpostMetrics                  []SignpostMetrics                `json:"signpostMetrics"`
	TimeStampBegin                   string                           `json:"timeStampBegin"`
	TimeStampEnd                     string                           `json:"timeStampEnd"`
}
