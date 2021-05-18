// 两个链表的第一个公共节点

// 双指针，追及相遇
func getIntersectionNode(headA, headB *ListNode) *ListNode {
    curA:=headA
    curB:=headB
    for curA!=curB{
        if curA==nil{
            curA=headA
        }else{
            curA=curA.Next
        }
        if curB==nil{
            curB=headB
        }else{
            curB=curB.Next
        }
    }
    return curA
}

// 解法二:map记录节点地址
func getIntersectionNode(headA, headB *ListNode) *ListNode {
    hashmap:=make(map[*ListNode]bool)
    for headA!=nil{
        hashmap[headA]=true
        headA=headA.Next
    }
    for headB!=nil{
        if hashmap[headB]{
            return headB
        }
        headB=headB.Next
    }
    return nil
}

