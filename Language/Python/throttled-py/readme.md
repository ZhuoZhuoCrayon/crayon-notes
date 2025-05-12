# throttled-py - 支持多种策略及存储选项的 Python 限流库

## throttled-py
> GitHub: https://github.com/ZhuoZhuoCrayon/throttled-py
> 
> 简介：🔧 支持多种算法（固定窗口，滑动窗口，令牌桶，漏桶 & GCRA）及存储（Redis、内存）的高性能 Python 限流库。


## ✨ 功能

* 提供线程安全的存储后端：[Redis](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#redis)、[内存（支持 Key 过期淘汰）](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#memory)。
* 支持多种限流算法：[固定窗口](https://github.com/ZhuoZhuoCrayon/throttled-py/tree/main/docs/basic#21-%E5%9B%BA%E5%AE%9A%E7%AA%97%E5%8F%A3%E8%AE%A1%E6%95%B0%E5%99%A8)、[滑动窗口](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#22-%E6%BB%91%E5%8A%A8%E7%AA%97%E5%8F%A3)、[令牌桶](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#23-%E4%BB%A4%E7%89%8C%E6%A1%B6)、[漏桶](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#24-%E6%BC%8F%E6%A1%B6) & [通用信元速率算法（Generic Cell Rate Algorithm, GCRA）](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/docs/basic/readme.md#25-gcra)。
* 支持[配置限流算法](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#3%E6%8C%87%E5%AE%9A%E9%99%90%E6%B5%81%E7%AE%97%E6%B3%95)，提供灵活的[配额设置](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#4%E6%8C%87%E5%AE%9A%E5%AE%B9%E9%87%8F)。
* 支持即刻返回及[等待重试](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#%E7%AD%89%E5%BE%85%E9%87%8D%E8%AF%95)，提供[函数调用](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#%E5%87%BD%E6%95%B0%E8%B0%83%E7%94%A8)、[装饰器](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#%E4%BD%9C%E4%B8%BA%E8%A3%85%E9%A5%B0%E5%99%A8)、[上下文管理器](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#%E4%B8%8A%E4%B8%8B%E6%96%87%E7%AE%A1%E7%90%86%E5%99%A8)。
* 良好的性能，单次限流 API 执行耗时换算如下（详见 [Benchmarks](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#-benchmarks)）：
  * 内存：约为 2.5 ~ 4.5 次 `dict[key] += 1` 操作。
  * Redis：约为 1.06 ~ 1.37 次 `INCRBY key increment` 操作。


## 🔰 安装

```shell
$ pip install throttled-py
```

如需使用扩展功能，可通过以下方式安装可选依赖项（多个依赖项用逗号分隔）：

```shell
$ pip install "throttled-py[redis]"

$ pip install "throttled-py[redis,in-memory]"
```


## 🎨 快速开始

### 1）通用 API

* `limit`：消耗请求，返回 [**RateLimitResult**](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#1ratelimitresult)。
* `peek`：获取指定 Key 的限流器状态，返回 [**RateLimitState**](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#2ratelimitstate)。

### 2）样例

```python
from throttled import RateLimiterType, Throttled, rate_limiter, store, utils

throttle = Throttled(
    # 📈 使用令牌桶作为限流算法。
    using=RateLimiterType.TOKEN_BUCKET.value,
    # 🪣 设置配额：每秒填充 1,000 个 Token（limit），桶大小为 1,000（burst）。
    quota=rate_limiter.per_sec(1_000, burst=1_000),
    # 📁 使用内存作为存储
    store=store.MemoryStore(),
)

def call_api() -> bool:
    # 💧消耗 Key=/ping 的一个 Token。
    result = throttle.limit("/ping", cost=1)
    return result.limited

if __name__ == "__main__":
    # 💻 Python 3.12.10, Linux 5.4.119-1-tlinux4-0009.1, Arch: x86_64, Specs: 2C4G.
    # ✅ Total: 100000, 🕒 Latency: 0.0068 ms/op, 🚀 Throughput: 122513 req/s (--)
    # ❌ Denied: 98000 requests
    benchmark: utils.Benchmark = utils.Benchmark()
    denied_num: int = sum(benchmark.serial(call_api, 100_000))
    print(f"❌ Denied: {denied_num} requests")
```

### 3）作为装饰器

```python
from throttled import Throttled, exceptions, rate_limiter

# 创建一个每分钟允许通过 1 次的限流器。
@Throttled(key="/ping", quota=rate_limiter.per_min(1))
def ping() -> str:
    return "ping"

ping()
try:
    ping()  # 当触发限流时，抛出 LimitedError。
except exceptions.LimitedError as exc:
    print(exc)  # Rate limit exceeded: remaining=0, reset_after=60, retry_after=60
```

### 4）上下文管理器

你可以使用「上下文管理器」对代码块进行限流，允许通过时，返回 [**RateLimitResult**](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#1ratelimitresult)。

触发限流或重试超时，抛出 [**LimitedError**](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/README_ZH.md#limitederror)。

```python
from throttled import Throttled, exceptions, rate_limiter

def call_api():
    print("doing something...")

throttle: Throttled = Throttled(key="/api/v1/users/", quota=rate_limiter.per_min(1))
with throttle as rate_limit_result:
    print(f"limited: {rate_limit_result.limited}")
    call_api()

try:
    with throttle:
        call_api()
except exceptions.LimitedError as exc:
    print(exc)  # Rate limit exceeded: remaining=0, reset_after=60, retry_after=60
```

