package core

import (
	"ffm/models"
	"github.com/lunny/log"
	"math"
	"math/rand"
)

func InitFFMModel(n, m int, param *models.FFMParameter) *models.FFMModel {
	model := models.FFMModel{}
	model.N = n
	model.M = m
	model.K = param.K
	model.Normalization = param.Normalization
	model.W = make([]float64, model.N*model.M*model.K*2)
	coef := 0.5 / math.Sqrt(float64(model.K))

	position := 0
	for j := 0; j < model.N; j++ {
		for f := 0; f < model.M; f++ {
			for d := 0; d < model.K; d++ {
				model.W[position] = coef * rand.Float64()
				position += 1
			}
			for d := model.K; d < 2*model.K; d++ {
				model.W[position] = 1.0
				position += 1
			}
		}
	}
	return &model
}

func normalize(problem *models.FFMProblem, normal bool) []float64 {
	R := make([]float64, problem.L)
	if normal {
		for i := 0; i < problem.L; i++ {
			var norm float64 = 0
			for p := problem.Page[i]; p < problem.Page[i+1]; p++ {
				norm += problem.X[p].V * problem.X[p].V
			}
			R[i] = 1 / norm
		}
	} else {
		for i := 0; i < problem.L; i++ {
			R[i] = 1
		}
	}
	return R
}

func randomization(l int, randFlag bool) []int {
	order := make([]int, l)
	for i := 0; i < len(order); i++ {
		order[i] = i
	}
	if randFlag {
		for i := len(order); i > 1; i-- {
			tmp := order[i-1]
			index := rand.Intn(i)
			order[i-1] = order[index]
			order[index] = tmp
		}
	}
	return order
}

func wTx(prob *models.FFMProblem, i int, r float64, model *models.FFMModel, kappa float64, eta float64, lambda float64, doUpdate bool) float64 {
	start := prob.Page[i]
	end := prob.Page[i+1]
	var t float64
	align0 := model.K * 2
	align1 := model.M * model.K * 2
	for n1 := start; n1 < end; n1++ {
		j1 := prob.X[n1].J
		f1 := prob.X[n1].F
		v1 := prob.X[n1].V
		if j1 >= model.N || f1 >= model.M {
			continue
		}
		for n2 := n1 + 1; n2 < end; n2++ {
			j2 := prob.X[n2].J
			f2 := prob.X[n2].F
			v2 := prob.X[n2].V
			if j2 >= model.N || f2 >= model.M {
				continue
			}
			w1Index := j1*align1 + f2*align0
			w2Index := j2*align1 + f1*align0
			v := 2 * v1 * v2 * r
			if doUpdate {
				wg1Index := w1Index + model.K
				wg2Index := w2Index + model.K
				kappav := kappa * v
				for d := 0; d < model.K; d++ {
					g1 := lambda*model.W[w1Index+d] + kappav*model.W[w2Index+d]
					g2 := lambda*model.W[w2Index+d] + kappav*model.W[w1Index+d]

					wg1 := model.W[wg1Index+d] + g1*g1
					wg2 := model.W[wg2Index+d] + g2*g2

					model.W[w1Index+d] = model.W[w1Index+d] - eta/(math.Sqrt(wg1))*g1
					model.W[w2Index+d] = model.W[w2Index+d] - eta/(math.Sqrt(wg2))*g2

					model.W[wg1Index+d] = wg1
					model.W[wg2Index+d] = wg2
				}
			} else {
				for d := 0; d < model.K; d++ {
					t += model.W[w1Index+d] * model.W[w2Index+d] * v
				}
			}
		}
	}
	return t
}

func FfmTrain(tr, va *models.FFMProblem, param *models.FFMParameter) *models.FFMModel {
	model := InitFFMModel(tr.N, tr.M, param)
	rTr := normalize(tr, param.Normalization)
	var rVa []float64
	if va != nil {
		rVa = normalize(va, param.Normalization)
	}
	for iter := 0; iter < param.NIters; iter++ {
		var trLoss float64 = 0
		order := randomization(tr.L, param.Random)
		for ii := 0; ii < tr.L; ii++ {
			i := order[ii]
			y := tr.Y[i]
			r := rTr[i]
			t := wTx(tr, i, r, model, 0, 0, 0, false)
			expnyt := math.Exp(-y * t)
			trLoss += math.Log(1 + expnyt)
			kappa := -y * expnyt / (1 + expnyt)
			//log.Infof("i:%3d, y:%.1f, t:%.3f, expynt:%.3f, kappa:%.3f\n", i, y, t, expnyt, kappa)

			wTx(tr, i, r, model, kappa, param.Eta, param.Lambda, true)
		}
		trLoss /= float64(tr.L)
		log.Infof("iter: %2d, trLoss: %.5f", iter+1, trLoss)

		if va != nil && va.L != 0 {
			var vaLoss float64
			for i := 0; i < va.L; i++ {
				y := va.Y[i]
				r := rVa[i]
				t := wTx(va, i, r, model, 0, 0, 0, false)
				expnyt := math.Exp(-y * t)
				vaLoss += math.Log(1 + expnyt)
			}
			vaLoss /= float64(va.L)
			log.Infof(", vaLoss: %.5f", vaLoss)
		}
	}
	return model
}
