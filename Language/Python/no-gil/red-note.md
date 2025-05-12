# 快但没完全快｜Python no-GIL 性能对比

❓Python 性能为什么上不去？GIL 使得多线程执行需频繁获取/释放 GIL，从而让并发得到近似单线程的效果。

📢 PEP703 提到在 CPython 中全局解释器锁（GIL）将成为可选项目（Making the Global Interpreter Lock Optional in CPython）。

🎉 接着上一篇笔记，这一期来系统分析 Python no-GIL 实验版本的性能表现。

1️⃣ 单线程场景（图 1）
>> ⚙️ Pyperformance 是 Python 的官方基准测试套件，用于测量和比较不同版本 Python 解释器的单线程运行性能。

>> 🧪 选取计算密集型的 benchmarks，用于反映单线程执行性能。

>> 🐌 no-GIL 实验版本在单线程场景下性能显著下降（⬇️ 100 %+），可能为保证线程安全，引入额外开销

2️⃣ 多线程场景（图 2）
>> ⏳构造计算密集、IO 密集型原子任务，在 8-threads 模式下进行性能分析。

>> ✅ no-GIL 在多线程处理计算密集型任务（is_prime）上具有较好的性能表现（⬆️ 290 %+）。

>> 💤 涉及申请大量内存（fibonacci）时性能下降明显（⬇️ 50 %+）。

>> 👍 IO 密集型场景下（Redis SET），性能显著提升（⬆️ 130 %+）。

>> ➖ numpy 底层为 C 实现，矩阵运算性能持平。

🧐 目前来看，no-GIL 在频繁申请内存的场景下仍有提升空间，单线程性能下滑问题有待解决～
