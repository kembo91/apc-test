package bfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// carState struct to track car position, velocity and num steps
type carState struct {
	X        int
	Y        int
	Velocity [2]int
	NumSteps int
}

// FindSpeedDiff calculates all possible changes of speed for a certain state
func (c *carState) FindSpeedDiff() [][2]int {
	diffX := []int{-1, 0, 1}
	diffY := []int{-1, 0, 1}
	rv := [][2]int{}
	for _, xdiff := range diffX {
		for _, ydiff := range diffY {
			newX := c.Velocity[0] + xdiff
			newY := c.Velocity[1] + ydiff
			if (newX < 4) && (newX > -4) && (newY < 4) && (newY > -4) {
				rv = append(rv, [2]int{newX, newY})
			}
		}
	}
	return rv
}

// carStates is a helper struct to ease work with states slice
type carStates []carState

// Next determines if there are elements in queue to process
func (c *carStates) Next() bool {
	if c != nil {
		if len(*c) > 0 {
			return true
		}
	}
	return false
}

// Scan returns state from queue and trims carStates slice from the front
func (c *carStates) Scan() *carState {
	rv := (*c)[0]
	*c = (*c)[1:]
	return &rv
}

// Add adds state to queue
func (c *carStates) Add(car carState) {
	*c = append(*c, car)
}

// RaceTrack is an exported interface so that clients have no way of working
// with a real struct
type RaceTrack interface {
	Race() int
}

// raceTrack is a struct for tracking raceTrack queue and state
type raceTrack struct {
	Queue  carStates
	StartX int
	StartY int
	Track  [][]int
}

// StateIsValid is a helper to determine if state position is valid
func (r *raceTrack) StateIsValid(st carState) bool {
	if st.X < len(r.Track) && st.Y < len(r.Track[0]) && st.X >= 0 && st.Y >= 0 {
		return true
	}
	return false
}

// Race starts out BFS over matrix, checks value on track,
// if we reach target, returns
// if position is invalid, continues
// for valid positions calculates possible changes of speed and adds those states with a hop
// in order to converge I assumed that max(numOptimalSteps) is less than number of squares on a grid X*Y
// that's probably not right, but that's the best I've got without digging into heuristic algorithms like ant colony optimization
func (r *raceTrack) Race() int {
	for r.Queue.Next() {
		cs := r.Queue.Scan()
		trackVal := r.Track[cs.X][cs.Y]
		switch trackVal {
		case 2:
			return cs.NumSteps
		case 1:
			continue
		}
		speeds := cs.FindSpeedDiff()
		for _, speed := range speeds {
			st := carState{
				Velocity: speed,
				X:        cs.X + speed[0],
				Y:        cs.Y + speed[1],
				NumSteps: cs.NumSteps + 1,
			}
			if r.StateIsValid(st) && st.NumSteps < len(r.Track)*len(r.Track[0]) {
				r.Queue.Add(st)
			}
		}
	}
	return -1
}

// NewRaceTracksFromFile is a file reader method that transforms data from a file to computable form
// it's a bit messy, but reading files is always messy
func NewRaceTracksFromFile(path string) ([]RaceTrack, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	numCases := 0
	if sc.Scan() {
		cases := sc.Text()
		numCases, err = strconv.Atoi(cases)
		if err != nil {
			return nil, fmt.Errorf("%s:%w", "wrong number of cases", err)
		}
	} else {
		return nil, fmt.Errorf("first line is empty")
	}
	rv := []RaceTrack{}
	for i := 0; i < numCases; i++ {
		obstacles := [][4]int{}
		var X, Y, startX, startY, targetX, targetY, numObstacles int
		for j := 0; j < 4; j++ {
			if sc.Scan() {
				line := sc.Text()
				switch j {
				case 0:
					coords := strings.Split(line, " ")
					if len(coords) != 2 {
						return nil, fmt.Errorf("invalid length in case %v line %v", i, j)
					}
					X, err = strconv.Atoi(coords[0])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid X", err)
					}
					Y, err = strconv.Atoi(coords[1])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid Y", err)
					}
				case 1:
					coords := strings.Split(line, " ")
					if len(coords) != 4 {
						return nil, fmt.Errorf("invalid length in case %v line %v", i, j)
					}
					startX, err = strconv.Atoi(coords[0])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid startX", err)
					}
					startY, err = strconv.Atoi(coords[1])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid startY", err)
					}
					targetX, err = strconv.Atoi(coords[2])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid targetX", err)
					}
					targetY, err = strconv.Atoi(coords[3])
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid targetY", err)
					}
				case 2:
					numObstacles, err = strconv.Atoi(line)
					if err != nil {
						return nil, fmt.Errorf("%s:%w", "invalid number of obstacles", err)
					}
					if numObstacles == 0 {
						j++
						continue
					}
				case 3:
					for k := 0; k < numObstacles; k++ {
						o, err := obstacleFromCoords(line)
						if err != nil {
							return nil, err
						}
						obstacles = append(obstacles, o)
					}
				}
			}
		}
		rv = append(rv, NewRaceTrack(X, Y, targetX, targetY, startX, startY, obstacles))
	}
	return rv, nil
}

func obstacleFromCoords(line string) ([4]int, error) {
	coords := strings.Split(line, " ")
	if len(coords) != 4 {
		return [4]int{0, 0, 0, 0}, fmt.Errorf("invalid number of obstacle coordinates")
	}
	x1, err := strconv.Atoi(coords[0])
	if err != nil {
		return [4]int{0, 0, 0, 0}, fmt.Errorf("invalid x1 obstacle coordinate")
	}
	x2, err := strconv.Atoi(coords[1])
	if err != nil {
		return [4]int{0, 0, 0, 0}, fmt.Errorf("invalid x2 obstacle coordinate")
	}
	y1, err := strconv.Atoi(coords[2])
	if err != nil {
		return [4]int{0, 0, 0, 0}, fmt.Errorf("invalid y1 obstacle coordinate")
	}
	y2, err := strconv.Atoi(coords[3])
	if err != nil {
		return [4]int{0, 0, 0, 0}, fmt.Errorf("invalid y2 obstacle coordinate")
	}
	return [4]int{x1, x2, y1, y2}, nil
}

// NewRaceTrack creates a computable RaceTrack instance
func NewRaceTrack(X, Y, targetX, targetY, startX, startY int, obstacles [][4]int) RaceTrack {
	track := make([][]int, X)
	for i := range track {
		track[i] = make([]int, Y)
	}
	for p := range obstacles {
		x1 := obstacles[p][0]
		x2 := obstacles[p][1]
		y1 := obstacles[p][2]
		y2 := obstacles[p][3]
		for i := x1; i <= x2; i++ {
			for j := y1; j <= y2; j++ {
				track[i][j] = 1
			}
		}
	}
	track[targetX][targetY] = 2
	return &raceTrack{
		Queue:  carStates{{X: startX, Y: startY, Velocity: [2]int{0, 0}, NumSteps: 0}},
		Track:  track,
		StartX: startX,
		StartY: startY,
	}
}
