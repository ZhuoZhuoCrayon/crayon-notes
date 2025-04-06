# 并发性能提升 200%+ ！Python no-GIL 实验版本性能实测

## 1. GIL 机制解析

### 1.1. 什么是 GIL？

全局解释器锁（Global Interpreter Lock，简称 GIL）是 CPython 解释器的核心同步机制，其本质是：
* 单进程级别的互斥锁。
* 控制 Python 字节码执行的准入机制。
* 确保解释器内部状态的线程安全。

### 1.2. GIL 的运行时瓶颈

1）每个线程执行计算任务前，都需要获取 GIL。

2）进入网络 IO 等待时释放 GIL，其他线程抢占后继续执行计算任务。

3）多线程执行需频繁获取/释放 GIL，使得 Python 并发得到近似单线程的效果。

```mermaid
sequenceDiagram
    participant Web
    participant Thread1
    participant GIL
    participant Thread2

    activate Thread1 #FFA500
    Thread1 ->> GIL : acquire
    Thread1 ->> Thread1: run
    GIL ->> Thread1 : release
    deactivate Thread1
    
    activate Thread1
    Thread1 ->> Web : request(I/O)
    Web ->> Thread1 : response
    deactivate Thread1
    
    %% Thread2 竞争阶段
    loop 计算密集型任务
        activate Thread2
        Thread2 ->> GIL: acquire
        Thread2 ->> Thread2: run
        Thread2 ->> GIL: release
        deactivate Thread2
    end

    activate Thread1

    activate Thread1 #FFA500
    Thread1 ->> GIL : acquire
    Thread1 ->> Thread1: run
    GIL ->> Thread1 : release
    deactivate Thread1
```

### 1.3.  no-GIL

