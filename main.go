package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// readDir reads the contents of the `exam` directory and calls processDir
func readDir(dirName string) (examArray []string) {
	dir, err := os.ReadDir(dirName)
	fmt.Println()
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
	return
}

// processDir reads the contents of each exam's dir and:
// 1. calls the folder creator
// 2. calls the parser
// 3. moves img's to img folder
func processDir(dirName string) {
	fmt.Printf("\n** %s **\n", dirName)
	// TODO Check if there is more than one html file (probably not)

	dir, err := os.ReadDir(dirName)

	if err != nil {
		fmt.Printf("Dir %s doesn't exist\n", dirName)
	}

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
				parseExam(file)
			}
		}
	}
}

// parseExam reads through the contents of the exam file and creates a struct/json
func parseExam(file []byte) {
	fmt.Println("Parsing file")

}

// makeDirs creates the necessary dirs for the outputs of parseExam
func makeDirs(outdir, examName string) error {

	fmt.Println("Making necessary dirs")
	path := fmt.Sprintf("./%s/%s/img", outdir, examName)
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func main() {

	dir := flag.String("dir", "./exams", "Path to the folder where the \"raw\" exams are located")
	dest := flag.String("dest", "./results", "Path to the output folder")
	// src := flag.String("src", "./", "The folder or path where the aztfexport files are located")
	// cref := flag.String("conf", "./", "The folder or path where the yaml config file is located")
	// check := flag.Bool("validate", false, "Validate the contents of the yaml config against the terraform file")
	//
	flag.Parse()
	fmt.Println(*dir)
	fmt.Println(*dest)

	exams := readDir(*dir)

	for _, exam := range exams {
		processDir(exam)

		fmt.Printf("*** dirname %s ***\n", exam)
		examA := strings.Split(exam, "/")
		err := makeDirs(*dest, examA[len(examA)-1])
		
		if err != nil {
			log.Fatal(err)
		}
	}

	// fmt.Println("wello")
}
