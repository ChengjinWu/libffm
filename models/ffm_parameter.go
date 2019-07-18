package models

type FFMParameter struct {
	// eta used for per-coordinate learning rate
	Eta float64
	// used for l2-regularization
	Lambda float64
	// max iterations
	NIters int
	// latent factor dim
	K int
	// instance-wise normalization
	Normalization bool
	// randomization training order of samples
	Random bool
}

func NewDefaultFFMParameter() *FFMParameter {
	return &FFMParameter{
		Eta:           0.1,
		Lambda:        0,
		NIters:        15,
		K:             4,
		Normalization: true,
		Random:        true,
	}
}
