package genalgo

import (
	"fmt"
	"math"
	"math/rand"

	"../logfile"
	"../vgoapi"
	"sort"
)

var (
	NB_GENE = 9
	F_CROSSOVER_PT = NB_GENE / 4
	S_CROSSOVER_PT = NB_GENE / 2
)

const (
	NB_MOTOR       = 3
	MAX_GENE_DIFF  = 10


	NB_INDIVIDUAL  = 70
	NB_GENERATION_MIN  = 50
	GMUTATE_PC     = 30
	MUTATE_PC      = 10
	CROSSOVER_RATE = 15
	BEST_IND_NB    = NB_INDIVIDUAL / 4
	ADD_GENE = 0
)

var (
	nb_generation = 0
	bestScore float32 = 0.0
	totalScore float32 = 0.0
	generation = 0
	bestIndividu  *Individual = nil

	bestIndividuTotal *Individual = &Individual{ID: NB_INDIVIDUAL, Distance: 0, Fitness:0, Score:0}
	nb_gene_dif = 0
)

type Individual struct {
	ID        int
	Distance  float32
	Fitness   float32
	Gene      []float32
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
	ind.Gene = make([]float32, NB_GENE)
	for i := 0; i != NB_GENE; i++ {
		ind.Gene[i] = float32(rand.Intn(300))
	}
}

func init() {
	Population = make([]*Individual, NB_INDIVIDUAL)
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

		forwardVector := [3]float64{-1, 0, 0}
		dirVector := [3]float64{0, 0, 0}

		// on recupere l'individu courant
		indivual := Population[ind]

		leftPosStart, rightPosStart := vgoapi.GetWheelsStarPosition()
		endPos, endOrient := vgoapi.StartRobotMovement(indivual.Gene)

		// distance parcouru par chaque individu

		leftPosEnd, rightPosEnd := vgoapi.GetWheelsEndPosition()

		dirVector[0] = float64(leftPosEnd[0] - leftPosStart[0])
		dirVector[1] = float64(leftPosEnd[1] - leftPosStart[1])
		dirVector[2] = float64(leftPosEnd[2] - leftPosStart[2])

		lenghtDirVector := math.Sqrt((dirVector[0] * dirVector[0]) + (dirVector[1] * dirVector[1]) + (dirVector[2] * dirVector[2]))
		scalarProduct := ((dirVector[0] * forwardVector[0]) + (dirVector[1] * forwardVector[1]) + (dirVector[2] * forwardVector[2]))
		angle := scalarProduct / lenghtDirVector

	//	fmt.Printf("angle = %f\n", angle)

		dist := math.Sqrt(math.Pow(float64(endPos[0])*(180.0/math.Pi)-float64(startPos[0])*(180.0/math.Pi), 2) + math.Pow(float64(endPos[1])*(180.0/math.Pi)-float64(startPos[1])*(180.0/math.Pi), 2))
		//vgoapi.FinishSimulation()
		indivual.Distance = float32(dist)
		indivual.Score = float32(indivual.Distance)

		if leftPosEnd[2] > rightPosStart[2]+0.02 || leftPosEnd[2] > leftPosStart[2]+0.02 || rightPosEnd[2] > rightPosStart[2]+0.02 || rightPosEnd[2] > leftPosEnd[2]+0.02 {
			indivual.Score = indivual.Score / 100.0
		} else {
			indivual.Score = indivual.Score * float32(180*(math.Abs(float64(endPos[2]*(180.0/math.Pi)))+0.01)/100)
		}

		if (angle >= 0.97 || angle <= -0.97){
		//	fmt.Printf("Bon angle = %f\n\n", angle)
			indivual.Score += float32(math.Abs(angle * 100) / 2)
		} else if (angle <= 0.70 && angle >= -0.70){
			angle = (math.Acos(angle)) * (180.0 / math.Pi)
		//	fmt.Printf("Mauvais angle = %f\n\n", angle)
			indivual.Score -= float32(angle / 2)
		}

		if indivual.Score > bestScore {
			bestScore = indivual.Score
			bestIndividu = indivual
		}
		//fmt.Printf("Score ind : %0.5f\tBestScore ind : %0.5f\n", indivual.Score, bestIndividuTotal.Score)
		if  indivual.Score > bestIndividuTotal.Score {
		//	fmt.Printf("ON SET Score ind : %0.5f\tBestScore ind : %0.5f\n", indivual.Score, bestIndividuTotal.Score)

			//bestIndividuTotal.ID = indivual.ID
			bestIndividuTotal.Gene = indivual.Gene
			bestIndividuTotal.Score = indivual.Score
			nb_gene_dif = 0
		}
		// on set l'individu avec les resultats de la simulation
		indivual.ObjOrient = endOrient
		indivual.ObjPos = endPos
		totalScore += indivual.Score
	}
	if bestScore <= bestIndividuTotal.Score {
		nb_gene_dif += 1
	}
	fmt.Printf("TOTAL SCORE : %0.5f\n", bestIndividuTotal.Score)
}

