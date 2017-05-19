testing
==================

本文是 Go 标准库中 testing 包文档的翻译，
原文地址为： 
https://golang.org/pkg/testing/


类型 B
---------------

..
    B is a type passed to Benchmark functions 
    to manage benchmark timing 
    and to specify the number of iterations to run.

B 是传递给基准测试函数的一种类型，
它用于管理基准测试的计时行为，
并指示应该迭代地运行测试多少次。

..
    A benchmark ends when its Benchmark function returns 
    or calls any of the methods FailNow, Fatal, Fatalf, SkipNow, Skip, or Skipf. 
    Those methods must be called only from the goroutine running the Benchmark function. 
    The other reporting methods, 
    such as the variations of Log and Error, 
    may be called simultaneously from multiple goroutines.

一个基准测试在它的基准测试函数返回时，
又或者在它的基准测试函数调用 FailNow 、 Fatal 、 Fatalf 、 SkipNow 、 Skip 或者 Skipf 中的任意一个方法时，
测试即宣告结束。
至于其他报告方法，
比如 Log 和 Error 的变种，
则可以在其他 goroutine 中同时进行调用。

..
    Like in tests, 
    benchmark logs are accumulated during execution 
    and dumped to standard error when done. 
    Unlike in tests, 
    benchmark logs are always printed, 
    so as not to hide output whose existence may be affecting benchmark results.

跟测试一样，
基准测试会在执行的过程中积累日志，
并在测试完毕时将日志转储到标准错误。
但跟测试不一样的是，
为了避免基准测试的结果受到日志打印操作的影响，
基准测试总是会把日志打印出来。

::

    type B struct {
        N int
        // contains filtered or unexported fields
    }

(\*B) Error 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Error(args ...interface{})

调用 Error 相当于在调用 Log 之后调用 Fail 。

(\*B) Errorf 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Errorf(format string, args ...interface{})

调用 Errorf 相当于在调用 Logf 之后调用 Fail 。

(\*B) Fail 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Fail()

将当前的测试函数标识为“失败”，
但仍然继续执行该函数。

(\*B) FailNow 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) FailNow()

将当前的测试函数标识为“失败”，
并停止执行该函数。
在此之后， 
测试过程将在下一个测试或者下一个基准测试中继续。

FailNow 必须在运行测试函数或者基准测试函数的 goroutine 中调用， 
而不能在测试期间创建的 goroutine 中调用。 
调用 FailNow 不会导致其他 goroutine 停止。

(\*B) Failed 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Failed() bool

Failed 用于报告测试函数是否已失败。

(\*B) Fatal 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Fatal(args ...interface{})

调用 Fatal 相当于在调用 Log 之后调用 FailNow 。

(\*B) Fatalf 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Fatalf(format string, args ...interface{})

调用 Fatalf 相当于在调用 Logf 之后调用 FailNow 。

(\*B) Log 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Log(args ...interface{})

Log 使用与 Printf 相同的格式化语法对它的参数进行格式化，
然后将格式化后的文本记录到错误日志里面。
如果输入的格式化文本最末尾没有出现新行，
那么将一个新行添加到格式化后的文本末尾。

对于测试来说，
Logf 产生的格式化文本只会在测试失败或者设置了 -test.v 标志的情况下被打印出来；
对于基准测试来说，
为了避免 -test.v 标志的值对测试的性能产生影响，
Logf 产生的格式化文本总会被打印出来。

(\*B) Logf 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Logf(format string, args ...interface{})

Log 使用与 Printf 相同的格式化语法对它的参数进行格式化，
然后将格式化后的文本记录到错误日志里面。
如果输入的格式化文本最末尾没有出现新行，
那么将一个新行添加到格式化后的文本末尾。

对于测试来说，
Logf 产生的格式化文本只会在测试失败或者设置了 -test.v 标志的情况下被打印出来；
对于基准测试来说，
为了避免 -test.v 标志的值对测试的性能产生影响，
Logf 产生的格式化文本总会被打印出来。

(\*B) Name 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Name() string

返回正在运行的测试或者基准测试的名字。

(\*B) ReportAllocs 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) ReportAllocs()

打开当前基准测试的内存统计功能，
与使用 -test.benchmem 设置类似，
但 ReportAllocs 只影响那些调用了该函数的基准测试。

(\*B) ResetTimer 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) ResetTimer()

对已经逝去的基准测试时间以及内存分配计数器进行清零。
对于正在运行中的计时器，
这个方法不会产生任何效果。

