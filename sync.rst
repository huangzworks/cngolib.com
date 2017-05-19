sync
==========

本文是 Go 标准库中 sync 包文档的翻译，
原文地址为： 
https://golang.org/pkg/sync/


概述
----------
..
    Package sync provides basic synchronization primitives such as mutual exclusion locks. 
    Other than the Once and WaitGroup types, 
    most are intended for use by low-level library routines. 
    Higher-level synchronization is better done via channels and communication.

    Values containing the types defined in this package should not be copied.

sync 包提供了诸如互斥锁这样的基本同步原语。
除了 Once 类型和 WaitGroup 类型之外，
包中的大多数其他类型都是为底层函数库程序准备的。
高层次的同步最好还是通过 channel 以及 communication 来完成。

包含本包定义的类型的值不应该进行拷贝。


Mutex 类型
--------------------

一个 Mutex 就是一个互斥锁，
这种锁可以用作其他结构的一部分。

Mutex 的零值是一个未上锁的互斥锁 。

::

    type Mutex struct {
        // contains filtered or unexported fields
    }

(\*Mutex) Lock 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (m *Mutex) Lock()

对 m 进行加锁。

如果 m 已经被加锁，
那么执行该方法的 goroutine 将被阻塞直到 m 可用为止。


(\*Mutex) Unlock 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (m *Mutex) Unlock()

解锁 m 。
如果 m 并未加锁，
那么引发一个运行时错误。

被加锁的 Mutex 并不与特定的 goroutine 绑定，
在一个 goroutine 里面对 Mutex 进行加锁，
然后在另一个 goroutine 里面对 Mutex 进行解锁，
这是完全可行的。


WaitGroup 类型
-------------------

..
    A WaitGroup waits for a collection of goroutines to finish. 
    The main goroutine calls Add to set the number of goroutines to wait for. 
    Then each of the goroutines runs and calls Done when finished. 
    At the same time, 
    Wait can be used to block until all goroutines have finished.

    A WaitGroup must not be copied after first use.

一个 WaitGroup 会等待一系列 goroutine 直到它们全部运行完毕为止。
主 goroutine 通过调用 Add 方法来设置需要等待的 goroutine 数量，
而每个运行的 goroutine 则在它们运行完毕时调用 Done 方法。
与此同时，
调用 Wait 方法可以阻塞直到所有 goroutine 都运行完毕为止。

::

    type WaitGroup struct {
        // 其他已过滤或者未导出字段……
    }

示例
^^^^^^^^^^

这个示例会并发地获取给定的多个 URL ，
并使用 WaitGroup 进行阻塞，
直到所有获取操作都已完成为止。

.. literalinclude:: code/sync/wait_group.go


(\*WaitGroup) Add 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (wg *WaitGroup) Add(delta int)

..
    Add adds delta, 
    which may be negative, 
    to the WaitGroup counter. 
    If the counter becomes zero, 
    all goroutines blocked on Wait are released. 
    If the counter goes negative, 
    Add panics.

为 WaitGroup 计数器加上给定的增量 delta ，
其中 delta 的值可以为负数。
当计数器的值变为 0 时，
所有被 Wait 阻塞的 goroutine 都将被释放。
当计数器的值变为负数时，
Add 调用将引发一个 panic 。

..
    Note that calls with a positive delta that occur when the counter is zero must happen before a Wait. 
    Calls with a negative delta, 
    or calls with a positive delta that start when the counter is greater than zero, 
    may happen at any time. 
    Typically this means the calls to Add should execute before the statement creating the goroutine or other event to be waited for. 
    If a WaitGroup is reused to wait for several independent sets of events, 
    new Add calls must happen after all previous Wait calls have returned. 
    See the WaitGroup example.

当计数器的值为 0 时，
delta 的值只能为正数，
并且对 Add 的调用必须出现在 Wait 调用之前；
当计数器的值大于 0 时，
delta 的值既可以是正数也可以是负数，
并且对 Add 的调用可以发生在任何时候。
简单来说，
这意味着 Add 必须在创建 goroutine 的语句之前调用，
又或者在其他需要等待的事件之前调用。

在重复使用同一个 WaitGroup 对不同的独立事件集合（independent sets of events）进行等待时，
新的 Add 调用必须发生在之前的所有 Wait 调用均已返回的情况下。

(\*WaitGroup) Done 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (wg *WaitGroup) Done()

对 WaitGroup 计数器执行减一操作。

(\*WaitGroup) Wait 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (wg *WaitGroup) Wait()

阻塞直至 WaitGroup 计数器的值为 0 。
