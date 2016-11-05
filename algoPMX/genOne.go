package main

import "fmt"
import "math/rand"
import "math"
import "time"
import "strings"

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz"

	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func main() {
  rand.Seed(time.Now().UTC().UnixNano())
  s := RandomString(50)
	fmt.Printf("base = %s\n", s)
	time.Sleep(2 * time.Second)
  population := generatePop(100)
  if evaluate(s, population) == true {
    return
  }
 	for key, value := range population{
    fmt.Printf("individu : %s  valeur: %f\n", key, value)
  }
  coreAlgo(s, population)
	fmt.Printf("base = %s\n", s)
}

type parents struct {
	father string
	mother string
}

func selectParents(population map[string]float64) parents{
  tmp := ""
  tmp2 := ""
	best := ""
	bestValue := 0.0
	for key,value := range population{
		if best == ""{
			best = key
			bestValue = value
		}
		if value >= population[best]{
			best = key
			bestValue = value
		}
	}
	randNum := math.Mod(rand.Float64(), bestValue)
  for key, value := range population{
    if value <= randNum && tmp != "" {
      if value > population[tmp] {
        tmp = key
      }
    } else {
			tmp = key
		}
  }
  for key, value := range population{
    if value <= bestValue && tmp2 != "" {
      if value > population[tmp2] && key != tmp {
        tmp2 = key
			}
			} else {
				tmp2 = key
			}
			if tmp2 == tmp{
				tmp2 = key
			}
    }
//		fmt.Printf("f = %s  m = %s\n", tmp, tmp2)
	//			time.Sleep(5 * time.Second)
	par := parents {father : tmp, mother : tmp2}
	return parents(par)
}

func searchInsert(index int, char byte, par parents, alleleOne int, alleleTwo int) int{
	for i := 0; i < 50; i++{
		if char == par.mother[i]{
			if i >= alleleOne && i <= alleleTwo{
				return int(searchInsert(i, par.father[i], par, alleleOne, alleleTwo))
			}else{
				return int(i)
			}
		}
	}
	return int(-1)
}

func loopCross(par parents, values map[byte]int, alleleOne int, alleleTwo int, child []byte)[]byte{
	isPMX := false
	isOk := false
	for i := 0; i < 50; i++{
		for j := 0; j < 50; j++{
			if par.father[i] == par.mother[j]{
				isOk = true
				break
			}
		}
		if isOk == false{
			isPMX = false
			break
		}
	}
	if isPMX == true{
	for key, value := range values{
		for ret := searchInsert(value, par.father[value], par, alleleOne, alleleTwo); ret != -1;{
			child[ret] = key
		}
	}
} else {
	tmp := rand.Int() % 1
	if tmp == 1{
		for i := 0; i < alleleOne; i++{
			child[i] = par.father[i]
		}
		for i := alleleTwo + 1; i < 50; i++{
			child[i] = par.father[i]
		}
	} else {
		for i := 0; i < alleleOne; i++{
			child[i] = par.mother[i]
		}
		for i := alleleTwo + 1; i < 50; i++{
			child[i] = par.mother[i]
		}
	}
}
return []byte(child)
}

