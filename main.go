package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const set1Min int = 20
const set1Max int = 33
const set2Min int = 34
const set2Max int = 48

var elevators [12]Elevator

type Elevator struct {
	minFloor, maxFloor, currFloor int
	inTransit                     bool
	name                          string
}

func constructElevators() {
	rand.Seed(10)
	for i := 0; i < 6; i++ {
		elevators[i] = Elevator{set1Min, set1Max,
			set1Min + rand.Intn(13), false, string(65 + i)}
		fmt.Printf("%+v\n", elevators[i])
	}
	for i := 0; i < 6; i++ {
		elevators[i+6] = Elevator{set2Min, set2Max,
			set2Min + rand.Intn(15), false, string(71 + i)}
		fmt.Printf("%+v\n", elevators[i+6])
	}
}

func inRange(floor int) int {
	if floor == 1 {
		return 0
	}
	if floor >= 20 && floor <= 32 {
		return 1
	} else if floor >= 34 && floor <= 48 {
		return 2
	}
	return -1
}

func findElevator(startingFloor, desiredFloor int) (string, int) {
	inRangeStart := inRange(startingFloor)
	inRangeEnd := inRange(desiredFloor)
	if startingFloor == desiredFloor {
		return "Error, attempting to travel to the same floor", -1
	} else if inRangeStart == -1 || inRangeEnd == -1 {
		return fmt.Sprintf(`Error, entered floor(s) are inaccessible. 
							Floor ranges are %v to %v and %v to %v`,
			set1Min, set1Max, set2Min, set2Max), -2
	} else if inRange(startingFloor) != inRange(desiredFloor) &&
		startingFloor != 1 && desiredFloor != 1 {
		return fmt.Sprintf(`Error, floors are not in the same ranges. 
							Floor ranges are %v to %v and %v to %v`,
			set1Min, set1Max, set2Min, set2Max), -3
	} else {
		elevatorResult := -1
		for {
			closestFloor := 100
			if inRangeStart == 0 || inRangeStart == 1 {
				for i := 0; i < 6; i++ {
					if math.Abs(float64(elevators[i].currFloor-startingFloor)) < math.Abs(float64(closestFloor-startingFloor)) &&
						!elevators[i].inTransit {
						elevatorResult = i
						closestFloor = elevators[i].currFloor
					}
				}
			} else if inRangeStart == 2 {
				for i := 6; i < 12; i++ {
					if math.Abs(float64(elevators[i].currFloor-startingFloor)) < math.Abs(float64(closestFloor-startingFloor)) &&
						!elevators[i].inTransit {
						elevatorResult = i
						closestFloor = elevators[i].currFloor
					}
				}
			}
			if elevatorResult != -1 {
				break
			}
		}
		return elevators[elevatorResult].name, elevatorResult
	}
}

var wg, wg1 sync.WaitGroup

