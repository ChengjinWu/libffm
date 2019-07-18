package models

type LogLossEvalutor struct {
	TestDataSize int
	Logloss      []float64
	Position     int
	TotalLoss    float64
	EnoughData   bool
}
