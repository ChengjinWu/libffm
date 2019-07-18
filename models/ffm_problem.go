package models

type FFMProblem struct {
	//data : field_num:feature_num:value
	// max(feature_num) + 1
	N int
	// max(field_num) + 1
	M int
	L int
	// X[ [P[0], P[1]) ], length=nnz
	X []*FFMNode
	// length=l+1
	Page []int
	// Y[0], length=l
	Y []float64
}
