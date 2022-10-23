package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(fileName string) ([]problem, error) {
	//read questions from csv
	//open file
	os.Open(fileName)
	if fObj, err := os.Open(fileName); err == nil {
		//create new reader
		csvR := csv.NewReader(fObj)
		//read file
		if cLines, err := csvR.ReadAll(); err == nil {
			//call the parseproblem function
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error in reading data in csv"+"format from %s file: %s", fileName, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", fileName, err.Error())
	}

}

func main() {
	//input the name of the file
	fName := flag.String("f", "quiz.csv", "path of csv file")
	//set duration of the timer
	timer := flag.Int("t", 30, "timer of the quiz")
	flag.Parse()
	//pull the problems from the file
	problems, err := problemPuller(*fName)
	//handle the error
	if err != nil {
		exit(fmt.Sprintf("Something went wrong:%s", err.Error()))
	}
	//create variable to count the correct answers
	correctAns := 0
	//initialize the timer
	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansChan := make(chan string)
	//loop through the questions, print on screen and accept the answers
problemLoop:

	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.q)

		go func() {
			fmt.Scanf("%s", &answer)
			ansChan <- answer
		}()
		select {
		case <-tObj.C: //timer ran out, break from loop
			fmt.Println()
			break problemLoop
		case iAns := <-ansChan: //if answer given matches ans in csv
			if iAns == p.a {
				correctAns++ //increament no of correct answers
			}
			if i == len(problems)-1 { //no of questions ran out
				close(ansChan)
			}
		}
	}
	//calculate and print out the result
	fmt.Printf("Your result is %d out of %d\n", correctAns, len(problems))
	fmt.Printf("Press enter to exit\n")
	<-ansChan
}

func parseProblem(lines [][]string) []problem {
	//go over questions and parse them
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{q: lines[i][0], a: lines[i][1]}
	}
	return r
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
