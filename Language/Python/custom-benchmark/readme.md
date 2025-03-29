# 📊 用 Python 实现一个 Benchmark 工具
> 源码：[throttled-py/throttled/utils.py Benchmark](https://github.com/ZhuoZhuoCrayon/throttled-py/blob/main/throttled/utils.py)，期待 Star 🌟。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/custom-benchmark/images/overview.png" width="50%">
</div>

## 1. 程序离不开 Benchmark

Benchmark 表示性能基准测试，通过一定量级的串行 / 并发调用，来测试程序的吞吐量、响应时间、资源消耗等指标，从而帮助我们了解程序的性能表现，为程序的优化提供方向。

### 1.1. Go 语言的 Benchmark

在 Go 语言中，非常方便就可以对某个模块进行 Benchmark 测试，例如：

```go
package main
import (
	"fmt"
	"testing"
)

func greeter(name string) string {
	return fmt.Sprintf("Hi, %s!", name)
}

func BenchmarkExample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		greeter("LiHua")
	}
}
```

执行 `go test -bench=.` 命令，就可以看到 Benchmark 的结果：

```shell
$ go test -bench=.
goos: darwin
goarch: arm64
pkg: xxx
cpu: Apple M1 Pro
BenchmarkExample
BenchmarkExample-10    	17006762	        66.51 ns/op
PASS
```

其中，`BenchmarkExample-10` 表示并发数，`17006762` 表示执行次数，`66.51 ns/op` 表示每次执行的平均耗时。

### 1.2. Python 语言的 Benchmark

Python 并不像 Go，没有内置的 Benchmark 工具，当然社区有不错的实现，例如：[pytest-benchmark](https://github.com/ionelmc/pytest-benchmark)：
* Python 版本支持 `>=3.8, <=3.12`。
* 支持串行测试。

最近在开发一个第三方库，需要对比模块在串行 / 并发 / Async 上的性能表现，并且在 Python 3.13 版本上运行，没有找到合适的工具，于是就自己实现了一个：
* 支持多种执行模式：串行 / 并发 / 异步并发。
* 支持多种 Python 版本：3.8 ~ 3.13。
* 支持性能对比。


## 2. 快速开始

### 2.1. 安装

```shell
# 依赖版本：throttled-py >= 1.0.2
$ pip install throttled-py==1.0.2
```

#### 2.2. 编写一个 Benchmark

编写一个 Benchmark 程序，对 `redis-py` 单次 `SET key value` 进行压测。

```python
import asyncio
import redis
from redis import asyncio as aioredis

from throttled import utils

url: str = "redis://127.0.0.1:6379/0"
client: redis.Redis = redis.Redis.from_url(url)
async_client: aioredis.Redis = aioredis.Redis.from_url(url)


def baseline():
    client.set("ping:baseline", 1)


async def async_baseline():
    await async_client.set("ping:baseline", 1)


async def main():
    benchmark: utils.Benchmark = utils.Benchmark()
    try:
        # 1️⃣ 串行：顺序执行 100,000 次。
        benchmark.serial(baseline, 100_000)
        # 2️⃣ 多线程并发：32 工作线程，执行 100,000 次。
        benchmark.concurrent(baseline, 100_000, 32)
        # 3️⃣ 异步并发：32 工作协程，执行 100,000 次。
        await benchmark.async_concurrent(async_baseline, 100_000, 32)
    finally:
        client.close()
        await async_client.aclose()


if __name__ == "__main__":
    asyncio.run(main())
```

执行得到结果：

```shell
Python 3.13.1 (main, Mar 29 2025, 16:29:36) [Clang 15.0.0 (clang-1500.3.9.4)]
Implementation: CPython
OS: Darwin 23.6.0, Arch: arm64 

✅ Total: 100000, 🕒 Latency: 0.0588 ms/op, 🚀 Throughput: 16837 req/s (--)
✅ Total: 100000, 🕒 Latency: 1.9953 ms/op, 💤 Throughput: 15966 req/s (⬇️-5.17%)
✅ Total: 100000, 🕒 Latency: 1.0244 ms/op, 🚀 Throughput: 25227 req/s (⬆️58.00%)
```

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/custom-benchmark/images/2.2.png" width="80%">
</div>

### 2.3. pytest 集成

`Benchmark` 可以作为一个 `fixture`，和 `pytest` 结合使用：

```python
import pytest
from throttled.utils import Benchmark

def ping() -> str:
  return "pong"

@pytest.fixture(scope="class")
def benchmark() -> Benchmark:
    return Benchmark()

class TestPing:
  def test_ping(benchmark: Benchmark):
    benchmark.concurrent(ping, 100)
```