func moveElevator(elevatorIndex, desiredFloor int) {
	defer wg.Done()
	fmt.Println("Elevator", elevators[elevatorIndex].name, "starting on floor", elevators[elevatorIndex].currFloor)
	elevators[elevatorIndex].inTransit = true
	goingDown := elevators[elevatorIndex].currFloor > desiredFloor
	diff := math.Abs(float64(desiredFloor - elevators[elevatorIndex].currFloor))
	for i := 0.0; i < diff; i++ {
		if goingDown {
			elevators[elevatorIndex].currFloor--
			fmt.Println("Elevator", elevators[elevatorIndex].name, "moved down")
		} else {
			elevators[elevatorIndex].currFloor++
			fmt.Println("Elevator", elevators[elevatorIndex].name, "moved up")
		}
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println("Elevator", elevators[elevatorIndex].name, "arrived at floor", desiredFloor)
	elevators[elevatorIndex].inTransit = false
}

func callElevator(startingFloor, desiredFloor int) {
	defer wg1.Done()
	result, foundElevator := findElevator(startingFloor, desiredFloor)
	if foundElevator < 0 {
		fmt.Println(result)
		return
	}
	wg.Add(1)
	go moveElevator(foundElevator, startingFloor)
	wg.Wait()
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Boarded elevator")
	wg.Add(1)
	go moveElevator(foundElevator, desiredFloor)
	wg.Wait()
}

type elevatorHandler struct {
	store *[12]Elevator
}

func main() {
	constructElevators()
	mux := http.NewServeMux()
	elevatorH := &elevatorHandler{
		store: &elevators,
	}
	mux.Handle("/ping", elevatorH)
	mux.Handle("/allinfo", elevatorH)
	mux.Handle("/elevatorinfo/", elevatorH)
	mux.Handle("/callelevator/", elevatorH)
	mux.Handle("/dropelevator/", elevatorH)
	http.ListenAndServe("localhost:8080", mux)
	// fmt.Println()
	// wg1.Add(1)
	// go callElevator(40, 35)
	// wg1.Add(1)
	// go callElevator(22, 30)
	// wg1.Wait()
}

/**
POST:	/callelevator/?startingFloor={currfloor}&desiredFloor={desiredfloor}
		/dropelevator/?name={name}
GET:	/elevatorinfo?name={name}
		/allinfo
		/ping
UPDATE:	/update/?name={name}&newLower={newLower}&newHigher={newHigher}
*/

var getElevatorInfo = regexp.MustCompile(`^\/elevatorinfo\/\?name=([A-Z])$`)
var callElevatorStr = regexp.MustCompile(`^\/callelevator\/\?startingFloor=(\d+)\&desiredFloor=(\d+)$`)
var dropElevatorStr = regexp.MustCompile(`^\/dropelevator\/\?name=([A-Z])$`)

func (h *elevatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		if r.URL.Path == "/ping" {
			h.ping(w, r)
		} else if r.URL.Path == "/allinfo" {
			h.allInfo(w, r)
		} else if r.URL.Path == "/elevatorinfo/" {
			h.elevatorInfo(w, r)
		}
	case "POST":
		if r.URL.Path == "/callelevator/" {
			h.callElevatorAPI(w, r)
		} else if r.URL.Path == "/dropelevator/" {
			h.dropElevator(w, r)
		}
	}
}

func (h *elevatorHandler) dropElevator(w http.ResponseWriter, r *http.Request) {
	matches := dropElevatorStr.FindStringSubmatch(r.URL.String())
	if len(matches) < 2 {
		return
	}
	desiredName := matches[1]
	desiredIndex := int(desiredName[0]) - 65
	if elevators[desiredIndex].currFloor == elevators[desiredIndex].minFloor {
		w.Write([]byte(fmt.Sprintln("This elevator can't drop any further.")))
	} else {
		elevators[desiredIndex].currFloor--
		w.Write([]byte(fmt.Sprintln("Elevator", desiredName, "dropped one floor. Wheee!")))
	}
}

func (h *elevatorHandler) callElevatorAPI(w http.ResponseWriter, r *http.Request) {
	matches := callElevatorStr.FindStringSubmatch(r.URL.String())
	if len(matches) < 3 {
		return
	}
	startingFloor, err1 := strconv.Atoi(matches[1])
	desiredFloor, err2 := strconv.Atoi(matches[2])
	fmt.Println(startingFloor, desiredFloor)
	if err1 != nil || err2 != nil {
		fmt.Println("Err")
		return
	}
	w.Write([]byte(fmt.Sprintln("Elevator called to ", startingFloor)))
	wg1.Add(1)
	go callElevator(startingFloor, desiredFloor)
	wg1.Wait()
	w.Write([]byte(fmt.Sprintln("Elevator arrived at floor ", desiredFloor)))
}

func (h *elevatorHandler) ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`PONG!`))
}

func (h *elevatorHandler) allInfo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("%+v", h.store)))
}

func (h *elevatorHandler) elevatorInfo(w http.ResponseWriter, r *http.Request) {
	matches := getElevatorInfo.FindStringSubmatch(r.URL.String())
	if len(matches) < 2 {
		return
	}
	desiredName := matches[1]
	desiredIndex := int(desiredName[0]) - 65
	if 0 <= desiredIndex && desiredIndex <= 11 {
		w.Write([]byte(fmt.Sprintf("%+v", h.store[desiredIndex])))
	} else {
		w.Write([]byte(fmt.Sprint("Error: Elevator", desiredName, "could not be found. Please select an elevator from the letters A through L.")))
	}
}
