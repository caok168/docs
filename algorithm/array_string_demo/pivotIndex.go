// 寻找数组的中心索引

/*
* 解题思路：
* 全部元素之和为total，当遍历到第i个元素时，左侧元素之和为sum，右侧元素之和为 total - num(i) - sum
* 所有相等 sum = total - num(i) - sum
* 2*sum + num(i) = total
*/

func pivotIndex(nums []int) int {
	total := 0
	for _, v := range nums {
		total += v
	}

	sum := 0
	for i, v := range nums {
		if 2*sum + v == total {
			return i
		}

		sum += v
	}

	return -1
}