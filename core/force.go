package core

import (
	"fmt"
	"github.com/golang/geo/r3"
)

// ForceReturn holds the index and force on a particle
type ForceReturn struct {
	Index int
	F r3.Vector
}

// PairwiseLennardJonesForce calculates the force vector on particle Ri due to Rj using the Lennard Jones potential.
func PairwiseLennardJonesForce(Ri, Rj r3.Vector, L float64) r3.Vector {
	if Ri == Rj {
		panic(fmt.Sprintf("%v and %v are equal, the pairwise force is infinite", Ri, Rj))
	}
	r := Displacement(Ri, Rj, L)
	R2 := r.Norm2()
	iR2 := 1.0 / R2
	iR8 := iR2 * iR2 * iR2 * iR2
	iR14 := iR8 * iR2 * iR2 * iR2
	f := 4 * (-12*iR14 + 6*iR8)
	return r.Mul(f)
}

// InternalForce calculates the total force vector on particle Ri due to the other particles in R due to a pairwise force.
func InternalForce(i int, R []r3.Vector, L float64, c chan ForceReturn) {
	F := r3.Vector{0, 0, 0}
	for j := range R {
		if i != j {
			F = F.Add(PairwiseLennardJonesForce(R[i], R[j], L))
		}
	}
	// return F
	c <- ForceReturn{i, F}
}
