package class

type Int64List []int64

func (l Int64List) Len() int           { return len(l) }
func (l Int64List) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Int64List) Less(i, j int) bool { return l[i] < l[j] }
