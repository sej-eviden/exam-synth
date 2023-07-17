package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

type ExamInfo struct {
	DirName string
	Total   int
}

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
func processDir(dirName string) ([]byte, string, []string, error) {
	fmt.Printf("\n%s\n", dirName)

	dir, err := os.ReadDir(dirName)

	var htmlFile []byte
	assetDirName := ""
	var imgs []string

	if err != nil {
		log.Fatal(err)
	}

	for _, d := range dir {
		if d.Name() == "imagenes" {
			continue
		}
		if d.IsDir() {
			// fmt.Printf("!!! Dir inside dir -> %s !!!\n", d.Name())
			assetDirName = fmt.Sprintf("%s/%s", dirName, d.Name())
			assetDir, err := os.ReadDir(assetDirName)

			if err != nil {
				fmt.Println("ERROR")
				log.Fatal(err)
			}
			for _, img := range assetDir {
				if strings.Index(img.Name(), ".png") >= 1 || strings.Index(img.Name(), ".jpg") >= 1 {
					// fmt.Println(img.Name())
					imgs = append(imgs, img.Name())
				}
			}
			// TODO copy paste images
		} else {
			// fmt.Printf("File inside dir -> %s\n", d.Name())
			fileName := fmt.Sprintf("%s/%s", dirName, d.Name())
			htmlFile, err = os.ReadFile(fileName)

			if err != nil {
				log.Fatalf("ERROR: %s doesn't exist", fileName)
			}
		}
	}
	if len(htmlFile) <= 0 {
		return make([]byte, 0), assetDirName, make([]string, 0), errors.New("no html file found")
	}
	return htmlFile, assetDirName, imgs, nil
}

