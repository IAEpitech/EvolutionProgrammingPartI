package genalgo

import (
	"fmt"
	"math"
	"math/rand"

	"../logfile"
	"../vgoapi"
	"sort"
	"time"
)

const (
	NB_GENE = 15
	MAX_GENE_DIFF = 5

	F_CROSSOVER_PT = NB_GENE / 3
	S_CROSSOVER_PT = 9

	NB_INDIVIDUAL = 150
	NB_GENERATION_MIN = 20
	//GMUTATE_PC     = 15
	MUTATE_PC = 1
	CROSSOVER_RATE = 10
	BEST_IND_NB = NB_INDIVIDUAL / 5
	LOGFILE = "basic"
)

var fileBest *logfile.File = logfile.New(LOGFILE + "_courbe_bestScore1")
var fileTotal *logfile.File = logfile.New(LOGFILE + "_courbe_totalScore1")
var fileMedian *logfile.File = logfile.New(LOGFILE + "_courbe_medianScore1")

var (
	nb_generation = 0
	bestScore float32 = 0.0
	totalScore float32 = 0.0
	generation = 0
	bestIndividu  *Individual = nil
	bestScoreTotal float32 = 0.0

	bestIndividuTotal *Individual = &Individual{ID: NB_INDIVIDUAL, Distance: 0, Score:0,
		Gene: []float32{258, 212, 196, 248, 103, 260, 174, 206, 240, 131, 282, 177, 202, 154, 283}}
	nb_gene_dif = 0
)

type Individual struct {
	ID       int
	Distance float32
	Gene     []float32
	Score    float32
}

func (ind *Individual) Copy() *Individual {
	ind2 := &Individual{Score:ind.Score, Distance: ind.Distance, ID:ind.ID, Gene:make([]float32, NB_GENE)}
	for i := 0; i < NB_GENE; i++ {
		ind2.Gene[i] = ind.Gene[i]
	}
	return ind2
}

var Population []*Individual

func createGene(ind *Individual) {
	ind.Gene = make([]float32, NB_GENE)
	for i := 0; i != NB_GENE; i++ {
		ind.Gene[i] = float32(rand.Intn(300))
	}
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	Population = make([]*Individual, NB_INDIVIDUAL)
	for x := 0; x < NB_INDIVIDUAL; x++ {
		tmp := &Individual{ID: x, Distance: 0.0}
		createGene(tmp)
		Population[x] = tmp
	}
}

func calculateScoreFromOrientation(wrist [3]float32, elb [3]float32, shd [3]float32) float32 {
	diff := (math.Abs(float64(shd[0] * (180.0 / math.Pi))) * 2 + math.Abs(math.Abs(float64(shd[1] * (180.0 / math.Pi))) - 90) + math.Abs(float64(shd[2] * (180.0 / math.Pi)))) * 5
	diff += ((math.Abs(math.Abs(float64(elb[0] * (180.0 / math.Pi))) - 180)) * 2 + math.Abs(math.Abs(float64(wrist[2] * (180.0 / math.Pi))) - 180)) * 5
	diff += math.Abs(math.Abs(float64(wrist[0] * (180.0 / math.Pi))) - 180) + math.Abs(float64(wrist[2] * (180.0 / math.Pi))) + math.Abs(float64(wrist[1] * (180.0 / math.Pi)))
	return float32(diff) / 100
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

		endPos, _ := vgoapi.StartRobotMovement(indivual.Gene)

		// distance parcouru par chaque individu

		leftPosEnd, rightPosEnd := vgoapi.GetWheelsEndPosition()
		dist := math.Sqrt(math.Pow(float64(endPos[0]) * (180.0 / math.Pi) - float64(startPos[0]) * (180.0 / math.Pi), 2) + math.Pow(float64(endPos[1]) * (180.0 / math.Pi) - float64(startPos[1]) * (180.0 / math.Pi), 2))
		indivual.Distance = float32(dist)
		indivual.Score = float32(indivual.Distance)

		wristOr, elbOr, shldOr := vgoapi.GetMotorsOrienation()

		oren := calculateScoreFromOrientation(wristOr, elbOr, shldOr)

		fmt.Printf("Orient : %0.5f\n", oren)
		if leftPosEnd[2] > leftPosStart[2] + 0.001 || leftPosEnd[2] < leftPosStart[2] - 0.001 ||
			rightPosEnd[2] > rightPosStart[2] + 0.001 || rightPosEnd[2] < rightPosStart[2] - 0.001 {
			indivual.Score = indivual.Score / 100.0
		} else {
			indivual.Score = indivual.Score * float32(180 * (math.Abs(float64(endPos[2] * (180.0 / math.Pi))) + 0.01) / 100)
		}
		if oren != 0 {
			indivual.Score = indivual.Score / float32(math.Pow(float64(oren), 2))
		}
		fmt.Printf("ID : %d\tScore : %0.5f\tDist : %0.5f\n", indivual.ID, indivual.Score, indivual.Distance)
		if indivual.Score > bestScore {
			bestScore = indivual.Score
			//bestIndividu = indivual
			bestIndividu = indivual.Copy()
			/*			bestIndividu = &Individual{Score:indivual.Score, Distance: indivual.Distance, ID:indivual.ID,
							Gene: make([]float32, NB_GENE)}
						for i := 0; i < 0; i++ {
							bestIndividu.Gene[i] = indivual.Gene[i]
						}*/
		}
		if indivual.Score > bestScoreTotal {
			nb_gene_dif = 0
			bestScoreTotal = indivual.Score
			bestIndividuTotal = indivual.Copy()

			//bestIndividuTotal = indivual
			/*			bestIndividuTotal = &Individual{Score:indivual.Score, Distance: indivual.Distance, ID:indivual.ID,
						Gene: make([]float32, NB_GENE)}
						for i := 0; i < 0; i++ {
							bestIndividuTotal.Gene[i] = indivual.Gene[i]
						}*/

		}
		totalScore += indivual.Score
	}
	if bestScore <= bestScoreTotal {
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
	selection = append(selection, bestIndividuTotal.Copy())
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
				selection = append(selection, ind.Copy())
				totalScore -= ind.Score
				//Population = append(Population[:index], Population[index+1:]...)
				break
			}

		}

	}

	sort.Sort(Indivudus(selection))

	for _, key := range selection {
		fmt.Printf("Parent choosen score : %0.5f\n", key.Score)
	}

	return selection
}

