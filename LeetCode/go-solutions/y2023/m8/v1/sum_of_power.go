package main

import (
	"sort"
)

func sumOfPower1(nums []int) int {
	res := 0
	mod := int(1e9 + 7)
	dp := make([]int, len(nums))
	preSum := make([]int, len(nums))
	for i := 0; i < len(nums); i++ {
		if i == 0 {
			// 初始值，只有一个元素的情况下，dp 是自己，preSum 也是
			dp[i] = nums[i]
			preSum[i] = dp[i]
		} else {
			dp[i] = (nums[i] + preSum[i-1]) % mod
			preSum[i] = (preSum[i-1] + dp[i]) % mod
		}
		res = (res + (nums[i]*nums[i]%mod*dp[i])%mod) % mod
	}
	return res
}

func sumOfPower2(nums []int) int {
	dp := 0
	res := 0
	preSum := 0
	mod := int(1e9 + 7)
	for i := 0; i < len(nums); i++ {
		dp = (nums[i] + preSum) % mod
		preSum = (dp + preSum) % mod
		res = (res + (nums[i] * nums[i] % mod * dp % mod)) % mod
	}
	return res
}

func sumOfPower(nums []int) int {
	// 结果和最大值最小值有关，先对 nums 进行排序
	sort.Ints(nums)

	// 由于 nums[0] ~ nums[i] 的序列最大值为 num[i]
	// 可推出结果 F(i) = F(i - 1) + nums[i]^2 * dp(i)，其中 dp(i) 表示序列 nums[0] ~ nums[i - 1] 的最小值之和
	// dp(i) = nums[i](自己的最小值是自己) + (sum(dp(i - 1), dp(i - 2), ...,dp(1)))
	return sumOfPower2(nums)
}

func main() {
	println(sumOfPower([]int{1, 2, 3, 4}))
}
