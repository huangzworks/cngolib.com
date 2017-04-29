net/http/httptest
=========================

本文是 Go 标准库中 net/http/httptest 包文档的翻译，
原文地址为： 
https://golang.org/pkg/net/http/httptest/


概述
-----------

httptest 包提供了进行 HTTP 测试所需的设施。


常量
--------------

DefaultRemoteAddr 是在 ResponseRecorder 没有显式地设置 DefaultRemoteAddr 时，
RemoteAddr 返回的默认远程地址。

::

    const DefaultRemoteAddr = "1.2.3.4"


NewRequest 函数
-----------------

::

    func NewRequest(method, target string, body io.Reader) *http.Request

返回一个新的服务器访问请求，
这个请求可以传递给 http.Handler 以便进行测试。

target 参数的值为 RFC 7230 中提到的“请求目标”（request-target)：
它可以是一个路径又或者一个绝对 URL 。
如果 target 是一个绝对 URL ，
那么 URL 中的主机名（host name）将被使用；
否则主机名将为 example.com 。

当 target 的模式为 https 时，
TLS 字段的值将被设置为一个非 nil 的随意值（dummy value）。

Request.Proto 总是为 HTTP/1.1 。

如果 method 参数的值为空，
那么使用 GET 方法作为默认值。

body 参数的值可以为 nil ；
另一方面，
如果 body 参数的值为 \*bytes.Reader 类型、 \*strings.Reader 类型或者 \*bytes.Buffer 类型，
那么 Request.ContentLength 将被设置。

为了使用的方便，
NewRequest 将在 panic 可以被接受的情况下，
使用 panic 代替错误。

如果你想要生成的不是服务器访问请求，
而是一个客户端 HTTP 请求，
那么请使用 net/http 包中的 NewRequest 函数。


ResponseRecorder 类型
-------------------------

ResponseRecorder 是一个 http.ResponseWriter 实现，
它会记录自身的变化以便在之后的测试中进行检验（inspect）。

::

    type ResponseRecorder struct {
        // Code 是由 WriteHeader 设置的 HTTP 响应码。
        //
        // 注意，如果一个处理器从来没有调用过 WriteHeader 或者 Write ，
        // 那么 Code 的值将会为 0 ，而不是隐含的 http.StatusOK 。
        // 如果用户想要在 Code 值为 0 时获得隐含的 http.StatusOK ，
        // 那么可以使用 Result 方法。
        Code int

        // HeaderMap 包含了处理器显式设置的各个 HTTP 首部。
        //
        // 如果用户想要获得诸如 Content-Type 等一系列由服务器隐式地进行设置的首部，
        // 那么可以使用 Result 方法。
        HeaderMap http.Header

        // Body 是处理器调用 Write 进行写入时所使用的缓冲区。
        // 如果它的值为 nil ，那么写入的数据将会被静默地丢弃（silently discarded）。
        Body *bytes.Buffer

        // Flushed 标记处理器是否调用了 Flush 方法。
        Flushed bool

        // 其他已过滤或者未导出的字段……
    }

使用示例：

.. literalinclude:: code/httptest/response_recorder.go

示例执行结果：

::

    200
    text/html; charset=utf-8
    <html><body>Hello World!</body></html>

NewRecorder 函数
^^^^^^^^^^^^^^^^^^^^^

::

    func NewRecorder() *ResponseRecorder

返回一个初始化后的 ResponseRecorder 。

(\*ResponseRecorder) Flush 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) Flush()

将 rw.Flushed 的值设置为 true 。

(\*ResponseRecorder) Header 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) Header() http.Header

返回响应的首部。

(\*ResponseRecorder) Result 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) Result() *http.Response

返回处理器生成的响应。

处理器返回的响应至少会对状态码（StatusCode）、首部（Header）、主体（Body）以及可选的 Trailer 进行设置。
因为未来可能会有更多字段被设置，
所以用户不应该在测试里面对结果调用 DeepEqual 。

Response.Header 是写入操作第一次调用时的首部快照（snapshot of the headers）；
另一方面，
如果处理器没有执行过写入操作，
那么 Response.Header 就是 Result 方法调用时的首部快照。

Response.Body 将被生成为一个非 nil 值，
而 Body.Read 则保证不会返回除 io.EOF 之外的其他任何错误。

Result 必须在处理器执行完毕之后调用。

(\*ResponseRecorder) Write 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) Write(buf []byte) (int, error)

在 buf 不为 nil 的情况下，
Write 总会成功地将 buf 中的内容写入到 rw.Body 当中。

(\*ResponseRecorder) WriteHeader 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) WriteHeader(code int)

对 rw.Code 进行设置。
在这个方法执行之后，
修改 rw.Header 将不会对 rw.HeaderMap 产生任何影响。

(\*ResponseRecorder) WriteString 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rw *ResponseRecorder) WriteString(str string) (int, error)

在 str 不为 nil 的情况下，
WriteString 总会成功地将 buf 中的内容写入到 rw.Body 当中。


Server 类型
-------------------

Server 是一个在本地回环接口（local loopback interface）的一个由系统选择的端口上进行监听的 HTTP 服务器，
该服务器可以用于端对端 HTTP 测试。

::

    type Server struct {
        // 一个末尾不包含斜线、格式为 http://ipaddr:port 的基本 URL
        Listener net.Listener

        // 可选的 TLS 配置选项（configuration），它将在 TLS 启动之后被设置成一个新的配置（config）。
        // 如果用户在调用 StartTLS 之前，对一个尚未启动的服务器的 TLS 配置进行了设置，
        // 那么已设置的字段将被拷贝至新配置里面。
        TLS *tls.Config

        // Config 的值有可能在调用 NewUnstartedServer 之后或者调用
        // Start 和 StartTLS 之前发生变化。
        Config *http.Server

        // 其他已过滤或者未导出的字段……
    }

示例代码：

.. literalinclude:: code/httptest/server.go

执行结果：

::

    Hello, client


NewServer 函数
^^^^^^^^^^^^^^^^^^

::

    func NewServer(handler http.Handler) *Server

启动并返回一个新的服务器。
调用者在使用完这个服务器之后需要调用 Close 方法将其关闭。

NewTLSServer 函数
^^^^^^^^^^^^^^^^^^^^^^

::

    func NewTLSServer(handler http.Handler) *Server

启动并返回一个使用 TLS 的新服务器。
调用者在使用完这个服务器之后需要调用 Close 方法将其关闭。

NewUnstartedServer 函数
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func NewUnstartedServer(handler http.Handler) *Server

返回一个未启动的新服务器，
用户应该在修改完服务器的配置选项之后，
调用 Start 方法或者 StartTLS 方法来启动它，
并在使用完这个服务器之后，
调用 Close 方法将其关闭。

(\*Server) Close 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Server) Close()

关闭服务器，
并且阻塞直到对该服务器的外部请求全部完成为止。

(\*Server) CloseClientConnections 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Server) CloseClientConnections()

关闭所有连接至测试服务器的 HTTP 连接。

(\*Server) Start 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Server) Start()

启动 NewUnstartedServer 返回的服务器。

(\*Server) StartTLS 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Server) StartTLS()

启动 NewUnstartedServer 返回的服务器上的 TLS 功能。
