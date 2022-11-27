package bfs

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BfsTestSuite struct {
	suite.Suite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(BfsTestSuite))
}

func (s *BfsTestSuite) TestRace() {
	obs := [][4]int{{1, 4, 2, 3}}
	rt := NewRaceTrack(5, 5, 4, 4, 4, 0, obs)
	s.Require().Equal(7, rt.Race())
}

func (s *BfsTestSuite) TestRace_NoSolution() {
	obs := [][4]int{{1, 1, 0, 2}, {0, 2, 1, 1}}
	rt := NewRaceTrack(3, 3, 2, 2, 0, 0, obs)
	s.Require().Equal(-1, rt.Race())
}

func (s *BfsTestSuite) TestFindSpeedDiff() {
	testState := carState{X: 0, Y: 0, Velocity: [2]int{0, 0}}
	got := testState.FindSpeedDiff()
	s.Require().Len(got, 9)
}

func (s *BfsTestSuite) TestFindStateIsValid() {
	testState1 := carState{X: 0, Y: 0, Velocity: [2]int{0, 0}}
	testState2 := carState{X: -1, Y: 0, Velocity: [2]int{0, 0}}
	testState3 := carState{X: 5, Y: 0, Velocity: [2]int{0, 0}}
	obs := [][4]int{{1, 4, 2, 3}}
	rt := NewRaceTrack(5, 5, 4, 4, 4, 0, obs)
	r := rt.(*raceTrack)
	s.Require().True(r.StateIsValid(testState1))
	s.Require().False(r.StateIsValid(testState2))
	s.Require().False(r.StateIsValid(testState3))
}

func (s *BfsTestSuite) TestNextScanAdd() {
	st := carState{X: 0, Y: 0, Velocity: [2]int{0, 0}}
	testStates := carStates{st}
	s.Require().True(testStates.Next())
	s.Require().EqualValues(*testStates.Scan(), st)
	s.Require().False(testStates.Next())
	s.Require().Len(testStates, 0)
	testStates.Add(st)
	s.Require().Len(testStates, 1)
	s.Require().True(testStates.Next())
	s.Require().EqualValues(*testStates.Scan(), st)
	s.Require().Len(testStates, 0)
}
