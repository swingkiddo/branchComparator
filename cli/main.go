/* 
	This is the CLI for using the branchComparator module. It parses arguments which are the branch names, and then gets packages of this branches, compares them and printing the result as a json string and also writes a result to a json file. 
	Packages are compared by their architectures. 
	The point of comparing is to get: 
		1) Packages of 1 branch which are not included in 2 branch 
		2) Packages of 2 branch which are not ibcluded in 1 branch
		3) Packages of 1 branch which are newer than packages in 2 branch
		
	Structure of resulting json:
		{
			"arch": {
				"branch1Differences": []Package,
				"branch2Differences": []Package,
				"branch1NewerPackages": []Package
			}...
		}
	
	The flags are:
		-b1
			The name of the branch you want to comapre
		-b2
			The name of the branch you want to compare against
*/


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