// SELECTION des parents selon leur score
//
type Indivudus []*Individual

func (slice Indivudus) Len() int {
	return len(slice)
}

func (slice Indivudus) Less(i, j int) bool {
	return slice[i].Score > slice[j].Score;
}

func (slice Indivudus) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func contains(score float32, selection []*Individual) bool {
	for _, key := range selection {
		if key.Score == score {
			return true
		}
	}
	return false
}

func Selection() []*Individual {
	var selection []*Individual
	selection = append(selection, bestIndividu)
	for i := 1; i < BEST_IND_NB; i++ {
		//shuffle
		for i := range Population {
			j := rand.Intn(i + 1)
			Population[i], Population[j] = Population[j], Population[i]
		}
		randNum := float32(math.Mod(rand.Float64(), float64(totalScore)))
		var tmpScore float32 = 0.0
		for _, ind := range Population {
			tmpScore = tmpScore + ind.Score
			if tmpScore >= randNum && !contains(ind.Score, selection) {
				selection = append(selection, ind)
				totalScore -= ind.Score
				//Population = append(Population[:index], Population[index+1:]...)
				break
			}

		}

	}

	sort.Sort(Indivudus(Population))


/*	for i := 0; i < BEST_IND_NB; i++ {

		selection = append(selection, Population[i])
	}*/

//	selection = append(selection, bestIndividuTotal)
/*	for i := 0 ; i < BEST_IND_NB; i++ {
>>>>>>> Stashed changes
		var tmpScore float32 = 0.0

		// shuffle
		for i := range Population {
			j := rand.Intn(i + 1)
			Population[i], Population[j] = Population[j], Population[i]
		}

		randNum := float32(math.Mod(rand.Float64(), float64(totalScore)))

		for index, key := range Population {
			tmpScore = tmpScore + key.Score
<<<<<<< Updated upstream
			if tmpScore >= randNum && key.ID != bestIndividu.ID {
=======
			if (tmpScore >= randNum && key.Score != bestIndividuTotal.Score) {
>>>>>>> Stashed changes
				selection = append(selection, key)
				Population = append(Population[:index], Population[index+1:]...)
				totalScore -= key.Score
				break
			}

		}
	}*/

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

	if (nb_generation > 0 && nb_generation % 5 == 0) {
		fmt.Println("NEW GENE ADDED")
		NB_GENE += ADD_GENE
		F_CROSSOVER_PT = NB_GENE / 4
		S_CROSSOVER_PT = NB_GENE / 2
	}
	//Population = make([]*Individual, NB_INDIVIDUAL)
	x := 0
	x2 := 1
	sellen := len(selection)
	fmt.Printf("overcross rate  : %d\n", NB_INDIVIDUAL * CROSSOVER_RATE / 100)
	//Implementing two point crossover methods
	for i:= 0; i < NB_INDIVIDUAL * CROSSOVER_RATE / 100; i++ {
		fmt.Printf("on add ind score : %0.5f\n", selection[i].Score)
		Population[i] = Population[i]
		Population[i].ID = i
	}

	for idx := NB_INDIVIDUAL * CROSSOVER_RATE / 100; idx < NB_INDIVIDUAL - 1; idx += 2 {
		tmp := &Individual{ID: idx, Distance: 0.0, Fitness: 0.0, Gene:make([]float32, NB_GENE)}
		tmp2 := &Individual{ID: idx + 1, Distance: 0.0, Fitness: 0.0, Gene:make([]float32, NB_GENE)}
		x = rand.Intn(sellen)
		if x == sellen-1 {
			x2 = 0
		} else {
			x2 = x + 1
		}
		//fmt.Println("2 New Individual")
		//fmt.Println("x value :", x)
		for i := 0; i < F_CROSSOVER_PT; i++ {
			if i >= len(selection[x].Gene) {
				tmp.Gene[i] = selection[x].Gene[i % 3]
			} else {
				tmp.Gene[i] = selection[x].Gene[i]
			}
			if i >= len(selection[x2].Gene) {
				tmp2.Gene[i] = selection[x2].Gene[i % 3]
			} else {
				tmp2.Gene[i] = selection[x2].Gene[i]
			}
		}
		for i := F_CROSSOVER_PT; i < S_CROSSOVER_PT; i++ {
			if i >= len(selection[x2].Gene) {
				tmp.Gene[i] = selection[x2].Gene[i % 3]
			} else {
				tmp.Gene[i] = selection[x2].Gene[i]
			}
			if i >= len(selection[x].Gene) {
				tmp2.Gene[i] = selection[x].Gene[i % 3]
			} else {
				tmp2.Gene[i] = selection[x].Gene[i]
			}
		}
		for i := S_CROSSOVER_PT; i < NB_GENE; i++ {
			if i >= len(selection[x].Gene) {
				tmp.Gene[i] = selection[x].Gene[i % 3]
			} else {
				tmp.Gene[i] = selection[x].Gene[i]
			}
			if i >= len(selection[x2].Gene) {
				tmp2.Gene[i] = selection[x2].Gene[i % 3]
			} else {
				tmp2.Gene[i] = selection[x2].Gene[i]
			}
		}


		Population[idx] = tmp
		Population[idx+1] = tmp2
		PrintIndGen(tmp)
		PrintIndGen(tmp2)
	}


	for i := 0; i < int(NB_INDIVIDUAL*MUTATE_PC/100); i++ {
		ind := Population[i]
		for j := 0; j < int(len(ind.Gene)*GMUTATE_PC/100); j++ {
			ind.Gene[rand.Intn(len(ind.Gene)-1)] =  float32(rand.Intn(300))

		}
	}
	fmt.Printf("SIWE POPULATION : %d\n", len(Population))
	for _, ind := range Population {
		fmt.Printf("NEW INDIVIDU ID : %d\t", ind.ID)
		for i:= 0; i < len(ind.Gene); i++ {
			fmt.Printf("%d : %0.5f\t", i, ind.Gene[i])
		}
		fmt.Printf("\n")

	}

}

//PrintIndGen print all genes from one Individual
func PrintIndGen(ind *Individual) {
/*	fmt.Println("Genes for Ind :\t", ind.ID)
	for i := 0; i < NB_GENE/3; i++ {
		fmt.Println("[", ind.Gene[i*3], "]\t", "[", ind.Gene[i*3+1], "]\t", "[", ind.Gene[i*3+2], "]\t")
	}*/
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
	logfile.Write_data(generation, bestIndividuTotal.Score)
	generation += 1
}

func IsEnd() bool {
	nb_generation++
	fmt.Printf("ALORS nb gene dif : %d et nb_generation : %d\n", nb_gene_dif, nb_generation)
	if nb_gene_dif > MAX_GENE_DIFF  && nb_generation > NB_GENERATION_MIN {
		return true
	}
	return false
}

func MoveForward() {
	vgoapi.MoveWhile(bestIndividuTotal.Gene)
}
