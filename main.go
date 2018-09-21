package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)


const (
	KEEP_GOING          = 3000
	GARBAGE_PROBABILITY = 0.01
)
var (
	names                     [10]string
	operationTypes            [4]string
	runningAvgs               map[string]float64
	runningAvgsBeforeOverflow map[string]float64
	totalsBeforeOverflow      map[string]uint64
	countsBeforeOverflow      map[string]uint64
	countsByOperation 		  map[string]uint64
	ErrOverflow                          = errors.New("integer overflow")
)

func init() {
	names = [...]string{"Carlos", "Juan", "Felipe", "Sonia", "Kung", "Keanu", "Anakin", "Dalí", "Gustavo", "David"}
	operationTypes = [...]string{"pago","cobro", "descuento", "inversión"}
	runningAvgs = make(map[string]float64)
	runningAvgsBeforeOverflow = make(map[string]float64)
	totalsBeforeOverflow = make(map[string]uint64)
	countsByOperation = make(map[string]uint64)
	countsBeforeOverflow = make(map[string]uint64)


}

func main() {
	keepGoing := 1
	var count uint64 = 0
	var overflowCount uint64 = 0
	overflow := false

	for ;; count++ {
		rnd := getRandomUint()
		user := getRandomUser()
		operationType := getRandomOperationType()

		if overflow {
			if keepGoing == KEEP_GOING {
				break
			}
			keepGoing ++
		}else{
			adding, e := Add64(totalsBeforeOverflow[operationType], rnd)
			if e!= nil || count > 101111  {
				overflowCount = count
				overflow = true
				for _, opType := range operationTypes{
					runningAvgsBeforeOverflow[opType] = runningAvgs[opType]
				}
			}else{
				totalsBeforeOverflow[operationType] = adding
				countsBeforeOverflow[operationType]++
			}
		}

		runningAvgs[operationType] = getRunningAvg(runningAvgs[operationType], countsByOperation[operationType], rnd)
		countsByOperation[operationType]++
		printData(rnd, user, operationType)
	}

	fmt.Printf("\n\n Overflow index: %d.  Total index: %d\n\n", overflowCount, count)

	for _, opType := range operationTypes{
		fmt.Printf("\nOperationType: %s - [running avg : %f]  [avg before overflow: %f]  [running avg before overflow: %f]", opType, runningAvgs[opType], float64(totalsBeforeOverflow[opType])/float64(countsBeforeOverflow[opType]) , runningAvgsBeforeOverflow[opType])
	}


}

func printData(rnd uint64, user, opType string){
	fmt.Printf("[user:%s] [type:%s] [ammount:%d]\n", user, opType, rnd)
	if getRandomFloat() < GARBAGE_PROBABILITY {
		printGarbageData()
	}
}

func printGarbageData(){
	fmt.Printf("[user:%s] Debug info %d\n", getRandomUser(), getRandomUint())
	fmt.Printf("[user:%s] Stack %d\n", getRandomUser(), getRandomInt(123456))
}


func getRunningAvg(oldAverage float64, averageCount, newNumber uint64) float64{
	return (float64(averageCount) * oldAverage + float64(newNumber)) / (float64(averageCount+1))
}
func getRandomUser() string{
	return names[getRandomInt(len(names))]
}

func getRandomOperationType() string{
	return operationTypes[getRandomInt(len(operationTypes))]
}


func getRandomUint() uint64{
	return rand.Uint64()/5000
}

func getRandomInt(n int) int{
	time.Sleep(189)
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(n)
}
func getRandomFloat() float64{
	return rand.Float64()
}

func Add64(left, right uint64) (uint64, error) {
	if left > math.MaxUint64-right {
		return 0, ErrOverflow
	}

	return left + right, nil
}