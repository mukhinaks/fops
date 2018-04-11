package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// CreateDirIfNotExist creates folders for output
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// ExperimentCompareProblemSolvingTime conducts experiments on computation time for 5 orienteering problems: OP, OPCV, OPTW, TDOP and OPFP.
func ExperimentCompareProblemSolvingTime(problems []string, algorithm string, outputFolderName string, numberOfProblemLaunches int) {
	fmt.Println("--------")
	fmt.Println(strings.ToUpper("Compare Problem Solving Time"))

	datasetSizes := []int{10, 50, 100, 500, 1000, 5000}

	for _, problem := range problems {
		fmt.Println("--------")
		fmt.Println(strings.ToUpper(problem))
		for _, datasetSize := range datasetSizes {
			fmt.Println("Dataset size:", datasetSize)
			switch algorithm {
			case "ACO":
				ProblemSolvingTime(problem, datasetSize, outputFolderName, numberOfProblemLaunches)

			case "RGA":
				ProblemSolvingTimeByRGA(problem, datasetSize, outputFolderName, numberOfProblemLaunches)

			default:
				fmt.Println("This algorithm is not implemented yet. Please, try one of those:")
				fmt.Println("RGA")
				fmt.Println("ACO")
			}
		}
		fmt.Println("--------")
	}
	fmt.Println("Done")
}

// ProblemSolvingTime runs specific problem of defined number of times; all resulting routes against with summary of each launch will be written in output folder.
func ProblemSolvingTime(problem string, datasetSize int, outputFolderName string, numberOfLaunches int) {
	algorithmFolder := "ACO"
	CreateDirIfNotExist(outputFolderName)
	CreateDirIfNotExist(filepath.Join(outputFolderName, algorithmFolder))
	CreateDirIfNotExist(filepath.Join(outputFolderName, algorithmFolder, problem))
	fileHandle, _ := os.Create(filepath.Join(outputFolderName, algorithmFolder, "experiment-"+problem+"-"+strconv.Itoa(datasetSize)+".txt"))
	writer := bufio.NewWriter(fileHandle)

	configPath := filepath.Join("experiments", "configs", "samples", "config-data-"+strconv.Itoa(datasetSize)+".json")

	switch problem {
	case "op":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveClassicalOP(configPath, 1, 3, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "tdop":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveTDOP(configPath, 1, 3, 600, 1000, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "optw":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPTW(configPath, 1, 3, 600, 1000, "0", filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "opcv":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPCV(configPath, []int{1, 0, 2, 3}, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	default:
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPFP(configPath, 1, 3, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	}
	writer.Flush()
}

// ProblemSolvingTime runs specific problem of defined number of times; all resulting routes against with summary of each launch will be written in output folder.
func ProblemSolvingTimeByRGA(problem string, datasetSize int, outputFolderName string, numberOfLaunches int) {
	algorithmFolder := "RGA"
	CreateDirIfNotExist(outputFolderName)
	CreateDirIfNotExist(filepath.Join(outputFolderName, algorithmFolder))
	CreateDirIfNotExist(filepath.Join(outputFolderName, algorithmFolder, problem))
	fileHandle, _ := os.Create(filepath.Join(outputFolderName, algorithmFolder, "experiment-"+problem+"-"+strconv.Itoa(datasetSize)+".txt"))
	writer := bufio.NewWriter(fileHandle)

	configPath := filepath.Join("experiments", "configs", "samples", "config-data-"+strconv.Itoa(datasetSize)+".json")

	switch problem {

	case "op":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveClassicalOPByRGA(configPath, 1, 3, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "tdop":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveTDOPByRGA(configPath, 1, 3, 600, 1000, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}
	case "optw":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPTWByRGA(configPath, 1, 3, 600, 1000, "0", filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}
	case "opcv":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPCVByRGA(configPath, []int{1, 0, 2, 3}, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}
	default:
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + strconv.Itoa(datasetSize) + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, algorithmFolder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPFPByNonameAlgorithm(configPath, 1, 3, 600, filePath) //SolveOPFPByNonameAlgorithm
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	}
	writer.Flush()
}

// Old stuff
func SomeLaunches() {

	/*
		ExperimentClassicalOPWithReference(solver, "", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "classic-op-with-ref")
		ExperimentOPCVMultipleDays(solver, "", []int{152, 3, 106, 105, 51, 63, 9, 127, 157, 158, 11, 13, 5191}, 600, 2, "opcv-multiple-days")
	*/
	//ExperimentCityBrand(solver, "", 720, 1000, 0, "Moscow")
}

func ExperimentOptimizeIterationAndAnts(problems []string, outputFolderName string, numberOfProblemLaunches int) {
	fmt.Println("--------")
	fmt.Println(strings.ToUpper("Experiment Optimize Iteration And Ants"))

	iterations := []int{10, 20, 50, 100, 200}
	ants := []string{"0.1", "0.25", "0.5", "1", "2"}

	for _, problem := range problems {
		fmt.Println("--------")
		fmt.Println(strings.ToUpper(problem))
		for _, iteration := range iterations {
			fmt.Println("Number of iterations:", iteration)
			for _, ant := range ants {
				fmt.Println("Number of ants:", ant)
				ProblemSolvingIterationsAndAntsTest(problem, iteration, ant, outputFolderName, numberOfProblemLaunches)
			}
		}
		fmt.Println("--------")
	}
	fmt.Println("Done")
}

func ProblemSolvingIterationsAndAntsTest(problem string, iterations int, ants string, outputFolderName string, numberOfLaunches int) {
	folder := "iterations"
	CreateDirIfNotExist(outputFolderName)
	CreateDirIfNotExist(filepath.Join(outputFolderName, folder))
	CreateDirIfNotExist(filepath.Join(outputFolderName, folder, problem))

	fileHandle, _ := os.Create(filepath.Join(outputFolderName, folder, "experiment-"+problem+"-iterations-"+strconv.Itoa(iterations)+"-ants-"+ants+".txt"))
	writer := bufio.NewWriter(fileHandle)

	configPath := filepath.Join("experiments", "configs", "iterations-ants", "config-iterations-"+strconv.Itoa(iterations)+"-ants-"+ants+".json")

	switch problem {
	case "op":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + problem + "-iterations-" + strconv.Itoa(iterations) + "-ants-" + ants + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, folder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveClassicalOP(configPath, 1, 3, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "tdop":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + problem + "-iterations-" + strconv.Itoa(iterations) + "-ants-" + ants + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, folder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveTDOP(configPath, 1, 3, 600, 1000, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "optw":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + problem + "-iterations-" + strconv.Itoa(iterations) + "-ants-" + ants + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, folder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPTW(configPath, 1, 3, 600, 1000, "0", filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	case "opcv":
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + problem + "-iterations-" + strconv.Itoa(iterations) + "-ants-" + ants + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, folder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPCV(configPath, []int{1, 0, 2, 3}, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	default:
		for i := 0; i < numberOfLaunches; i++ {
			fileName := "experiment-" + problem + "-iterations-" + strconv.Itoa(iterations) + "-ants-" + ants + "-" + strconv.Itoa(i) + ".json"
			filePath := filepath.Join(outputFolderName, folder, problem, fileName)
			t := time.Now()
			score, routeTime := SolveOPFP(configPath, 1, 3, 600, filePath)
			fmt.Fprintln(writer, score, routeTime, time.Since(t))
		}

	}
	writer.Flush()
}
