
class ListNode {
      int val;
      ListNode next;
      ListNode(int x) {
          val = x;
          next = null;
 }
 }

public class Solution {
    public ListNode detectCycle(ListNode head) {
        if (head == null || head.next == null) {
            return null;
        }

        ListNode commonNode = getCommonNode(head);

        if (commonNode == null) return null;

        ListNode currentNode = head;
        while (currentNode != commonNode) {
            currentNode = currentNode.next;
            commonNode = commonNode.next;
        }
        return commonNode;
    }

    private ListNode getCommonNode(ListNode head) {
        ListNode slow = head;
        ListNode fast = head;

        while (fast != null && fast.next != null) {
            slow = slow.next;
            fast = fast.next.next;

            if (slow == fast) {
                return slow;
            }
        }

        return null;
    }
}
