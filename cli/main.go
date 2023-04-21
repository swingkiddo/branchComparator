package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"io/ioutil"

	"github.com/swingkiddo/branchComparator"
)


func main() {
	b1Name := flag.String("b1", "", "Name of the branch for compare")
	b2Name := flag.String("b2", "", "Name of the branch to compare against")
	flag.Parse()

	b1 := branchComparator.NewBranch(*b1Name)
	b2 := branchComparator.NewBranch(*b2Name)
	result := branchComparator.CompareBranches(b1, b2)
	js, _ := json.MarshalIndent(result, "", "    ")
	ioutil.WriteFile("test.json", js, 0644)
	fmt.Println(string(js))
}