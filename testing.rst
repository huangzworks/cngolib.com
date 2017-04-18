testing
==================

本文是 Go 标准库中 testing 包文档的翻译，
原文地址为： 
https://golang.org/pkg/testing/


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
并报告 f 是否成功。
Run 将一直阻塞直到 f 的所有并行测试执行完毕。

Run 可以在多个 goroutine 里面同时进行调用，
但这些调用必须发生在 t 的外部函数返回之前。

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
