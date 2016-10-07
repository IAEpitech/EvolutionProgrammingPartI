package genalgo

import (
	"math/rand"
	"fmt"
	"math"
	"../vgoapi"
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

//
// ALGO pour set la distance, la position et l'orientation final du robot pour chaque individu
// et set le score de chaque individu
// PENSER a check l'orientation du robot (= sa stabilite) pour set son score
//
func Evaluate() {
	for ind := 0; ind < NB_INDIVIDUAL; ind++ {

		// position de depart du robot
		startPos := [3]float32{ 0,0,0}

		// distance parcouru par chaque individu
		dist := 0.0

		// on recupere l'individu courant
		indivual := Population[ind]
		newPos := make([]float32, NB_MOTOR)

		// position et orientation du robot apres chaque simulation
		var endPos [3]float32
		var endOrient [3]float32

		// on repete (3 fois pour l instant) un movement afin d'avoir plus d'infos sur l'individu
		for i := 0; i < NB_GENE; i += NB_MOTOR {
			// on set la nouvelle position pour chaque motor (wrist, elbow, shoulder)
			newPos[i / 3] = indivual.Gene[i]
			newPos[i / 3] =  indivual.Gene[i + 1]
			newPos[i / 3] = indivual.Gene[i + 2]

			vgoapi.StartSimulation()

			// start la simulation
			endPos, endOrient = vgoapi.StartRobotMovement(newPos)

			vgoapi.FinishSimulation()

			// calcul de la distance
			dist += math.Sqrt(math.Pow(float64(endPos[0]) * (180.0 / math.Pi) - float64(startPos[0]) * (180.0 / math.Pi), 2)+ math.Pow(float64(endPos[1])* (180.0 / math.Pi) - float64(startPos[1]) * (180.0 / math.Pi), 2))
		}
		// on set l'individu avec les resultats de la simulation
		indivual.Distance = float32(dist)
		indivual.ObjOrient = endOrient
		indivual.ObjPos = endPos
	}
}

//
// SELECTION des parents selon le score attribue -- Regarder algorithm de Roulette Wheel
//
func SelectParent() (*Individual, *Individual) {
	// pour l'instant on return deux randoms individus
	return Population[rand.Intn(NB_INDIVIDUAL)], Population[rand.Intn(NB_INDIVIDUAL)]
}

//
// MERGE des parents pour creer un enfant avec system de mutation
//
func CreateChild(parent1 *Individual, parent2 *Individual) *Individual {
	// pour l'instant on return un random individu
	return Population[rand.Intn(NB_INDIVIDUAL)]
}

//
// GENERATION de la nouvelle population
//
func GenerateNewPopulation() {

}

func PrintPopulation() {
	for _, key := range Population {
		fmt.Printf("value robotPos :  x = %0.5f\ty = %0.5f\tz = %0.5f\tObjOrient x : %0.5f\tObjOrient y : %0.5f\tObjOrient z : %0.5f\ndistance = %0.5f\n\n",
			key.ObjPos[0], key.ObjPos[1], key.ObjPos[2], key.ObjOrient[0], key.ObjOrient[1], key.ObjOrient[2], key.Distance)
	}
}



