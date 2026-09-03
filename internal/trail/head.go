package trail

import "sort"

// HeadDate is the commit date of HEAD as YYYY-MM-DD, or "" in a repo with no
// commits.
func (p Project) HeadDate() string {
	day, err := git(p.Root, "log", "-1", "--format=%cd", "--date=short")
	if err != nil {
		return ""
	}
	return day
}

// NewestChangeDate is the date of the newest changelog artifact, or "".
func (p Project) NewestChangeDate() string {
	arts := p.Artifacts(RowChangelog.Class())
	if len(arts) == 0 {
		return ""
	}
	dates := make([]string, 0, len(arts))
	for _, a := range arts {
		if a.Owner == RowChangelog {
			dates = append(dates, a.Date)
		}
	}
	if len(dates) == 0 {
		return ""
	}
	sort.Strings(dates)
	return dates[len(dates)-1]
}

// CoversHead reports whether the changelog trail covers the last commit. This
// is the commit-boundary check: work in progress is nobody's emergency, but a
// commit that no changelog entry describes is a trail that has fallen behind
// the code.
func (p Project) CoversHead() (ok bool, head, newest string) {
	head = p.HeadDate()
	newest = p.NewestChangeDate()
	if head == "" {
		return true, head, newest // no commits yet: vacuously covered
	}
	if newest == "" {
		return false, head, newest
	}
	return newest >= head, head, newest
}