//GeneratePopulation : MERGE des parents pour creer une nouvelle population d'enfant avec system de mutation
func GeneratePopulation(selection []*Individual) {
	x := 0
	x2 := 1
	sellen := len(selection)
	fmt.Printf("overcross rate  : %d\n", NB_INDIVIDUAL * CROSSOVER_RATE / 100)
	//Implementing two point crossover methods
	for i := 0; i < NB_INDIVIDUAL * CROSSOVER_RATE / 100; i++ {
		fmt.Printf("on add ind score : %0.5f\n", selection[i].Score)
		Population[i] = selection[i]
		Population[i].ID = i
	}

	for idx := NB_INDIVIDUAL * CROSSOVER_RATE / 100; idx < NB_INDIVIDUAL; idx += 2 {
		/*		if idx == NB_INDIVIDUAL - 1 {
					break;
				}*/
		tmp := &Individual{ID: idx, Distance: 0.0, Gene:make([]float32, NB_GENE)}
		tmp2 := &Individual{ID: idx + 1, Distance: 0.0, Gene:make([]float32, NB_GENE)}
		x = rand.Intn(sellen)
		x2 = x
		for x2 == x {
			x2 = rand.Intn(sellen)
		}
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
		if idx != NB_INDIVIDUAL - 1 {
			Population[idx + 1] = tmp2
		}
		PrintIndGen(tmp)
		PrintIndGen(tmp2)
	}

	for index1, ind := range Population {
		if index1 != 0 {
			for index, _ := range ind.Gene {
				if rand.Intn(100) <= MUTATE_PC {

					if (rand.Intn(2) == 0) {
						ind.Gene[index] = float32(int((ind.Gene[index] + float32(rand.Intn(50)))) % 300)
					} else {
						ind.Gene[index] -= float32(rand.Intn(50))
						if ind.Gene[index] < 0 {
							ind.Gene[index] += 300
						}
					}
					fmt.Printf("On mutate index : %d\tind ID : %f\tScore : %0.5f\n", index1, ind.ID, ind.Score)
				}
			}
		}
	}
	Population[0], Population[NB_INDIVIDUAL - 1] = Population[NB_INDIVIDUAL - 1], Population[0]

	fmt.Printf("SIWE POPULATION : %d\n", len(Population))
	for _, ind := range Population {
		fmt.Printf("NEW INDIVIDU ID : %d\t", ind.ID)
		for i := 0; i < len(ind.Gene); i++ {
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
	sort.Sort(Indivudus(Population))
	fmt.Printf("BEST SCORE : %0.5f\n\n", bestScore)
	fileBest.Write_data(generation, bestIndividuTotal.Score)
	fileTotal.Write_data(generation, totalScore)
	fileMedian.Write_data(generation, Population[NB_INDIVIDUAL / 2].Score)
	generation += 1
}

func IsEnd() bool {
	nb_generation++
	fmt.Printf("ALORS nb gene dif : %d et nb_generation : %d\tTotal score : %0.5f\n", nb_gene_dif, nb_generation, bestIndividuTotal.Score)
	for i := 0; i < len(bestIndividuTotal.Gene); i++ {
		fmt.Printf("%d : %0.5f\t", i, bestIndividuTotal.Gene[i])
	}
	if nb_gene_dif >= MAX_GENE_DIFF  && nb_generation >= NB_GENERATION_MIN {
		fileBest.Close()
		fileMedian.Close()
		fileTotal.Close()
		return true
	}
	return false
}

func MoveForward() {
	vgoapi.MoveWhile(bestIndividuTotal.Gene)
}