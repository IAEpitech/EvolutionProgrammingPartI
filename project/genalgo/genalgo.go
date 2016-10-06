package genalgo

import (
	"math/rand"
	"fmt"
	"math"
)

var (
	NB_GENE = 9
	NB_MOTOR = 3
	NB_INDIVIDUAL = 5
	NB_GENERATION = 1
	i = 0
)

type Individual struct {
	Distance float32
	Fitness float32
	Gene []float32
	ObjOrient [3]float32
	ObjPos [3]float32
}

var popMaxScore float32
var maxIndScore float32

var Population []*Individual

func createGene(ind *Individual) {
	for i:= 0; i != NB_GENE; i++ {
		ind.Gene[i] = float32(rand.Intn(300)) * (math.Pi / 180)
	}
}


func init() {
	Population = make([]*Individual, NB_INDIVIDUAL)
	i = 0
	for x := 0; x < NB_INDIVIDUAL; x++ {
		tmp := &Individual{Distance: 0.0, Fitness: 0.0, Gene:make([]float32, NB_GENE)}
		createGene(tmp)
		Population[x] = tmp
	}
}

func Evaluate() {
/*	for _, key := range Population {
		if key.Distance > popMaxScore {
			popMaxScore = key.Distance
		}

		if key.Distance() > maxIndScore {
			maxIndScore = key.Distance
		}
	}*/
}
func PrintPopulation() {
	for _, key := range Population {
		fmt.Printf("value wrist : %0.5f\value elbow : %0.5f\tvalue elbow : %0.5f\tdistance : %0.5f\tObjOrient x : %0.5f\tObjOrient y : %0.5f\tObjOrient z : %0.5f\n",
		key.Gene[0], key.Gene[1], key.Gene[2], key.Distance, key.ObjOrient[0], key.ObjOrient[1], key.ObjOrient[2])
	}
}

func IsEnd() bool {
	if i == NB_GENERATION {
		return true
	}
	i++
	fmt.Printf("nb generation : %d\n", i)
	return false
}



