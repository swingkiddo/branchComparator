# Branch Comparator
This repository contains a branchComparator module and a CLI for using it. \
Packages are compared by their architectures. \
The point of comparing is to get: \
	1) Packages of 1 branch which are not included in 2 branch \
	2) Packages of 2 branch which are not ibcluded in 1 branch \
	3) Packages of 1 branch which are newer than packages in 2 branch \

Structure of resulting json:
```
{
	"arch": {
		"branch1Differences": []Package,
		"branch2Differences": []Package,
		"branch1NewerPackages": []Package
	}...
}
```

## Clone the project
```
$ git clone https://github.com/swingkiddo/branchComparator.git
$ cd branchComparator
```

## Run the CLI
```
$ go run cli/main.go -b1 "branch1Name" -b2 "branch2Name"
```

The flags are required. Be sure you have set the names in quotes.
The API provides only these branches - p9, p10, sysiphus