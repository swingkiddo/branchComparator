package branchComparator

import (
	"fmt"
	"net/http"
	// "io"
	"encoding/json"
	"github.com/hashicorp/go-version"
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
	fmt.Printf("Initializing %s branch...\n", name)

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

		result[arch]["branch1Differences"] = b1Differences
		result[arch]["branch2Differences"] = b2Differences
		result[arch]["branch1NewerPackages"] = b1NewerPackages
	}
	return result
}


func getBranchPackages(branch_name string, branch interface{}) {
	fmt.Printf("Getting %s branch data \n", branch_name)
	url := EXPORT_BRANCH_URL + branch_name
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	json_err := json.NewDecoder(resp.Body).Decode(branch)
	if json_err != nil {
		fmt.Println(err)
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
			if pack.Release == p.Release {
				package1Version, _ := version.NewVersion(pack.Version)
				package2Version, _ := version.NewVersion(p.Version)
				if package1Version != nil && package2Version != nil {
					if package1Version.GreaterThan(package2Version) {
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


