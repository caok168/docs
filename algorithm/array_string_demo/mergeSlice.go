// 合并区间

type MergeSlice [][]int

func (s MergeSlice) Len() int {return len(s)}
func (s MergeSlice) Swap(i, j int) {s[i], s[j] = s[j], s[i]}
func (s MergeSlice) Less(i, j int) bool {return s[i][0] < s[j][0]}

func merge(intervals [][]int) [][]int {
	sort.Sort(MergeSlice(intervals))

	retList := make([][]int, 0)
	size := len(intervals)

	var i int
	for i < size {
		left := intervals[i][0]
		right := intervals[i][1]

		j := i+1
		for j < size {
			if intervals[j][0] <= right {
				right = Max(intervals[j][1], right)
				j++
			} else {
				break
			}
		}

		ret := []int{left, right}
		retList = append(retList, ret)
		i = j
	}

	return retList


}

