container/heap
===================

本文是 Go 标准库中 container/heap 包文档的翻译，
原文地址为： 
https://golang.org/pkg/container/heap/


概述
---------

..
    Package heap provides heap operations for any type that implements heap.Interface. 
    A heap is a tree with the property that each node is the minimum-valued node in its subtree.

    The minimum element in the tree is the root, at index 0.

包 heap 为所有实现了 heap.Interface 的类型提供堆操作。
一个堆即是一棵树，
这棵树的每个节点的值都比它的子节点的值要小，
而整棵树最小的值位于树根（root），
也即是索引 0 的位置上。

..
    A heap is a common way to implement a priority queue. 
    To build a priority queue, 
    implement the Heap interface with the (negative) priority as the ordering for the Less method, 
    so Push adds items while Pop removes the highest-priority item from the queue. 
    The Examples include such an implementation; 
    the file example_pq_test.go has the complete source.

堆是实现优先队列的一种常见方法。
为了构建优先队列，
用户在实现堆接口时，
需要让 Less() 方法返回逆序的结果，
这样就可以在使用 Push 添加元素的同时，
通过 Pop 移除队列中优先级最高的元素了。
具体的实现请看接下来展示的优先队列例子。

示例：整数堆
^^^^^^^^^^^^^^^^^^

.. literalinclude:: code/heap/int_heap.go

执行结果：

.. literalinclude:: code/heap/int_heap.txt

示例：优先队列
^^^^^^^^^^^^^^^^^^^^

.. literalinclude:: code/heap/priority_queue.go

执行结果：

.. literalinclude:: code/heap/priority_queue.txt


Fix 函数
------------

::

    func Fix(h Interface, i int)

..
    Fix re-establishes the heap ordering after the element at index i has changed its value. 
    Changing the value of the element at index i and then calling Fix is equivalent to, but less expensive than, 
    calling Remove(h, i) followed by a Push of the new value. 
    The complexity is O(log(n)) where n = h.Len().

在索引 i 上的元素的值发生变化之后，
重新修复堆的有序性。
先修改索引 i 上的元素的值然后再执行 Fix ，
跟先调用 Remove(h, i) 然后再使用 Push 操作将新值重新添加到堆里面的做法具有同等的效果，
但前者所需的计算量稍微要少一些。

Fix 函数的复杂度为 O(log(n)) ，
其中 n 等于 h.Len() 。


Init 函数
------------

::

    func Init(h Interface)

..
    A heap must be initialized before any of the heap operations can be used. 
    Init is idempotent with respect to the heap invariants 
    and may be called whenever the heap invariants may have been invalidated. 
    Its complexity is O(n) where n = h.Len().

在执行任何堆操作之前，
必须对堆进行初始化。
Init 操作对于堆不变性（invariants）具有幂等性，
无论堆不变性是否有效，
它都可以被调用。

Init 函数的复杂度为 O(n) ，
其中 n 等于 h.Len() 。


Pop 函数
----------

::

    func Pop(h Interface) interface{}

..
    Pop removes the minimum element (according to Less) from the heap and returns it. 
    The complexity is O(log(n)) where n = h.Len(). 
    It is equivalent to Remove(h, 0).

Pop 函数根据 Less 的结果，
从堆中移除并返回具有最小值的元素，
等同于执行 Remove(h, 0) 。

Pop 函数的复杂度为 O(log(n)) ，
其中 n 等于 h.Len() 。


Push 函数
------------

::

    func Push(h Interface, x interface{})

..
    Push pushes the element x onto the heap. 
    The complexity is O(log(n)) where n = h.Len().

Push 函数将值为 x 的元素推入到堆里面，
该函数的复杂度为 O(log(n)) ，
其中 n 等于 h.Len() 。


Remove 函数
--------------

::

    func Remove(h Interface, i int) interface{}

..  
    Remove removes the element at index i from the heap. 
    The complexity is O(log(n)) where n = h.Len().

Remove 函数将移除堆中索引为 i 的元素，
该函数的复杂度为 O(log(n)) ，
其中 n 等于 h.Len() 。


Interface 类型
-----------------

..
    Any type that implements heap.Interface may be used as a min-heap with the following invariants 
    (established after Init has been called or if the data is empty or sorted):

任何实现了 heap.Interface 接口的类型，
都可以用作带有以下不变性的最小堆，
（换句话说，
这个堆在为空、已排序或者调用 Init 之后，
应该具有以下性质）：

::

    !h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
    

..
    Note that Push and Pop in this interface are for package heap's implementation to call. 
    To add and remove things from the heap, 
    use heap.Push and heap.Pop.

注意，
这个接口中的 Push 和 Pop 都是由 heap 包的实现负责调用的。
因此用户在向堆添加元素又或者从堆中移除元素时，
需要使用 heap.Push 以及 heap.Pop ：

::

    type Interface interface {
        sort.Interface
        Push(x interface{}) // 将 x 添加为元素 Len()
        Pop() interface{}   // 移除并返回元素 Len() - 1
    }
