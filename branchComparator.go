package branchComparator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	// "github.com/hashicorp/go-version"
)

const (
	API_URL = "https://rdb.altlinux.org/api/"
	EXPORT_BRANCH_URL = API_URL + "/export/branch_binary_packages/"
)


type Branch struct {
	Name string `json:"name"`
	Length int `json:"length"`
	Packages []Package `json:"packages"`
}

func NewBranch(name string) Branch {
	fmt.Printf("Initializing %s branch started\n", name)

	branch := Branch {Name: name}
	getBranchPackages(name, &branch)
	
	fmt.Printf("%s branch initialized\n", name)
	return branch
}

type Package struct {
	Name string `json:"name"`
	Epoch int `json:"epoch"`
	Version string `json:"version"`
	Release string `json:"release"`
	Arch string `json:"arch"`
	DistTag string `json:"disttag"`
	Buildtime int `json:"buildtime"`
	Source string `json:"source"`
}

func (b Branch) sortPackagesByArchs() map[string][]Package {
	sorted := make(map[string][]Package)
	for _, p := range b.Packages {
		_, ok := sorted[p.Arch]
		if !ok {
			sorted[p.Arch] = []Package {p}
		} else {
			sorted[p.Arch] = append(sorted[p.Arch], p)
		}
	}
	return sorted
}


func CompareBranches(b1, b2 Branch) map[string]map[string][]Package {
	var result = map[string]map[string][]Package{}
	b1SortedPackages := b1.sortPackagesByArchs()
	b2SortedPackages := b2.sortPackagesByArchs()
	for arch, b1Packages := range b1SortedPackages {
		result[arch] = make(map[string][]Package)
		b2Packages, _ := b2SortedPackages[arch]
		b1Differences, b2Differences, b1NewerPackages := comparePackages(b1Packages, b2Packages)

		result[arch][fmt.Sprintf("%sDifferences", b1.Name)] = b1Differences
		result[arch][fmt.Sprintf("%sDifferences", b2.Name)] = b2Differences
		result[arch][fmt.Sprintf("%sNewerPackages", b1.Name)] = b1NewerPackages
	}
	return result
}

func getBranchPackages(branch_name string, branch interface{}) {
	fmt.Printf("Getting %s branch data... \n", branch_name)
	url := EXPORT_BRANCH_URL + branch_name
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	json_err := json.NewDecoder(resp.Body).Decode(branch)
	if json_err != nil {
		log.Fatalln(json_err)
	}
}

func isGreater(r1, r2 string) bool {
	re := regexp.MustCompile("[0-9]+")
	r1Versions := re.FindAllString(r1, -1)
	r2Versions := re.FindAllString(r2, -1)

	if len(r1Versions) <= len(r2Versions) {
		for i, v := range r1Versions {
			intV1, _ := strconv.ParseInt(v, 10, 64)
			intV2, _ := strconv.ParseInt(r2Versions[i], 10, 64)
			if intV1 > intV2 {
				return true
			} else if intV1 < intV2 {
				return false
			} 
		}
		return false
	} else {
		for i, v := range r2Versions {
			intV1, _ := strconv.ParseInt(r1Versions[i], 10, 64)
			intV2, _ := strconv.ParseInt(v, 10, 64)
			if intV1 > intV2 {
				return true
			} else if intV1 < intV2 {
				return false
			} 
		}
		return true
	}
}

func comparePackages(a, b []Package) ([]Package, []Package, []Package) {
	branch2MappedPackages := make(map[string]Package, len(b))
    for _, pack:= range b {
        branch2MappedPackages[pack.Name] = pack
    }
	var (
		branch1Differences []Package
		branch2Differences []Package
		branch1NewerPackages []Package

	)

    for _, pack := range a {
        if p, found := branch2MappedPackages[pack.Name]; !found {
            branch1Differences = append(branch1Differences, pack)
        } else {
			// if package versions are not equal, than we need to compare them
			// if they are equal, we need to compare 
			if pack.Version != p.Version {
				if isGreater(pack.Version, p.Version) {
					branch1NewerPackages = append(branch1NewerPackages, pack)
				} else if pack.Version == p.Version && pack.Release != p.Release {
				 	if isGreater(pack.Release, p.Release) {
						branch1NewerPackages = append(branch1NewerPackages, pack)
				 	}
				}
			}

			delete(branch2MappedPackages, pack.Name)
		}
    }

	for _, pack := range branch2MappedPackages {
		branch2Differences = append(branch2Differences, pack)
	}
    return branch1Differences, branch2Differences, branch1NewerPackages
}


