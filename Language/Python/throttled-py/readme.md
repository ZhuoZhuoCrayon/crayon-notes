# throttled-py - 支持多种策略及存储选项的 Python 限流库

## throttled-py
> GitHub: https://github.com/ZhuoZhuoCrayon/throttled-py
> 
> 简介：🔧 支持多种算法（固定窗口，滑动窗口，令牌桶，漏桶 & GCRA）及存储（Redis、内存）的高性能 Python 限流库。

## 🚀 功能

* 提供线程安全的存储后端：Redis（基于 Lua 实现限流算法）、内存（基于 threading.RLock，支持 Key 过期淘汰）。
* 支持多种限流算法：[固定窗口](https://github.com/ZhuoZhuoCrayon/throttled-py/tree/main/docs/basic#21-%E5%9B%BA%E5%AE%9A%E7%AA%97%E5%8F%A3%E8%AE%A1%E6%95%B0%E5%99%A8)、[滑动窗口](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#22-%E6%BB%91%E5%8A%A8%E7%AA%97%E5%8F%A3)、[令牌桶](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#23-%E4%BB%A4%E7%89%8C%E6%A1%B6)、[漏桶](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#24-%E6%BC%8F%E6%A1%B6) & [通用信元速率算法（Generic Cell Rate Algorithm, GCRA）](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#25-gcra)。
* 提供灵活的限流策略、配额设置 API，文档详尽。
* 支持装饰器模式。
* 良好的性能，单次限流 API 执行耗时换算如下（详见 [Benchmarks](https://github.com/ZhuoZhuoCrayon/throttled-py?tab=readme-ov-file#-benchmarks)）：
  * 内存：约为 2.5 ~ 4.5 次 `dict[key] += 1` 操作。
  * Redis：约为 1.06 ~ 1.37 次 `INCRBY key increment` 操作。


## 🔰 安装

```shell
$ pip install throttled-py
```

## 🔥 快速开始

### 1）通用 API

* `limit`：消耗请求，返回 [**RateLimitResult**](https://github.com/ZhuoZhuoCrayon/throttled-py?tab=readme-ov-file#1ratelimitresult)。
* `peek`：获取指定 Key 的限流器状态，返回 [**RateLimitState**](https://github.com/ZhuoZhuoCrayon/throttled-py?tab=readme-ov-file#2ratelimitstate)。

### 2）样例

```python
from throttled import RateLimiterType, Throttled, rate_limter, store, utils

throttle = Throttled(
    # 📈 使用令牌桶作为限流算法。
    using=RateLimiterType.TOKEN_BUCKET.value,
    # 🪣 设置配额：每分钟填充 1000 个 Token（limit），桶大小为 1000（burst）。
    quota=rate_limter.per_sec(1_000, burst=1_000),
    # 📁 使用内存作为存储
    store=store.MemoryStore(),
)


def call_api() -> bool:
    # 💧 消耗 Key=/ping 的一个 Token。
    result = throttle.limit("/ping", cost=1)
    return result.limited


if __name__ == "__main__":
    # ✅ Total: 100000, 🕒 Latency: 0.5463 ms/op, 🚀 Throughput: 55630 req/s (--)
    # ❌ Denied: 96314 requests
    benchmark: utils.Benchmark = utils.Benchmark()
    denied_num: int = sum(benchmark.concurrent(call_api, 100_000, workers=32))
    print(f"❌ Denied: {denied_num} requests")
```
