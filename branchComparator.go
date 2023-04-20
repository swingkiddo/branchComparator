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

func (b1 Branch) CompareWithAnotherBranch(b2 Branch) map[string][]Package {
	result := make(map[string][]Package)
	b1SortedPackages := b1.sortPackagesByArchs()
	b2SortedPackages := b2.sortPackagesByArchs()
	for arch, b1Packages := range b1SortedPackages {
		b2Packages, _ := b2SortedPackages[arch]
		diff, newer := comparePackages(b1Packages, b2Packages)
		result["difference"] = diff
		result["newerPackages"] = newer
	}
	return result
}


func getBranchPackages(branch_name string, branch interface{}) {
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

func comparePackages(a, b []Package) ([]Package, []Package) {
	mb := make(map[string]Package, len(b))
    for _, pack:= range b {
        mb[pack.Name] = pack
    }
    var diff []Package
	var newerPackages []Package
    for _, pack := range a {
        if p, found := mb[pack.Name]; !found {
            diff = append(diff, pack)
        } else {
			if pack.Release == p.Release {
				package1Version, _ := version.NewVersion(pack.Version)
				package2Version, _ := version.NewVersion(p.Version)
				if package1Version != nil && package2Version != nil {
					if package1Version.GreaterThan(package2Version) {
						newerPackages = append(newerPackages, pack)
					}
				}
			}

		}
    }
    return diff, newerPackages
}

func NewBranch(name string) Branch {
	branch := Branch {Name: name}
	getBranchPackages(name, &branch)
	return branch
}
