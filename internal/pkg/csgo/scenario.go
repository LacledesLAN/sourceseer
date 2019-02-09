package csgo

type Scenario func(*CSGO) *CSGO

func ClinchableMapCycle(maps []string) Scenario {
	if len(maps) == 0 {
		panic("At least one map must be provided")
	}

	if len(maps)%2 == 0 {
		panic("Must provide an odd number of maps")
	}

	if err := validateStockMapNames(maps); err != nil {
		panic(err)
	}

	return func(gs *CSGO) *CSGO {
		gs.srcds.AddLaunchArg("+map " + maps[0])

		return gs
	}
}

func CompetitiveWarmUp(mpTeamname1, mpTeamname2 string) Scenario {
	var args []string

	// Process mpTeamname1
	mpTeamname1 = SanitizeTeamName(mpTeamname1)
	if len(mpTeamname1) > 0 {
		args = append(args, "+mp_teamname_1", mpTeamname1)
	}

	// Process mpTeamname2
	mpTeamname2 = SanitizeTeamName(mpTeamname2)
	if len(mpTeamname2) > 0 {
		args = append(args, "+mp_teamname_2", mpTeamname2)
	}

	args = append(args, `+hostname "`+HostnameFromTeamNames(mpTeamname1, mpTeamname2)+`"`)

	return func(gs *CSGO) *CSGO {
		gs.srcds.AddLaunchArg(args...)

		return gs
	}
}
