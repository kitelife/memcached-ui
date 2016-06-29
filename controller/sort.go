package controller

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	kb, mb, gb, tb, pb float64
)

func init() {
	kb = 1024
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
	pb = 1024 * tb
}

func ToHuman(size int64) string {
	fsize := float64(size)
	var negative bool
	if fsize < 0 {
		fsize = -fsize
		negative = true
	}

	var toHuman string
	switch {
	case fsize > pb:
		toHuman = fmt.Sprintf("%.3f PB", fsize/pb)
	case fsize > tb:
		toHuman = fmt.Sprintf("%.3f TB", fsize/tb)
	case fsize > gb:
		toHuman = fmt.Sprintf("%.3f GB", fsize/gb)
	case fsize > mb:
		toHuman = fmt.Sprintf("%.3f MB", fsize/mb)
	case fsize > kb:
		toHuman = fmt.Sprintf("%.3f KB", fsize/kb)
	default:
		toHuman = fmt.Sprintf("%d B", size)
	}
	if negative {
		toHuman = fmt.Sprintf("-%s", toHuman)
	}
	return toHuman
}

type Hosts []string

func (h Hosts) Len() int {
	return len(h)
}
func (h Hosts) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h Hosts) Less(i, j int) bool {
	return hostLess(h[i], h[j])
}
func hostLess(host1, host2 string) bool {
	host1s := strings.SplitN(host1, ":", 2)
	host1 = host1s[0]
	host2s := strings.SplitN(host2, ":", 2)
	host2 = host2s[0]
	parts1 := strings.SplitN(host1, ".", 5)
	parts2 := strings.SplitN(host2, ".", 5)
	if len(parts1) < 4 || len(parts2) < 4 {
		return host1 < host2
	}
	for i := 0; i < 4; i++ {
		ipInt1, _ := strconv.Atoi(parts1[i])
		ipInt2, _ := strconv.Atoi(parts2[i])
		if ipInt1 == ipInt2 {
			continue
		} else {
			return ipInt1 < ipInt2
		}
	}
	return false
}

type ValSorter struct {
	Keys []string
	Vals []int
}

func NewValSorter(m map[string]int) *ValSorter {
	vs := &ValSorter{
		Keys: make([]string, 0, len(m)),
		Vals: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		vs.Vals = append(vs.Vals, v)
	}
	return vs
}
func (vs *ValSorter) Len() int           { return len(vs.Vals) }
func (vs *ValSorter) Less(i, j int) bool { return vs.Vals[i] < vs.Vals[j] }
func (vs *ValSorter) Swap(i, j int) {
	vs.Vals[i], vs.Vals[j] = vs.Vals[j], vs.Vals[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
