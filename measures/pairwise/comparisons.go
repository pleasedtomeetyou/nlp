package pairwise

import (
	"math"

	"github.com/james-bowman/sparse"
	"gonum.org/v1/gonum/mat"
)

// Comparer is a type of function that compares two mat.Vector types and
// returns a value indicating how similar they are.
type Comparer func(a, b mat.Vector) float64

// CosineSimilarity calculates the cosine of the angles of 2 vectors i.e. how
// similar they are.  Possible values range up to 1 (exact match)
func CosineSimilarity(a, b mat.Vector) float64 {
	// Cosine angle between two vectors is equal to their dot product divided by
	// the product of their L2 norms
	dotProduct := sparse.Dot(a, b)
	norma := sparse.Norm(a, 2.0)
	normb := sparse.Norm(b, 2.0)

	return (dotProduct / (norma * normb))
}

// CosineDistance is the complement of CosineSimilarity in the positive space.
// 	CosineDistance = 1.0 - CosineSimilariy
// It should be noted that CosineDistance is not strictly a valid distance measure
// as it does not obey triangular inequality.  For applications requiring a distance
// measure that conforms with the strict definition then AngularDistance or
// Euclidean distance (with all vectors L2 normalised first) should be used instead.
// Whilst these distance measures may give different values, they will rank the same
// as CosineDistance.
func CosineDistance(a, b mat.Vector) float64 {
	return 1.0 - CosineSimilarity(a, b)
}

// AngularDistance is a distance measure closely related to CosineSimilarity.
// It measures the difference between the angles of 2 vectors by taking
// the inverse cosine (acos) of the CosineSimilarity and dividing by Pi.
// Unlike CosineSimilarity, this distance measure is a valid distance measure
// as it obeys triangular inequality.
// See https://en.wikipedia.org/wiki/Cosine_similarity#Angular_distance_and_similarity
func AngularDistance(a, b mat.Vector) float64 {
	cos := CosineSimilarity(a, b)
	if cos > 1 {
		cos = 1.0
	}
	theta := math.Acos(cos)
	return theta / math.Pi
}

// AngularSimilarity is the inverse of AngularDistance.
// 	AngularSimilarity = 1.0 - AngularDistance
func AngularSimilarity(a, b mat.Vector) float64 {
	return 1.0 - AngularDistance(a, b)
}

// HammingDistance is a distance measure sometimes referred to as the
// `Matching Distance` and measures how different the 2 vectors are
// in terms of the number of non-matching elements. This measurement
// is normalised to provide the distance as proportional to the total
// number of elements in the vectors.  If a and b are not the same
// shape then the function will panic.
func HammingDistance(a, b mat.Vector) float64 {
	ba, aok := a.(*sparse.BinaryVec)
	bb, bok := b.(*sparse.BinaryVec)

	if aok && bok {
		return float64(ba.DistanceFrom(bb)) / float64(ba.Len())
	}

	var count float64
	for i := 0; i < a.Len(); i++ {
		if a.AtVec(i) != b.AtVec(i) {
			count++
		}
	}
	return count / float64(a.Len())
}

// HammingSimilarity is the inverse of HammingDistance (1-HammingDistance)
// and represents the proportion of elements within the 2 vectors that
// exactly match.
func HammingSimilarity(a, b mat.Vector) float64 {
	return 1.0 - HammingDistance(a, b)
}