(\*B) Run 方法
^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) Run(name string, f func(b *B)) bool

执行名字为 name 的子基准测试（subbenchmark）f ，
并报告 f 在执行过程中是否出现了任何失败。

子基准测试跟其他普通的基准测试一样。
一个调用了 Run 方法至少一次的基准测试将不会对其自身进行测量（measure），
并且在 N 为 1 时，
这个基准测试将只会被执行一次。

Run 可以同时在多个 goroutine 中被调用，
但这些调用必须发生在 b 的外部基准函数（outer benchmark function）返回之前。

(\*B) RunParallel 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) RunParallel(body func(*PB))

以并行的方式执行给定的基准测试。
RunParallel 会创建出多个 goroutine ，
并将 b.N 个迭代分配给这些 goroutine 执行，
其中 goroutine 数量的默认值为 GOMAXPROCS 。
用户如果想要增加非CPU受限（non-CPU-bound）基准测试的并行性，
那么可以在 RunParallel 之前调用 SetParallelism 。
RunParallel 通常会与 -cpu 标志一同使用。

body 函数将在每个 goroutine 中执行，
这个函数需要设置所有 goroutine 本地的状态，
并迭代直到 pb.Next 返回 false 值为止。
因为 StartTimer 、 StopTimer 和 ResetTimer 这三个函数都带有全局作用，
所以 body 函数不应该调用这些函数；
除此之外，
body 函数也不应该调用 Run 函数。

执行示例：

.. literalinclude:: code/testing/run_parallel_test.go

输出：

::

    $ go test -v -bench .
    PASS
    ok      _/Users/huangz/cngolib.com/code/testing 0.006s

(\*B) SetBytes 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) SetBytes(n int64)

记录在单个操作中处理的字节数量。
在调用了这个方法之后，
基准测试将会报告 ns/op 以及 MB/s 。

(\*B) SetParallelism 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) SetParallelism(p int)

将 RunParallel 使用的 goroutine 数量设置为 p*GOMAXPROCS ，
如果 p 小于 1 ，
那么调用将不产生任何效果。

CPU受限（CPU-bound）的基准测试通常不需要调用这个方法。

(\*B) Skip 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Skip(args ...interface{})

调用 Skip 相当于在调用 Log 之后调用 SkipNow 。

(\*B) SkipNow 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) SkipNow()

将当前测试标识为“被跳过”并停止执行该测试。 
如果一个测试在失败（参考 Error 、 Errorf 和 Fail）之后被跳过了， 
那么它还是会被判断为是“失败的”。

在停止当前测试之后， 
测试过程将在下一个测试或者下一个基准测试中继续， 
具体请参考 FailNow 。

SkipNow 必须在运行测试的 goroutine 中进行调用， 
而不能在测试期间创建的 goroutine 中调用。 
调用 SkipNow 不会导致其他 goroutine 停止。

(\*B) Skipf 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Skipf(format string, args ...interface{})

调用 Skipf 相当于在调用 Logf 之后调用 SkipNow 。

(\*B) Skipped 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *B) Skipped() bool

报告测试是否已被跳过。

(\*B) StartTimer 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) StartTimer()

开始对测试进行计时。

这个函数在基准测试开始时会自动被调用，
它也可以在调用 StopTimer 之后恢复进行计时。

(\*B) StopTimer 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (b *B) StopTimer()

停止对测试进行计时。

当你需要执行一些复杂的初始化操作，
并且你不想对这些操作进行测量时，
就可以使用这个方法来暂时地停止计时。


类型 T
---------------

..
    T is a type passed to Test functions to manage test state and support formatted test logs. 
    Logs are accumulated during execution 
    and dumped to standard output when done.

T 是传递给测试函数的一种类型，
它用于维护测试状态并支持格式化的测试日志。
测试日志会在执行测试的过程中不断累积，
并在测试完成时转储至标准输出。

..
    A test ends when its Test function returns 
    or calls any of the methods FailNow, Fatal, Fatalf, SkipNow, Skip, or Skipf. 
    Those methods, 
    as well as the Parallel method, 
    must be called only from the goroutine running the Test function.

当一个测试的测试函数返回时，
又或者当一个测试函数调用 FailNow 、 Fatal 、 Fatalf 、 SkipNow 、 Skip 或者 Skipf 中的任意一个时，
该测试即宣告结束。
跟 Parallel 方法一样，
以上提到的这些方法只能在运行测试函数的 goroutine 中调用。

..
    The other reporting methods, 
    such as the variations of Log and Error, 
    may be called simultaneously from multiple goroutines.