func crossover(population map[string]float64, par parents) map[string]float64{
	popTmp := map[string]float64{}
	best := ""
	bestValue := 0.0
	for key,value := range population{
		if best == ""{
			best = key
			bestValue = value
		}
		if value >= population[best]{
			best = key
			bestValue = value
		}
	}
	if bestValue < 0.010 {
		cnt := 0
		isOk := true
		fmt.Printf("\nTrash best = %s %f\n", best, population[best])
		for x := 0; x < 100; x++{
			alleleOne := rand.Int() % (len(par.father) - 1)
			alleleTwo := 0
			for alleleTwo = rand.Int() % (len(par.father) - 1); alleleOne == alleleTwo; alleleTwo = rand.Int() % len(par.father){}
			if alleleOne > alleleTwo {
				tmp := alleleOne
				alleleOne = alleleTwo
				alleleTwo = tmp
			}
			child := []byte("                                                  ")
			values := map[byte]int{}
			for i := alleleOne; i <= alleleTwo; i++ {
				child[i] = byte(par.father[i])
			}
		//	size := alleleTwo - alleleOne
			for i := alleleOne; i <= alleleTwo; i++{
				for j := alleleOne; j < alleleTwo; j++{
					if child[i] == byte(par.mother[j]){
						values[child[i]] = j;
					}
				}
			}
			if isOk == true {
				loopCross(par, values, alleleOne, alleleTwo, child)
			} else{
				tmp := generatePop(1)
				for key := range tmp{
					child = []byte(key)
				}
			}
			if val, ok := popTmp[string(child)]; ok{
				x--
				cnt++
				if cnt > 1000{
					isOk = false
				}
				_ = val
			} else {
				cnt = 0
			}
		//	fmt.Printf("\n chiuld %s\n", child)
			popTmp[string(child)] = 0.0
		}
} else{
		cnt := 0
		isOk := true
		fmt.Printf("\nbest = %s %f\n", best, population[best])
		for x := 0; x < 99; x++{
			alleleOne := rand.Int() % (len(par.father) - 1)
			alleleTwo := 0
			for alleleTwo = rand.Int() % (len(par.father) - 1); alleleOne == alleleTwo; alleleTwo = rand.Int() % len(par.father){}
			if alleleOne > alleleTwo {
				tmp := alleleOne
				alleleOne = alleleTwo
				alleleTwo = tmp
			}
			child := []byte("                                                  ")
			values := map[byte]int{}
			for i := alleleOne; i <= alleleTwo; i++ {
				child[i] = byte(par.father[i])
			}
		//	size := alleleTwo - alleleOne
			for i := alleleOne; i <= alleleTwo; i++{
				for j := alleleOne; j < alleleTwo; j++{
					if child[i] == byte(par.mother[j]){
						values[child[i]] = j;
					}
				}
			}
			if isOk == true {
				loopCross(par, values, alleleOne, alleleTwo, child)
			} else{
				tmp := generatePop(1)
				for key := range tmp{
					child = []byte(key)
				}
			}
			if val, ok := popTmp[string(child)]; ok{
				x--
				cnt++
				if cnt > 1000{
					isOk = false
				}
				_ = val
			} else {
				cnt = 0
			}
		//	fmt.Printf("\n chiuld %s\n", child)
			popTmp[string(child)] = 0.0
		}
		popTmp[best] = 0.0
	}
	return map[string]float64(popTmp)
}

func coreAlgo(s string, population map[string]float64) bool{
	over := false
	i := 0
	for over == false{
		//time.Sleep(1 * time.Second)
	  par := selectParents(population)
	  population = crossover(population, par)
	  over = evaluate(s, population)
		i++
	}
		fmt.Printf("generation = %d\n", i)
		for key, value := range population{
			fmt.Printf("individu : %s  valeur: %f\n", key, value)
		}
  return bool(true)
}

func evaluate(end string, population map[string]float64) bool{
  size := len(end)
  total := 0.0
  for key := range population{
		cnt := 0
    result := 0.0
    for i := 0; i < size && i < len(key); i++ {
      if end[i] == key[i] {
				cnt++
        result += 1.0 / float64(size)
      } else if strings.ContainsAny(end, string(key[i])){
				result += 0.5 / float64(size)
			}
      if cnt == 50 {
				fmt.Printf("FOUND %s\n", key)
        return bool(true)
      }
    }
    population[key] = result
    total += result
  }
  for key, value := range population{
    population[key] = value / total
  }
  return bool(false)
}

func generatePop(size int) map[string]float64 {
  population := map[string]float64{}
  for i := 0; i < size; i++ {
    population[RandomString(50)] = 0;
  }
return map[string]float64(population)
}
