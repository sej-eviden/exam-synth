package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// readDir reads the contents of the `exam` directory (2)
func readDir(dirName string) []string {
	dir, err := os.ReadDir(dirName)
	fmt.Println()
	examArray := make([]string, 0)
	if err != nil {
		fmt.Printf("Dir %s doesn't exist\n", dirName)
	}

	for _, d := range dir {
		if !d.IsDir() {
			continue
		} else {
			folderPath := fmt.Sprintf("%s/%s", dirName, d.Name())
			examArray = append(examArray, folderPath)
		}
	}
	// REMEMBER path.filepath.walk
	return examArray
}

// processDir reads the contents of each exam's dir and returns:
//
//	. raw HTML file
//	. Contents of Assets folder
func processDir(dirName string) ([]byte, error) {
	fmt.Printf("\n** %s **\n", dirName)

	dir, err := os.ReadDir(dirName)

	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if d.IsDir() {
			fmt.Printf("Dir inside dir -> %s\n", d.Name())
			// TODO copy paste images
		} else {
			fmt.Printf("File inside dir -> %s\n", d.Name())
			fileName := fmt.Sprintf("%s/%s", dirName, d.Name())
			file, err := os.ReadFile(fileName)

			if err != nil {
				log.Fatalf("ERROR: %s doesn't exist", fileName)
			} else {
				fmt.Printf("Reading exam questions from %s\n", fileName)
				// parseExam(file)
				return file, nil
			}
		}
	}
	return make([]byte, 0), errors.New("No html file found")
}

type question struct {
	Title   string   `json:"title"`
	Body    []string `json:"body"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type questions map[string]question

// parseExam reads through the contents of the exam file and creates a struct/json
func parseExam(file []byte) questions {
	fmt.Println("Parsing file")
	rawFile := string(file[:])
	lines := strings.Split(rawFile, "\n")
	var questionTitle, examTitle string
	var currentTitleI, questionI, currentBodyI, currentOptionI int

	var question question
	Questions := make(questions, 0)

	for i, line := range lines {
		if i > 1 && strings.Index(lines[i-1], "<h1") >= 0 {
			examTitle = strings.Trim(strings.Split(line, "Exam Actual Questions")[0], " ")
		}

		if strings.Index(line, "<div class=\"card-header text-white bg-primary\">") >= 0 {
			currentTitleI = i
			questionI += 1
		}

		if i == currentTitleI+1 && questionI > 0 {
			titleA := strings.Fields(strings.Trim(strings.Replace(line, "#", "", 1), " "))
			questionTitle = titleA[0] + fmt.Sprintf(" %03s", titleA[1])

			// fmt.Print(title)
		}

		if i == currentTitleI+3 && questionI > 0 {
			topicA := strings.Fields(strings.Trim(line, " "))
			topic := topicA[0] + fmt.Sprintf(" %03s", topicA[1])
			questionTitle = fmt.Sprintf("%s %s", topic, questionTitle)
			// fmt.Printf(" | %s\n", topic)

			question.Title = questionTitle
			// title = " ".join([word.rjust(3, "0") for word in title.split(" ") if word != "|"])
		}

		// GET Body
		if strings.Index(line, "<p class=\"card-text\">") >= 0 {
			currentBodyI = i
		}

		if i == currentBodyI+3 && question.Title != "" {
			paragraphs := strings.Split(strings.Trim(line, " "), "<br>")

			for _, p := range paragraphs {
				if strings.Index(p, "img") >= 0 {
					src := strings.Split(strings.Split(p, "\"")[1], "/")
					p = fmt.Sprintf("<img>/%s/img/%s<img>", examTitle, src[len(src)-1])
					fmt.Println(questionTitle)
				}
				question.Body = append(question.Body, p)
			}
		}

		if strings.Index(line, "<span class=\"inquestion-subtitle mb-0 mt-3\">Question</span>") >= 0 {
			question.Body = append(question.Body, "Question:")
			currentBodyI = i
		}

		// GET options
		if strings.Index(line, "<li class=\"multi-choice-item") >= 0 {
			currentOptionI = i
		}

		if i == currentOptionI+6 && question.Title != "" {
			question.Options = append(question.Options, strings.Trim(line, " "))
		}

		// GET answer
		if strings.Index(line, "<div class=\"vote-bar progress-bar bg-primary\"") >= 0 && question.Title != "" {
			firstSplit := strings.Split(line, "<div class=\"vote-bar progress-bar bg-primary\"")[1]
			ansStart := strings.Index(firstSplit, ">")
			secondSplit := firstSplit[ansStart:]
			ans := secondSplit[1:strings.Index(secondSplit, " ")]

			question.Answer = ans
			Questions[questionTitle] = question
			question.Answer = ""
			question.Title = ""
			question.Options = make([]string, 0)
			question.Body = make([]string, 0)
		}

		if strings.Index(line, "<span class=\"correct-answer\"><img") >= 0 {
			src := strings.Split(strings.Split(line, "/")[2], "\"")[0]

			p := fmt.Sprintf("<img>/%s/img/%s<img>", examTitle, src)
			fmt.Println(p)
			question.Options = append(question.Options, p)

			Questions[questionTitle] = question
			question.Answer = ""
			question.Title = ""
			question.Options = make([]string, 0)
			question.Body = make([]string, 0)
		}

		/*
			            if line.find("<span class=\"correct-answer\"><img") >= 0:
							src = line.split('/')[2].split('"')[0]
							questions[title]["options"].append(f"<img>/{file_title}/img/{src}<img>")
		*/

	}
	// fmt.Println("\n*****************")
	// fmt.Printf("File title: %s\nNumber of questions: %d\nLast title: %s\n", fileTitle, questionI, title)
	// fmt.Print("*****************\n\n")
	return Questions
}

func createJson(q questions, path string) {
	for _, v := range q {
		questionTitle := strings.ReplaceAll(v.Title, " ", "_")
		contents, err := json.Marshal(v)

		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(
			fmt.Sprintf("%s/%s.json", path, strings.ReplaceAll(questionTitle, " ", "_")),
			contents,
			os.ModePerm,
		)

		if err != nil {
			log.Fatal(err)
		}

	}
}

// makeDirs creates the necessary dirs:
//
//	1. Creates "outdir" directory
//  2. Creates a directory per exam (with the `img` directory included)
func makeDirs(outdir, examName string) (string, error) {

	fmt.Println("Making necessary dirs")
	examPath := fmt.Sprintf("./%s/%s", outdir, examName)
	path := fmt.Sprintf("%s/img", examPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return "", err
	}

	return examPath, nil
}

func main() {
	dir := flag.String("dir", "./exams", "Path to the folder where the \"raw\" exams are located")
	dest := flag.String("dest", "./results", "Path to the output folder")
	// src := flag.String("src", "./", "The folder or path where the aztfexport files are located")
	// cref := flag.String("conf", "./", "The folder or path where the yaml config file is located")
	// check := flag.Bool("validate", false, "Validate the contents of the yaml config against the terraform file")

	flag.Parse()
	fmt.Println(*dir)
	fmt.Println(*dest)

	exams := readDir(*dir)

	for _, exam := range exams {
		fmt.Printf("*** dirname %s ***\n", exam)
		examA := strings.Split(exam, "/")
		examPath, err := makeDirs(*dest, examA[len(examA)-1])

		rawFile, _ := processDir(exam)

		questions := parseExam(rawFile)
		createJson(questions, examPath)

		if err != nil {
			log.Fatal(err)
		}
	}

	// fmt.Println("wello")
}
