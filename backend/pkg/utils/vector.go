package utils

import (
	"slices"
	"solution/internal/domain/dto"
	"strings"

	"github.com/pgvector/pgvector-go"
)

func VectorFromTags(genres []string) pgvector.Vector {
	res := make([]float32, 46, 46)
	for _, genre := range genres {
		idx := slices.IndexFunc(dto.VectorTag[:], func(f string) bool {
			return strings.EqualFold(strings.ToLower(f), strings.ToLower(genre))
		})
		if idx == -1 {
			continue
		}
		res[idx] = 1
	}
	return pgvector.NewVector(res)
}

// dest if user vector, src is film vector
func SumVectors(dest, src *pgvector.Vector) pgvector.Vector {
	res := make([]float32, 46, 46)
	ds, ss := dest.Slice(), src.Slice()
	for ind := range ds {
		sum1 := ds[ind] + ss[ind]
		res[ind] = sum1
	}
	return pgvector.NewVector(res)
}

// dest if user vector, src is film vector
func SubVectors(dest, src *pgvector.Vector) pgvector.Vector {
	res := make([]float32, 46, 46)
	ds, ss := dest.Slice(), src.Slice()
	for ind := range ds {
		res[ind] = ds[ind] - ss[ind]
	}
	return pgvector.NewVector(res)
}

// get avg vector from users vector list
func AvgVector(src []pgvector.Vector) pgvector.Vector {
	vecAmt := float32(len(src[0].Slice()))
	res := make([]float32, 46, 46)
	for ind := range 46 {
		var sum float32
		for _, v := range src {
			sum += v.Slice()[ind]
		}
		sum /= vecAmt
		res[ind] = sum
	}
	return pgvector.NewVector(res)
}

func MultiplyVector(v *pgvector.Vector, weigth float32) pgvector.Vector {
	res := make([]float32, 46, 46)
	for ind, val := range v.Slice() {
		res[ind] = val * weigth
	}
	return pgvector.NewVector(res)
}

// maps 0-10 rating to -3 - +3 scale of weigth
func AdjustRatingWeigth(weigth int32) float32 {
	return -3 + (6 * (float32(weigth) / 10))
}
