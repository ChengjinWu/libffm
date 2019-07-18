package main

import (
	"encoding/json"
	"ffm/core"
	"ffm/models"
	"flag"
	"fmt"
)

var (
	vaPath string
	trPath string
	param  *models.FFMParameter = models.NewDefaultFFMParameter()
)

func init() {
	flag.Float64Var(&param.Eta, "eta", 0.1, "")
	flag.Float64Var(&param.Lambda, "Lambda", 0, "")
	flag.IntVar(&param.NIters, "NIters", 15, "")
	flag.IntVar(&param.K, "k", 4, "")
	flag.BoolVar(&param.Normalization, "Normalization", true, "")
	flag.BoolVar(&param.Random, "Random", true, "")
	flag.StringVar(&vaPath, "va", "", "")
	flag.StringVar(&trPath, "tr", "", "")
	flag.Parse()
}
func main() {
	tr, _ := core.ReadFFMProblem(trPath)
	va, _ := core.ReadFFMProblem(vaPath)
	model := core.FfmTrain(tr, va, param)
	jb, _ := json.Marshal(model)
	fmt.Println(string(jb))
}
