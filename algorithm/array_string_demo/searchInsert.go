// 搜索插入位置

func searchInsert(nums []int, target int) int {
	for i, v := range nums {
		if v >= target {
			return i
		}
	}

	return len(nums)
}
