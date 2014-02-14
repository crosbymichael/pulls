package gordon

import (
	"code.google.com/p/go.codereview/patch"
	"io"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
)

// ReviewPatch reads a git-formatted patch from `src`, and for each file affected by the patch
// it assign its Maintainers based on the current repository tree directories
// The list of Maintainers are generated when the MaintainerManager object is instantiated.
//
// The result is a map where the keys are the paths of files affected by the patch,
// and the values are the maintainers assigned to review that partiular file.
//
// There is no duplicate checks: the same maintainer may be present in multiple entries
// of the map, or even multiple times in the same entry if the MAINTAINERS file has
// duplicate lines.
func ReviewPatch(src io.Reader, maintainersDirMap *map[string][]*Maintainer) (reviewers map[string][]*Maintainer, err error) {
	reviewers = make(map[string][]*Maintainer)
	input, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	set, err := patch.Parse(input)
	if err != nil {
		return nil, err
	}
	for _, f := range set.File {
		for _, target := range []string{f.Dst, f.Src} {
			if target == "" {
				continue
			}
			target = path.Clean(target)
			if _, exists := reviewers[target]; exists {
				continue
			}
			targetDir := "."
			items := strings.Split(target, "/")
			for i := 0; i < len(items)-1; i++ {
				if targetDir == "." {
					targetDir = items[i]
				} else {
					targetDir = path.Join(targetDir, items[i])
				}
			}
			maintainers := (*maintainersDirMap)[targetDir]
			reviewers[target] = maintainers
		}
	}
	return reviewers, nil
}

type MaintainerFile map[string][]*Maintainer

type Maintainer struct {
	Username string
	FullName string
	Email    string
	Target   string
	Active   bool
	Lead     bool
	Raw      string
}

func parseMaintainer(line string) *Maintainer {
	const (
		commentIndex  = 1
		targetIndex   = 3
		fullnameIndex = 4
		emailIndex    = 5
		usernameIndex = 7
	)
	re := regexp.MustCompile("^[ \t]*(#|)((?P<target>[^: ]*) *:|) *(?P<fullname>[a-zA-Z][^<]*) *<(?P<email>[^>]*)> *(\\(@(?P<username>[^\\)]+)\\)|).*$")
	match := re.FindStringSubmatch(line)
	return &Maintainer{
		Active:   match[commentIndex] == "",
		Target:   path.Base(path.Clean(match[targetIndex])),
		Username: strings.Trim(match[usernameIndex], " \t"),
		Email:    strings.Trim(match[emailIndex], " \t"),
		FullName: strings.Trim(match[fullnameIndex], " \t"),
		Raw:      line,
	}
}
