package core

import (
	"ffm/models"
	"math"
)

func NewLogLossEvalutor(testDataSize int) *models.LogLossEvalutor {
	return &models.LogLossEvalutor{
		TestDataSize: testDataSize,
		Logloss:      make([]float64, testDataSize),
		Position:     0,
		TotalLoss:    0,
	}
}

func addLogLoss(lle *models.LogLossEvalutor, loss float64) {
	lle.TotalLoss = lle.TotalLoss + loss - lle.Logloss[lle.Position]
	lle.Logloss[lle.Position] = loss
	lle.Position += 1
	if lle.Position >= lle.TestDataSize {
		lle.Position = 0
		lle.EnoughData = true
	}
}

func getAverageLogLoss(lle *models.LogLossEvalutor) float64 {
	if lle.EnoughData {
		return lle.TotalLoss / float64(lle.TestDataSize)
	} else {
		return lle.TotalLoss / float64(lle.Position)
	}
}

func calLogLoss(prob float64, y float64) float64 {
	p := math.Max(math.Min(prob, 1-1e-15), 1e-15)
	if y == 1 {
		return -math.Log(p)
	} else {
		return -math.Log(1. - p)
	}
}
