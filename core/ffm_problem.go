package core

import (
	"bufio"
	"ffm/models"
	"github.com/lunny/log"
	"os"
	"strconv"
	"strings"
)

func ReadFFMProblem(path string) (*models.FFMProblem, error) {
	problem := &models.FFMProblem{}
	l := 0
	nnz := 0
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("Cannot open text file:%s, err:%s", path, err)
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // or
		l += 1
		fields := strings.Split(line, "\t")
		nnz += len(fields) - 1
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("Cannot scanner text file: %s, err: %s", path, err)
		return nil, err
	}
	file.Close()
	log.Infof("reading %s, instance_num: %d, nnz: %d", path, l, nnz)
	problem.L = l
	problem.X = make([]*models.FFMNode, nnz)
	problem.Y = make([]float64, l)
	problem.Page = make([]int, l+1)
	problem.Page[0] = 0
	file, err = os.Open(path)
	if err != nil {
		log.Errorf("Cannot open text file:%s, err:%s", path, err)
		return nil, err
	}
	p := 0
	i := 0
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text() // or
		fields := strings.Split(line, "\t")
		problem.Y[i], _ = strconv.ParseFloat(fields[0], 64)
		for j := 1; j < len(fields); j++ {
			fields[j] = strings.Trim(fields[j], " ")
			subFields := strings.Split(fields[j], ":")
			node := models.FFMNode{}
			node.F, _ = strconv.Atoi(subFields[0])
			node.J, _ = strconv.Atoi(subFields[1])
			node.V, _ = strconv.ParseFloat(subFields[2], 64)
			problem.X[p] = &node
			p++
			problem.M = max(problem.M, node.F+1)
			problem.N = max(problem.N, node.J+1)
		}
		problem.Page[i+1] = p
		i++
	}
	if err := scanner.Err(); err != nil {
		log.Errorf("Cannot scanner text file: %s, err: %s", path, err)
		return nil, err
	}
	file.Close()
	return problem, nil
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
