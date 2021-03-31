package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Question string
	Result string
}

func ProblemsReader(title string) map[int]Problem{
	file, err := os.Open(title + ".csv")
	problems := make(map[int]Problem)
	if err != nil{
		return nil
	}

	i := 0
	err = nil
	r := csv.NewReader(file)
	for {
		record, err := r.Read()

		if err == io.EOF{
			fmt.Println("Question read is over...")
			break
		}

		if err != nil{
			log.Fatal(err)
		}
		problems[i] = Problem{record[0], record[1]}
		i++
	}

	return problems
}

func doQuestion(problems map[int]Problem, answer *[]string, countdown time.Duration){
	c := make(chan bool)
	go wait(countdown, c)

	for i := 0; i < len(problems); i++{
		fmt.Printf("what %s equals to?\n", problems[i].Question)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter answer> ")
		tmp, _ := reader.ReadString('\n')

		select{
		case <-c:
			fmt.Println("\nTime's up, last answer will be unsubmitted.\nExiting...")
			return
		default:

			*answer = append(*answer, tmp)
		}
	}
}

func checkAnswers(answer []string, problems map[int]Problem) (int, int){
	correct := 0
	total := len(problems)

	for i := 0; i < len(answer); i++{
		if strings.Compare(problems[i].Result + "\n", answer[i]) == 0{
			correct++
		}
	}
	return correct, total
}

func wait(sec time.Duration, c chan bool){
	tick := time.Tick(sec * time.Second)
	for{
		select {
		case <-tick:
			c <- true
			return
		}
	}
}

func main(){
	var answers []string

	problems := ProblemsReader("problems")

	doQuestion(problems, &answers, 30)

	correct, total := checkAnswers(answers, problems)

	fmt.Printf("The number of correct answer is: %v\nThe total question was: %v", correct, total)
}