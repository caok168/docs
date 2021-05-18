// 回文链表
// 时间复杂度和空间复杂度都是O(n)

func isPalindrome(head *ListNode) bool {
    vals := []int{}
    for ; head != nil; head = head.Next {
        vals = append(vals, head.Val)
    }
    n := len(vals)
    for i, v := range vals[:n/2] {
        if v != vals[n-1-i] {
            return false
        }
    }
    return true
}


