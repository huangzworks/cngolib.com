errors
=============

本文是 Go 标准库中 errors 包文档的翻译，
原文地址为： 
https://golang.org/pkg/errors/


概述
-------

errors 包实现了用于处理错误的函数。

示例：

.. literalinclude:: code/errors/example.go


示例执行结果：

.. literalinclude:: code/errors/example.txt


New 函数
-------------

::

    func New(text string) error

根据给定的文本返回一个错误。

示例：

.. literalinclude:: code/errors/new.go

示例执行结果：

.. literalinclude:: code/errors/new.txt

fmt 包的 Errorf 函数可以让用户使用该包的格式化功能来创建描述错误的消息。

示例：

.. literalinclude:: code/errors/errorf.go

示例执行结果：

.. literalinclude:: code/errors/errorf.txt


