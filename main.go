package main

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)


const (
	KEEP_GOING          = 1
	GARBAGE_PROBABILITY = 0.01
)
var (
	names                     [10]string
	operationTypes            [4]string
	runningAvgs               map[string]float64
	runningAvgsBeforeOverflow map[string]float64
	totalsBeforeOverflow      map[string]float64
	ErrOverflow                          = errors.New("integer overflow")
)

func init() {
	names = [10]string{"Carlos", "Juan", "Felipe", "Sonia", "Kung", "Keanu", "Anakin", "Dalí", "Gustavo", "David"}
	operationTypes = [4]string{"pago", "cobro", "descuento", "inversión"}
	runningAvgs = make(map[string]float64)
	runningAvgsBeforeOverflow = make(map[string]float64)
	totalsBeforeOverflow = make(map[string]float64)
}

func main() {
	keepGoing := 1
	var count uint64 = 0
	var beforeOverflowIndex uint64 = 0
	overflow := false

	for ;; count++ {
		rnd := getRandomUint()
		operationType := getRandomOperationType()

		if !overflow{
			adding, e := Add64(uint64(totalsBeforeOverflow[operationType]), rnd)
			if e!= nil {
				overflow = true
				beforeOverflowIndex = count - 1
				for _, opType := range operationTypes{
					runningAvgsBeforeOverflow[opType] = runningAvgs[opType]
				}
			}else{
				totalsBeforeOverflow[operationType] = float64(adding)
			}
		}else{
			if keepGoing == KEEP_GOING{
				for _, opType := range operationTypes{
					fmt.Printf("\n%f %f",runningAvgsBeforeOverflow[opType], float64(totalsBeforeOverflow[opType]) / float64(beforeOverflowIndex))
				}
				break
			}
			keepGoing ++
		}
		//runningAvgs[operationType] = (float64(count) * runningAvgs[operationType] + float64(rnd)) / float64(count+1)

		runningAvgs[operationType] += float64(float64(rnd) - runningAvgs[operationType]) / float64(count +1);

		printData(rnd)
	}

	fmt.Printf("\n\n Overflow index: %d.  Total index: %d\n\n", beforeOverflowIndex, count)

	for _, opType := range operationTypes{
		fmt.Printf("\nOperationType: %s - [running avg : %d]  [avg before overflow: %f]  [running avg before overflow: %d]", opType, uint64(runningAvgs[opType]), float64(totalsBeforeOverflow[opType])/float64(beforeOverflowIndex) , uint64(runningAvgsBeforeOverflow[opType]))
	}


}

func printData(rnd uint64){
	fmt.Printf("[user:%s] [type:%s] [ammount:%d]\n", getRandomUser(), getRandomOperationType(), rnd)
	if getRandomFloat() < GARBAGE_PROBABILITY {
		printGarbageData()
	}
}

func printGarbageData(){
	fmt.Printf("[user:%s] Debug info %d\n", getRandomUser(), getRandomUint())
	fmt.Printf("[user:%s] Stack %d\n", getRandomUser(), getRandomInt(123456))
}

func getRandomUser() string{
	return names[getRandomInt(10)]
}

func getRandomOperationType() string{
	return operationTypes[getRandomInt(4)]
}


func getRandomUint() uint64{
	return rand.Uint64()/500
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
	if right > 0 {
		if left > math.MaxUint64-right {
			return 0, ErrOverflow
		}
	} else {
		if left < math.MaxUint64-right {
			return 0, ErrOverflow
		}
	}
	return left + right, nil
}