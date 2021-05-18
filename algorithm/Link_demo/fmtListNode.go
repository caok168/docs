// 从尾到头打印链表

func fmtListNode(head *ListNode) {
	l := list.New()
	for ; head != nil; head = head.Next {
	   l.PushFront(head.Val)
	}
 ​
	for item := l.Front(); item != nil; item = item.Next() {
	   fmt.Println(item.Value)
	}
 }
 
