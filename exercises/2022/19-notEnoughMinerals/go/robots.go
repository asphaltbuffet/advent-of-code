package exercises

func (s *state) farm() {
	s.ore += s.oreRobots
	s.clay += s.clayRobots
	s.obsidian += s.obsidianRobots
	s.geode += s.geodeRobots
}

func (s *state) makeGeodeRobot() {
	if s.ore < s.oreForGeodeRobot || s.obsidian < s.obsidianForGeodeRobot {
		panic("not enough minerals for geode robot")
	}
	s.ore -= s.oreForGeodeRobot
	s.obsidian -= s.obsidianForGeodeRobot

	s.geodeRobots++
}

func (s *state) makeOreRobot() {
	if s.ore < s.oreForOreRobot {
		panic("not enough minerals for ore robot")
	}

	s.ore -= s.oreForOreRobot

	s.oreRobots++
}

func (s *state) makeClayRobot() {
	if s.ore < s.oreForClayRobot {
		panic("not enough minerals for clay robot")
	}

	s.ore -= s.oreForClayRobot

	s.clayRobots++
}

func (s *state) makeObsidianRobot() {
	if s.ore < s.oreForObsidianRobot || s.clay < s.clayForObsidianRobot {
		panic("not enough minerals for obsidian robot")
	}

	s.ore -= s.oreForObsidianRobot
	s.clay -= s.clayForObsidianRobot

	s.obsidianRobots++
}