📢 [PEP703](https://discuss.python.org/t/a-steering-council-notice-about-pep-703-making-the-global-interpreter-lock-optional-in-cpython/30474) 提到在 CPython 中全局解释器锁（GIL）将成为可选项目（Making the Global Interpreter Lock Optional in CPython）：

* a. Short term, we add the no-GIL build as an experimental build mode, presumably in 3.13 (if it slips to 3.14, that is not a problem).

  > 短期内，将 no-GIL 构建添加为实验性构建模式。

* b. Mid-term, after we have confidence that there is enough community support to make production use of no-GIL viable, we make the no-GIL build supported but not the default (yet).

  >  中期，在我们确信有足够的社区支持使无 GIL 的生产使用变得可行之后，我们将支持无 GIL 构建，但暂时不将其作为默认构建。

* c. Long-term, we want no-GIL to be the default, and to remove any vestiges of the GIL (without unnecessarily breaking backward compatibility).

  >  长期来看，我们希望 no-GIL 成为默认方式，并删除 GIL 的所有痕迹。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/no-gil/images/1-3-x.jpg" width="30%">
</div>

🎉 好消息是，`3.13t-dev`  作为 no-GIL 实验版本已具备较强的可测试性，下文将对其性能进行对比。


## 2. 安装实验版本

1）通过 [pyenv](https://github.com/pyenv/pyenv) 安装实验版本：

```shell
$ pyenv install 3.13t-dev
```

2）验证已关闭 GIL：

```shell
python -c 'import sys; print(sys._is_gil_enabled())'
```


## 3. 性能实测

> 通过构造一个 Redis IO 密集型操作的场景，对比 Python 3.8.12 & 3.13t-dev 在串行 / 并行上的性能表现。

### 3.1. 对比基准

1）开发环境：MacBook M1 Pro 32 GB。

2）运行一个本地 Redis，版本：`7.0.12`。

```shell
$ redis-server&
```

3）获取本地 Redis `SET` 命令性能基准：**53,561 requests/sec**。

```shell
$ redis-benchmark -t set -n 10000

# Summary:
#  throughput summary: 53561.86 requests per second
#  latency summary (msec):
#          avg       min       p50       p95       p99       max
#        0.505     0.176     0.511     0.679     1.207     2.751
```

### 3.2. 编写 Python benchmark

1）安装 [throttled-py](https://github.com/ZhuoZhuoCrayon/throttled-py)，该库提供 Redis 及 Benchmark 依赖：

```shell
# 依赖版本：throttled-py >= 1.0.2
$ pip install throttled-py==1.0.2
```

2）编写一个压测程序：

* 场景（Redis）：`SET key value`。
* 串行：顺序执行 100,000 次。
* 多线程并发：16 工作线程，执行 100,000 次。

```python
import sys
import redis
from throttled import utils

url: str = "redis://127.0.0.1:6379/0"
client: redis.Redis = redis.Redis(connection_pool=redis.ConnectionPool.from_url(url=url))


def redis_baseline():
    client.set("ping:baseline", 1)


def main():
    try:
        print(f"version: {sys.version} \nis_gil_enabled: {sys._is_gil_enabled()}")
    except AttributeError:
        print(f"version: {sys.version} \nis_gil_enabled: True")

    benchmark: utils.Benchmark = utils.Benchmark()
    # 串行测试（单线程顺序执行）
    benchmark.serial(redis_baseline, 100_000)
    # 并发测试（16 工作线程）
    benchmark.current(redis_baseline, 100_000, workers=16)


if __name__ == "__main__":
    main()
```


### 3.3. 性能对比

#### 3.3.1. QPS

在不同 Python 版本间，运行 3.2 的 Benchmark 程序，记录 QPS。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/no-gil/images/3-2-1-x.png" width="65%">
</div>

| Python 版本        | 串行 x 100,000        | 线程池（16 workers）x 100,000 |
|------------------|---------------------|--------------------------|
| [Redis SET 基准]   | 53,561 requests/sec | 53,561 requests/sec      |
| Python 3.8.12    | 18,640 requests/sec | 13,069 requests/sec      |
| Python 3.13t-dev | 17,609 requests/sec | 🚀 39,494 requests /sec  |

对比得出：

* Python 3.8.12 串行甚至比多线程并发更快。
* Python 3.13t-dev 较 Python 3.8 多线程提升 **200%+**。

#### 3.3.2. CPU 使用率

1）Python 3.8.12：🐌 只能跑满一个核心（假多线程实锤）。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/no-gil/images/3-3-2-1x.png" width="60%">
</div>

2）Python 3.13t-dev：🔥跑满 6 个核。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/no-gil/images/3-3-2-2x.png" width="60%">
</div>


## 4. 扩展 Benchmarks 分析

### 4.1. Pyperformance / 单线程

#### 4.1.1. 什么是 Pyperformance？

[Pyperformance](https://github.com/python/pyperformance) 是 Python 的官方基准测试套件，用于测量和比较不同版本 Python 解释器的**单线程**运行性能。

#### 4.1.2. 运行 Benchmarks

1）运行 Benchmarks

```shell
# no-GIL
$ pyperformance run -b nbody,regex_v8,crypto_pyaes,json_dumps,logging -o no_gil_results.json --python=python3.13t-dev 
# with-GIL
$ pyperformance run -b nbody,regex_v8,crypto_pyaes,json_dumps,logging -o with_gil_results.json --python=python3.13.1
# Compare
$ pyperformance compare with_gil_results.json no_gil_results.json
```

2）性能对比

| 基准测试 *[1]*         | Python 3.13.1（with-GIL）*[2]* | 🐌 Python 3.13t-dev（no-GIL）*[3]* | 对比 *[4]*   |
|--------------------|------------------------------|----------------------------------|------------|
| **crypto_pyaes**   | 58.8 ms                      | 85.1 ms                          | ⬇️ 44.7 %  |
| **json_dumps**     | 7.96 ms                      | 10.05 ms                         | ⬇️ 26.3 %  |
| **logging_format** | 4.25 μs                      | 8.86 μs                          | ⬇️ 108.5 % |
| **logging_silent** | 75.6 ns                      | 156.0 ns                         | ⬇️ 106.3 % |
| **logging_simple** | 3.85 μs                      | 8.00 μs                          | ⬇️ 108.8 % |
| **nbody**          | 70.8 ms                      | 192.5 ms                         | ⬇️ 171.9 % |
| **regex_v8**       | 19.4 ms                      | 20.3 ms                          | ⬇️ 4.6 %   |

* *[1] 在 Pyperformance 中选取计算密集型的 benchmarks，用于反映单线程执行性能。*
* *[2] Python version: 3.13.1 (64-bit), Report on macOS-14.7.1-arm64-arm-64bit-Mach-O, Number of logical CPUs: 10。*
* *[3] Python version: 3.13.2+ (64-bit) revision 646b453, Report on macOS-14.7.1-arm64-arm-64bit-Mach-O, Number of logical CPUs: 10。*
* *[4] no-GIL 实验版本在单线程场景下性能显著下降，可能为保证线程安全，引入额外开销。*

### 4.2. 多线程场景

构造计算密集、IO 密集型原子任务，在 8-threads 模式下进行性能分析，详见 [no-gil/benchmarks/main.py](https://github.com/ZhuoZhuoCrayon/crayon-notes/blob/master/Language/Python/no-gil/benchmarks/main.py)。

<div align="left">
  <img src="https://github.com/ZhuoZhuoCrayon/crayon-notes/raw/master/Language/Python/no-gil/images/4.2.png" width="60%">
</div>

| 基准测试 *[1]*          | Python 3.13.1（with-GIL） | Python 3.13t-dev（no-GIL） | 对比 *[2]*       |
|---------------------|-------------------------|--------------------------|----------------|
| **is_prime**        | 2,493 requests/sec      | 9,768 requests/sec       | ⬆️ 292%  *[2]* |
| **fibonacci**       | 462 requests/sec        | 215 requests/sec         | ⬇️ 53.5% *[2]* |
| **matrix_multiply** | 108 requests/sec        | 103 requests/sec         | ➖ 持平 *[3]*     |
| **redis_set**       | 15,923 requests/sec     | 38,020 requests/sec      | ⬆️ 139% *[4]*  |

* *[1] is_prime、fibonacci、matrix_multiply 为计算密集型任务，redis_set 为 IO 密集型任务。*
  * *is_prime：求解  `2 ^ 29 - 3`  是否为素数。*
  * *fibonacci：生成长度为 n 的斐波那契数列。*
  * *matrix_multiply：n 阶矩阵乘法（numpy）。*
  * *redis_set：执行 `SET KEY VALUE`。*
* *[2] no-GIL 在多线程处理计算密集型任务（is_prime）上具有较好的性能表现，涉及申请大量内存（fibonacci）时性能表现不佳。*
* *[3] numpy 底层为 C 实现，性能持平。*
* *[4] IO 密集型场景下，性能显著提升。*


## 5. 结语
* GIL 的存在使得过往部分线程不安全的代码得以正常运行，这可能会是未来升级 no-GIL 的隐患。
* no-GIL 在 IO 密集型任务上具有较好的性能表现，但计算密集型任务上性能表现不佳，具有较大优化空间。

