package genalgo

import (
	"fmt"
	"math"
	"math/rand"

	"../logfile"
	"../vgoapi"
)

var (
	BEST_IND_NB    = 3
	NB_GENE        = 9
	NB_MOTOR       = 3
	F_CROSSOVER_PT = NB_GENE / 4
	S_CROSSOVER_PT = NB_GENE / 2
	NB_INDIVIDUAL  = 80
	NB_GENERATION  = 10
	GMUTATE_PC     = 30
	MUTATE_PC      = 5

	i                        = 0
	bestScore    float32     = 0.0
	totalScore   float32     = 0.0
	generation               = 0
	bestIndividu *Individual = nil
)

type Individual struct {
	ID        int
	Distance  float32
	Fitness   float32
	Gene      [9]float32
	ObjOrient [3]float32
	ObjPos    [3]float32
	Score     float32
}

var Population []*Individual

func (ind *Individual) printInfos() {
	fmt.Printf("Individu ID : %d\n", ind.ID)
	fmt.Printf("Score : %0.5f\n", ind.Score)
}

func createGene(ind *Individual) {
	for i := 0; i != NB_GENE; i++ {
		ind.Gene[i] = float32(rand.Intn(300))
	}
}

func init() {
	Population = make([]*Individual, NB_INDIVIDUAL)
	i = 0
	for x := 0; x < NB_INDIVIDUAL; x++ {
		tmp := &Individual{ID: x, Distance: 0.0, Fitness: 0.0}
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
	bestIndividu = nil
	bestScore = 0.0
	totalScore = 0.0
	for ind := 0; ind < NB_INDIVIDUAL; ind++ {

		// position de depart du robot
		startPos := [3]float32{0, 0, 0}

		// on recupere l'individu courant
		indivual := Population[ind]

		leftPosStart, rightPosStart := vgoapi.GetWheelsStarPosition()

		endPos, endOrient := vgoapi.StartRobotMovement(indivual.Gene)

		// distance parcouru par chaque individu

		leftPosEnd, rightPosEnd := vgoapi.GetWheelsEndPosition()
		dist := math.Sqrt(math.Pow(float64(endPos[0])*(180.0/math.Pi)-float64(startPos[0])*(180.0/math.Pi), 2) + math.Pow(float64(endPos[1])*(180.0/math.Pi)-float64(startPos[1])*(180.0/math.Pi), 2))
		//vgoapi.FinishSimulation()
		indivual.Distance = float32(dist)
		indivual.Score = float32(indivual.Distance)

		if leftPosEnd[2] > rightPosStart[2]+0.02 || leftPosEnd[2] > leftPosStart[2]+0.02 || rightPosEnd[2] > rightPosStart[2]+0.02 || rightPosEnd[2] > leftPosEnd[2]+0.02 {
			indivual.Score = indivual.Score / 100.0
		} else {
			indivual.Score = indivual.Score * float32(180*(math.Abs(float64(endPos[2]*(180.0/math.Pi)))+0.01)/100)
		}
		if indivual.Score > bestScore {
			bestScore = indivual.Score
			bestIndividu = indivual
		}
		indivual.Fitness = indivual.Score / bestScore
		// on set l'individu avec les resultats de la simulation
		indivual.ObjOrient = endOrient
		indivual.ObjPos = endPos
		totalScore += indivual.Score
	}
}

// SELECTION des parents selon leur score
func Selection() []*Individual {
	var selection []*Individual
	selection = append(selection, bestIndividu)

	for i := 0; i < BEST_IND_NB; i++ {
		var tmpScore float32 = 0.0
		randNum := float32(math.Mod(rand.Float64(), float64(totalScore)))

		for index, key := range Population {
			tmpScore = tmpScore + key.Score
			if tmpScore >= randNum && key.ID != bestIndividu.ID {
				selection = append(selection, key)
				Population = append(Population[:index], Population[index+1:]...)
				totalScore -= key.Score
				break
			}

		}
	}

	/*
		randNum := float32(math.Mod(rand.Float64(), float64(bestScore)))
		var tmp *Individual = nil
		for _, key := range Population {
			if key.Score <= randNum && tmp != nil {
				if key.Score > tmp.Score {
					tmp = key
				}
			} else {
				tmp = key
			}
		}
		selection = append(selection, tmp)
	}*/
	for _, key := range selection {
		fmt.Printf("Parent choosen score : %0.5f\n", key.Score)
	}

	return selection
}

//GeneratePopulation : MERGE des parents pour creer une nouvelle population d'enfant avec system de mutation
func GeneratePopulation(selection []*Individual) {
	Population = make([]*Individual, NB_INDIVIDUAL)
	x := 0
	x2 := 1
	i = 0
	sellen := len(selection)
	fmt.Printf("sellen : %d\n", sellen)
	//Implementing two point crossover methods
	for idx := 0; idx < NB_INDIVIDUAL-1; idx += 2 {
		tmp := &Individual{ID: idx, Distance: 0.0, Fitness: 0.0}
		tmp2 := &Individual{ID: idx + 1, Distance: 0.0, Fitness: 0.0}
		x = rand.Intn(sellen)
		if x == sellen-1 {
			x2 = 0
		} else {
			x2 = x + 1
		}
		fmt.Println("2 New Individual")
		fmt.Println("x value :", x)
		for i := 0; i < F_CROSSOVER_PT; i++ {
			tmp.Gene[i] = selection[x].Gene[i]
			tmp2.Gene[i] = selection[x2].Gene[i]
		}
		for i := F_CROSSOVER_PT; i < S_CROSSOVER_PT; i++ {
			tmp.Gene[i] = selection[x2].Gene[i]
			tmp2.Gene[i] = selection[x].Gene[i]
		}
		for i := S_CROSSOVER_PT; i < NB_GENE; i++ {
			tmp.Gene[i] = selection[x].Gene[i]
			tmp2.Gene[i] = selection[x2].Gene[i]
		}

		Population[idx] = tmp
		Population[idx+1] = tmp2
		PrintIndGen(tmp)
		PrintIndGen(tmp2)
	}
	for i := 0; i < int(NB_INDIVIDUAL*MUTATE_PC/100); i++ {
		ind := Population[i]
		for j := 0; j < int(NB_GENE*GMUTATE_PC/100); j++ {
			ind.Gene[rand.Intn(NB_GENE-1)] = float32(rand.Intn(300))
		}
	}

}

//PrintIndGen print all genes from one Individual
func PrintIndGen(ind *Individual) {
	fmt.Println("Genes for Ind :\t", ind.ID)
	for i := 0; i < NB_GENE/3; i++ {
		fmt.Println("[", ind.Gene[i*3], "]\t", "[", ind.Gene[i*3+1], "]\t", "[", ind.Gene[i*3+2], "]\t")
	}
}

func PrintPopulation() {
	for _, key := range Population {
		fmt.Printf("STARTING simlation -- Robot ID : %d\n", key.ID)
		fmt.Printf("Score : %0.5f\tDist : %0.5f\n", key.Score, key.Distance)
		/*		fmt.Printf("value robotPos :  x = %0.5f\ty = %0.5f\tz = %0.5f\tObjOrient x : %0.5f\tObjOrient y : %0.5f\tObjOrient z : %0.5f\ndistance = %0.5f\n" +
				"Score : %d\n",
				key.ObjPos[0] * (180.0 / math.Pi), key.ObjPos[1]* (180.0 / math.Pi), key.ObjPos[2]* (180.0 / math.Pi),
				key.ObjOrient[0]* (180.0 / math.Pi), key.ObjOrient[1]* (180.0 / math.Pi), key.ObjOrient[2]* (180.0 / math.Pi), key.Distance, key.Score)*/
		fmt.Printf("FINISHING simlation -- Robot ID : %d\n\n", key.ID)

	}
	fmt.Printf("BEST SCORE : %0.5f\n\n", bestScore)
	logfile.Write_data(generation, bestScore)
	generation += 1
}
