package branchComparator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isGreaterVersion(t *testing.T) {
	cases := []struct{
		name string 
		version1 string
		version2 string 
		expected bool
	}{
		{
			name: "version is less",
			version1: "3.7.11",
			version2: "3.9.16",
			expected: false,
		},
		{
			name: "version are equal",
			version1: "4.4.1",
			version2: "4.4.1",
			expected: false,
		},
		{
			name: "version is greater",
			version1: "2.0.3",
			version2: "1.4.1",
			expected: true,
		},
		{
			name: "versions as date",
			version1: "20220418",
			version2: "20210721",
			expected: true,
		},
		{
			name: "different amount of dotted numbers case 1",
			version1: "1.9",
			version2: "2.3.1",
			expected: false,
		},
		{
			name: "different amount of dotted numbers case 2",
			version1: "4.3.5",
			version2: "4.3",
			expected: true,
		},
		{
			name: "complex name case 1",
			version1: "4.6.3.0.16.git5ecb40bc",
			version2: "4.4.2",
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := isGreater(c.version1, c.version2)
			assert.Equal(t, c.expected, got)
		})
	}
}

func Test_isGreaterRelease(t *testing.T) {
	cases := []struct{
		name string
		release1 string
		release2 string
		expected bool
	}{
		{
			name: "release is less",
			release1: "alt4",
			release2: "alt5",
			expected: false,
		},
		{
			name: "releases are equal",
			release1: "alt4",
			release2: "alt4",
			expected: false,
		},
		{
			name: "release is greater",
			release1: "alt5",
			release2: "alt4",
			expected: true,
		},
		{
			name: "first release name with underscore",
			release1: "alt1_4",
			release2: "alt2",
			expected: false,
		},
		{
			name: "second release name with underscore",
			release1: "alt3",
			release2: "alt2_5",
			expected: true,
		},
		{
			name: "both release names with underscored",
			release1: "alt4_2",
			release2: "alt4_1",
			expected: true,
		},
		{
			name: "complex names case 1",
			release1: "alt1_0.1.D20170810git19dd1a7d",
			release2: "alt1_0.1.D20170810git19dd1a7d.6",
			expected: false,
		},
		{
			name: "complex names case 2",
			release1: "alt8_6jpp9",
			release2: "alt8_6jpp8.M90P.1",
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func( t *testing.T) {
			got := isGreater(c.release1, c.release2)
			assert.Equal(t, c.expected, got)
		})
	}
}
