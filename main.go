package main

func main() {
	algorithms := AvailableAlgortihms{}
	algorithms = algorithms.Init()
	testedProblems := []string{"op", "opcv", "optw", "tdop", "opfp"}

	ExperimentCompareProblemSolvingTime(testedProblems, algorithms.RGA, "benchmarks", 50)
}
