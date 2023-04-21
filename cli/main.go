package main

import (
	"flag"
	"fmt"
	"os"
	"encoding/json"
	"io/ioutil"
	"log"
	"github.com/swingkiddo/branchComparator"
)

var (
	logFileName = "info.log"
	jsonFileName = "result.json"
	branches = map[string]string {"p9":"", "p10":"", "sisyphus":""}
)

func main() {
	logFile, _ := os.Create(logFileName)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	log.SetOutput(logFile)
	
	b1Name := flag.String("b1", "", "Name of the branch for compare")
	b2Name := flag.String("b2", "", "Name of the branch to compare against")
	flag.Parse()

	if *b1Name == "" || *b2Name == "" {
		fmt.Println("You need to set the flags to run the program")
		os.Exit(1)
	}

	if _, ok := branches[*b1Name]; !ok {
		fmt.Printf("%s is unacceptable branch name. API provides only these names: p9, p10, sisyphus\n", *b1Name)
		os.Exit(1)
	}

	if _, ok := branches[*b2Name]; !ok {
		fmt.Printf("%s is unacceptable branch name. API provides only these names: p9, p10, sisyphus\n", *b2Name)
		os.Exit(1)
	}

	b1 := branchComparator.NewBranch(*b1Name)
	b2 := branchComparator.NewBranch(*b2Name)

	result := branchComparator.CompareBranches(b1, b2)

	js, json_err := json.MarshalIndent(result, "", "    ")
	if json_err != nil {
		log.Fatalln(json_err)
	} else {
		ioutil.WriteFile(jsonFileName, js, 0644)
		fmt.Println(string(js))
	}
}