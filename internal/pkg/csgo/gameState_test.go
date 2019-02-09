package csgo

import (
	"testing"
)

func Test_mapChanged(t *testing.T) {
	sut := CSGO{
		mpTeamname1: "Angleyne",
		mpTeamname2: "Nikolai"}

	sut.mapChanged("one")
	mapOneAddress := sut.currentMap

	if sut.currentMap == nil {
		t.Error("sut.currentMap should be initialized.")
	}
	if sut.currentMap.name != "one" {
		t.Errorf("Current map's name should be 'one' buy was %q.", sut.currentMap.name)
	}
	if sut.currentMap.started.IsZero() {
		t.Error("Current map's start time should have been set.")
	}
	if sut.currentMap.mpTeam1.name != sut.mpTeamname1 {
		t.Errorf("Current map's mpTeam1 name should be %q but was %q.", sut.mpTeamname1, sut.currentMap.mpTeam1.name)
	}
	if sut.currentMap.mpTeam2.name != sut.mpTeamname2 {
		t.Errorf("Current map's mpTeam2 name should be %q but was %q.", sut.mpTeamname2, sut.currentMap.mpTeam2.name)
	}

	sut.mapChanged("two")
	mapTwoAddress := sut.currentMap

	if &mapOneAddress == &mapTwoAddress {
		t.Errorf("The memory address for the first map (%X) should not match the address for the second map (%X).", &mapOneAddress, &mapTwoAddress)
	}

	if sut.currentMap.name != "two" {
		t.Errorf("Current map's name should be 'two' but was %q.", sut.currentMap.name)
	}

	if sut.maps[0].ended.IsZero() {
		t.Error("Previous map's end time should have been set.")
	}
}
