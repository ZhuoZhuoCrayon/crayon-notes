package main

// flipgame from: https://leetcode.cn/problems/card-flipping-game/description/
func flipgame(fronts []int, backs []int) int {
	blackSet := make(map[int]bool)

	for i := 0; i < len(fronts); i++ {
		// 正反面数字相同，但在翻某个具体数字的时候，这种牌是没办法「规避」的，反之，只要去掉这类 case，任何牌都可以翻
		// 记录正反面数字一样的情况，具有相同数字 & 且朝正面的牌一定会「相同」，此时不可能作为最小值
		if fronts[i] == backs[i] {
			blackSet[fronts[i]] = true
		}
	}

	// 剩下的牌，只要翻或不翻时，现在正面的数字不和黑名单的值冲突，则都可以作为最小值
	min := 2001
	for i := 0; i < len(fronts); i++ {
		// 任意张进行翻转，可能是翻或不翻
		// 翻牌的情况，如果正面数字不和黑名单冲突，说明当前场上「正面」「具有相同数字」的牌可以通过翻面，构造符合「与任意一张卡片的正面的数字都不同」的条件
		if !blackSet[fronts[i]] {
			if fronts[i] < min {
				min = fronts[i]
			}
		}
		// 不翻牌的情况，如果反面数字不和黑名单冲突，说明当前场上「正面」「具有相同数字」的牌可以通过翻面，构造符合「与任意一张卡片的正面的数字都不同」的条件
		if !blackSet[backs[i]] {
			if backs[i] < min {
				min = backs[i]
			}
		}
	}

	if min == 2001 {
		return 0
	} else {
		return min
	}
}
