import os
import gc
import time
import math
import random
from typing import List

import redis
import numpy as np
from throttled.utils import Benchmark


url: str = "redis://127.0.0.1:6379/0"
client: redis.Redis = redis.Redis.from_url(url)


# 素数求解
def is_prime(n: int) -> bool:
    if n <= 1:
        return False
    if n in (2, 3):
        return True
    if n % 2 == 0:
        return False
    max_divisor = math.isqrt(n) + 1
    for i in range(3, max_divisor, 2):
        if n % i == 0:
            return False
    return True


# 生成长度为 n 的斐波那契数列
def fibonacci(n: int) -> List[int]:
    """斐波那契数列"""
    if n <= 0:
        return []
    elif n == 1:
        return [0]
    elif n == 2:
        return [0, 1]
    else:
        fib = [0, 1]
        for i in range(2, n):
            fib.append(fib[i-1] + fib[i-2])
        return fib


# 生成长度为 n 的斐波那契数列
def matrix_multiply(n: int) -> np.ndarray:
    """矩阵乘法"""
    a = np.random.rand(n, n)
    b = np.random.rand(n, n)
    c = np.dot(a, b)
    return c


# 执行 `SET KEY VALUE`
def redis_set():
    client.set("ping:baseline", 1)


def mixed(n: int) -> int:
    """混合 I/O 计算"""
    time.sleep(random.uniform(0.01, 0.05))
    result: int = sum(x ** 2 for x in range(n))
    return result


def benchmark_is_prime():
    benchmark: Benchmark = Benchmark()
    # benchmark.serial(is_prime, batch=1_000, n=536870909)
    benchmark.concurrent(is_prime, batch=1_000, workers=8, n=536870909)


def benchmark_fibonacci():
    benchmark: Benchmark = Benchmark()
    # benchmark.serial(fibonacci, batch=1_000, n=10_000)
    benchmark.concurrent(fibonacci, batch=1_000, workers=8, n=1_000)


def benchmark_matrix_multiply():
    benchmark: Benchmark = Benchmark()
    # benchmark.serial(matrix_multiply, batch=1_000, n=1000)
    benchmark.concurrent(matrix_multiply, batch=1_000, workers=8, n=1000)


def benchmark_mixed():
    benchmark: Benchmark = Benchmark()
    # benchmark.serial(mixed, batch=1_000, n=1000)
    benchmark.concurrent(mixed, batch=1_000, workers=8, n=1000)


def benchmark_redis():
    benchmark: Benchmark = Benchmark()
    # benchmark.serial(redis_set, batch=100_000)
    benchmark.concurrent(redis_set, batch=100_000, workers=8)


if __name__ == '__main__':
    # benchmark_is_prime()
    # gc.collect()
    #
    benchmark_fibonacci()
    gc.collect()
    #
    # benchmark_matrix_multiply()
    # gc.collect()
    #
    # benchmark_mixed()
    # gc.collect()

    # benchmark_redis()
    # gc.collect()