type question struct {
	Title   string   `json:"title"`
	Body    []string `json:"body"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

type questionMap map[string]question

// parseExam reads through the contents of the exam file and creates a struct/json
func parseExam(file []byte) (questionMap, string) {
	// TODO extract repeated code into mini functions
	// fmt.Println("Parsing file")
	rawFile := string(file[:])
	lines := strings.Split(rawFile, "\n")
	var questionTitle, examTitle string
	var currentTitleI, questionI, currentBodyI, currentOptionI int

	var question question
	Questions := make(questionMap, 0)

	for i, line := range lines {
		if i > 1 && strings.Contains(lines[i-1], "<h1") {
			examTitle = strings.Trim(strings.Split(line, "Exam Actual Questions")[0], " ")
			// fmt.Println("Title", examTitle)
		}

		if strings.Contains(line, "<div class=\"card-header text-white bg-primary\">") {
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
		if strings.Contains(line, "<p class=\"card-text\">") {
			currentBodyI = i
		}

		if i == currentBodyI+3 && question.Title != "" {
			paragraphs := strings.Split(strings.Trim(line, " "), "<br>")

			for _, p := range paragraphs {
				if strings.Contains(p, "img") {
					src := strings.Split(strings.Split(p, "\"")[1], "/")
					p = fmt.Sprintf("<img>/%s/img/%s<img>", examTitle, src[len(src)-1])
					// fmt.Println(questionTitle)
				}
				question.Body = append(question.Body, p)
			}
		}

		if strings.Contains(line, "<span class=\"inquestion-subtitle mb-0 mt-3\">Question</span>") {
			question.Body = append(question.Body, "Question:")
			currentBodyI = i
		}

		// GET options
		if strings.Contains(line, "<li class=\"multi-choice-item") {
			currentOptionI = i
		}

		if i == currentOptionI+6 && question.Title != "" {
			question.Options = append(question.Options, strings.Trim(line, " "))
		}

		// GET answer
		if strings.Contains(line, "<div class=\"vote-bar progress-bar bg-primary\"") && question.Title != "" {
			firstSplit := strings.Split(line, "<div class=\"vote-bar progress-bar bg-primary\"")[1]
			ansStart := strings.Index(firstSplit, ">")
			secondSplit := firstSplit[ansStart:]
			ans := secondSplit[1:strings.Index(secondSplit, " ")]

			question.Answer = ans
		}

		// GET img answer
		if strings.Contains(line, "<span class=\"correct-answer\"><img") {
			src := strings.Split(strings.Split(line, "/")[2], "\"")[0]

			p := fmt.Sprintf("<img>/%s/img/%s<img>", examTitle, src)
			// fmt.Println(p)
			question.Options = append(question.Options, p)

		}

		if strings.Contains(line, "<!-- / Question  -->") {
			Questions[questionTitle] = question
			question.Answer = ""
			question.Title = ""
			question.Options = make([]string, 0)
			question.Body = make([]string, 0)
		}

	}
	return Questions, examTitle
}

func createJson(q questionMap, path string) {
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
//  1. Creates "outdir" directory
//  2. Creates a directory per exam (with the `img` directory included)
func makeDirs(outdir, examName string) (string, error) {

	// fmt.Println("Making necessary dirs")
	examPath := fmt.Sprintf("./%s/%s", outdir, examName)
	path := fmt.Sprintf("%s/img", examPath)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return "", err
	}

	return examPath, nil
}

func copyImg(srcPath, destPath string) error {
	inputFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't open the source file: %s", err)
	}

	stats, statErr := inputFile.Stat()
	if statErr != nil {
		return statErr
	}

	if stats.Size() <= 0 {
		return errors.New("empty file")
	}
	outputFile, err := os.Create(destPath)

	if err != nil {
		inputFile.Close()
		return fmt.Errorf("couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("writing to output file failed: %s", err)
	}

	// Removes the original file. Don't really want this, but good to know
	// err = os.Remove(srcPath)
	// if err != nil {
	// 	return fmt.Errorf("Failed removing original file: %s", err)
	// }
	return nil
}

func main() {
	// TODO create subcommand to show help
	// "<bin> info/i" with argv[0] I guess
	dir := flag.String("dir", "./exams", "Path to the folder where the \"raw\" exams are located")
	dest := flag.String("dest", "./results", "Path to the output folder")

	flag.Parse()
	fmt.Println("Exam files ->\t", *dir)
	fmt.Println("Destination folder ->\t", *dest)

	exams := readDir(*dir)
	masterInfoArr := []ExamInfo{}

	for _, exam := range exams {
		// fmt.Printf("*** dirname %s ***\n", exam)
		rawFile, assetDirName, imgs, _ := processDir(exam)

		questions, examTitle := parseExam(rawFile)
		examName := strings.ReplaceAll(strings.Split(examTitle, " - ")[0], " ", "_")
		fmt.Printf("Dest exam folder ->\t%s/%s\n", *dest, examName)

		examPath, err := makeDirs(*dest, examName)
		masterInfoArr = append(masterInfoArr, ExamInfo{DirName: examName, Total: len(questions)})

		if err != nil {
			log.Fatal(err)
		}

		createJson(questions, examPath)

		for _, img := range imgs {
			srcPath := fmt.Sprintf("%s/%s", assetDirName, img)
			destPath := fmt.Sprintf("%s/%s/img/%s", *dest, examName, img)
			e := copyImg(srcPath, destPath)
			if e != nil {
				log.Fatal(e)
			}
		}

		fmt.Println("All files successfully created.")
	}
	tmplFile := "master.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)

	if err != nil {
		fmt.Println("This is wrong")
	}

	masterFile, err := os.Create(*dest + "/master.txt")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(masterInfoArr)
	err = tmpl.Execute(masterFile, masterInfoArr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	fmt.Println()
	fmt.Print("\nTo continue with the process of adding new exams,\nplease refer to the steps detailed on the README.md file in the exam repo.\n")
	fmt.Print("\nAll the necessary changes/additions inside the files of the exam source code\nare located in the results/master.txt file.\n")
	fmt.Print("\n* Move the img folder inside each exam folder and paste them inside the\n  public directory, in the corresponding exams folder")
	fmt.Print("\n* Copy the json files corresponding to each exam and paste them\n  inside the proper folder in the src/content folder")
	fmt.Print("\n* Modify the source code with the lines specified in master.txt")
	fmt.Print("\n  NOTE: only copy the line of the arrays or objects where the exam information is located")
	fmt.Print("\n* Commit with a message that starts with \"Update:\" and push the modified code.")
	fmt.Println()
	fmt.Println()
	// err = tmpl.Execute(os.Stdout, masterInfoArr)

	// if err != nil {
	// 	fmt.Println(err)
	// }

}