至于其他报告方法，
比如 Log 以及 Error 的变种，
则可以在多个 goroutine 中同时进行调用。

::

    type T struct {
        // contains filtered or unexported fields
    }


(\*T) Error 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Error(args ...interface{})

调用 Error 相当于在调用 Log 之后调用 Fail 。

(\*T) Errorf 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Errorf(format string, args ...interface{})

调用 Errorf 相当于在调用 Logf 之后调用 Fail 。

(\*T) Fail 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    
    func (c *T) Fail()

将当前测试标识为失败，
但是仍继续执行该测试。

(\*T) FailNow 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) FailNow()

将当前测试标识为失败并停止执行该测试，
在此之后，
测试过程将在下一个测试或者下一个基准测试中继续。

FailNow 必须在运行测试函数或者基准测试函数的 goroutine 中调用，
而不能在测试期间创建的 goroutine 中调用。
调用 FailNow 不会导致其他 goroutine 停止。


(\*T) Failed 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Failed() bool

Failed 用于报告测试函数是否已失败。

(\*T) Fatal 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Fatal(args ...interface{})

调用 Fatal 相当于在调用 Log 之后调用 FailNow 。

(\*T) Fatalf 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Fatalf(format string, args ...interface{})

调用 Fatalf 相当于在调用 Logf 之后调用 FailNow 。

(\*T) Log 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Log(args ...interface{})

Log 使用与 Println 相同的格式化语法对它的参数进行格式化，
然后将格式化后的文本记录到错误日志里面：

- 对于测试来说，
  格式化文本只会在测试失败或者设置了 -test.v 标志的情况下被打印出来；

- 对于基准测试来说，
  为了避免 -test.v 标志的值对测试的性能产生影响，
  格式化文本总会被打印出来。

(\*T) Logf 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Logf(format string, args ...interface{})

Log 使用与 Printf 相同的格式化语法对它的参数进行格式化，
然后将格式化后的文本记录到错误日志里面。
如果输入的格式化文本最末尾没有出现新行，
那么将一个新行添加到格式化后的文本末尾。

对于测试来说，
Logf 产生的格式化文本只会在测试失败或者设置了 -test.v 标志的情况下被打印出来；
对于基准测试来说，
为了避免 -test.v 标志的值对测试的性能产生影响，
Logf 产生的格式化文本总会被打印出来。

(\*T) Name 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Name() string

返回正在运行的测试或基准测试的名字。

(\*T) Parallel 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (t *T) Parallel()

Parallel 用于表示当前测试只会与其他带有 Parallel 方法的测试并行进行测试。

(\*T) Run 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (t *T) Run(name string, f func(t *T)) bool

执行名字为 name 的子测试 f ，
并报告 f 在执行过程中是否出现了任何失败。
Run 将一直阻塞直到 f 的所有并行测试执行完毕。

Run 可以在多个 goroutine 里面同时进行调用，
但这些调用必须发生在 t 的外层测试函数（outer test function）返回之前。

(\*T) Skip 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Skip(args ...interface{})

调用 Skip 相当于在调用 Log 之后调用 SkipNow 。

(\*T) SkipNow 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) SkipNow()

将当前测试标识为“被跳过”并停止执行该测试。
如果一个测试在失败（参考 Error 、 Errorf 和 Fail）之后被跳过了，
那么它还是会被判断为是“失败的”。

在停止当前测试之后，
测试过程将在下一个测试或者下一个基准测试中继续，
具体请参考 FailNow 。

SkipNow 必须在运行测试的 goroutine 中进行调用，
而不能在测试期间创建的 goroutine 中调用。
调用 SkipNow 不会导致其他 goroutine 停止。

(\*T) Skipf 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Skipf(format string, args ...interface{})

调用 Skipf 相当于在调用 Logf 之后调用 SkipNow 。


(\*T) Skipped 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (c *T) Skipped() bool

Skipped 用于报告测试函数是否已被跳过。



TB 类型
-------------

TB 类型同时拥有 T 类型和 B 类型提供的接口。

::

    type TB interface {
        Error(args ...interface{})
        Errorf(format string, args ...interface{})
        Fail()
        FailNow()
        Failed() bool
        Fatal(args ...interface{})
        Fatalf(format string, args ...interface{})
        Log(args ...interface{})
        Logf(format string, args ...interface{})
        Name() string
        Skip(args ...interface{})
        SkipNow()
        Skipf(format string, args ...interface{})
        Skipped() bool
        // contains filtered or unexported methods
    }
