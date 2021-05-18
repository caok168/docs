// 中间节点

//思路一，双指针法
func middleNode(head *ListNode) *ListNode {
    fast, slow := head, head
    for {
        if fast.Next != nil && fast.Next.Next != nil {
            fast = fast.Next.Next
            slow = slow.Next
        } else if fast.Next != nil {
            return slow.Next
        } else {
            return slow
        }
    }
}


//思路二：转存数组取中间值
func middleNode(head *ListNode) *ListNode {
    var list []*ListNode
    //转存数组
    for head!=nil {
       list= append(list,head)
       head =  head.Next
    }
    return list[len(list)/2]
}


