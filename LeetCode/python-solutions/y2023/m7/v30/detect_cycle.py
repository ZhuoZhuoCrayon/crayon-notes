# -*- coding: utf-8 -*-
from typing import Optional


# Definition for singly-linked list.
class ListNode:
    def __init__(self, x):
        self.val = x
        self.next = None


class Solution:
    def detectCycle(self, head: Optional[ListNode]) -> Optional[ListNode]:
        if head is None or head.next is None:
            return None

        common_node: Optional[ListNode] = self.getCommonNode(head)
        if common_node is None:
            return None
        current_node: Optional[ListNode] = head
        while current_node != common_node:
            current_node, common_node = current_node.next, common_node.next

        return current_node

    @classmethod
    def getCommonNode(cls, head: ListNode) -> Optional[ListNode]:
        slow, fast = head, head

        while not (fast is None or fast.next is None):
            slow, fast = slow.next, fast.next.next

            if slow == fast:
                return slow

        return None
