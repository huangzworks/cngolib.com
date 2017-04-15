database/sql
================

本文是 Go 标准库中 database/sql 包文档的翻译，
原文地址为： 
https://golang.org/pkg/database/sql/


概述
-----------------

sql 包为操作 SQL 以及类 SQL 数据库提供了一个通用的接口。

sql 包必须与数据库驱动一同使用，
https://golang.org/s/sqldrivers 列出了可用的数据库驱动列表。

不支持上下文取消操作（context cancelation）的驱动会在查询完成之后才返回。

更多 sql 包使用示例请见 wiki 页面：
https://golang.org/s/sqlwiki


变量
----------------

当 QueryRow 没有返回一个行时，
Scan 将返回 ErrNoRows 。
在这种情况下，
QueryRow 将返回一个 \*Row 占位符值，
从而将这个错误推延至 Scan 。

::

    var ErrNoRows = errors.New("sql: no rows in result set")

当用户尝试在事务中执行一个已经提交又或者已经回滚的操作时，
就会返回一个 ErrTxDone ：

::

    var ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")


Drivers 函数
--------------

::

    func Drivers() []string

以有序列表的形式返回所有已注册驱动的名字。



Register 函数
----------------

::

    func Register(name string, driver driver.Driver)

将带有给定名字的驱动设置为可用。
如果用户以相同的名字重复调用 Register ，
又或者 driver 参数的值为 nil ，
那么 Register 函数将引发一个 panic 错误。



ColumnType 类型
------------------

ColumnType 包含了数据列的名字以及类型：

::

    type ColumnType struct {
                // contains filtered or unexported fields
    }


(\*ColumnType) DatabaseTypeName 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) DatabaseTypeName() string

返回数据列类型的数据库系统名字。

返回空字符表示该驱动类型名字并未被支持。

驱动的数据类型可以通过查看驱动的文档来得到。

长度指示器不会被包含在内（Length specifiers are not included）。

通用的类型包括 "VARCHAR"， "TEXT"， "NVARCHAR"， "DECIMAL"， "BOOL"， "INT"， "BIGINT" 。

(\*ColumnType) DecimalSize 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) DecimalSize() (precision, scale int64, ok bool)

返回小数类型的范围和精度。

如果给定的类型不可用又或者不支持，
那么将返回值 ok 参数的值设置为 false 。

(\*ColumnType) Length 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) Length() (length int64, ok bool)

对于诸如 text 和 binary 字段类型这种可变长度的数据列类型，
返回它们的数据列类型长度。

如果该类型的长度可以是无限大的，
那么值将被设置为 math.MaxInt64 ，
但具体的长度仍然由数据库本身决定。

如果给定的类型并不是可变长度类型，
比如 int ，
又或者驱动并不支持该类型，
那么返回值 ok 参数的值将被设置为 false 。

(\*ColumnType) Name 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) Name() string

返回数据列的名字又或者别名。

(\*ColumnType) Nullable 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) Nullable() (nullable, ok bool)

检查给定的数据列能否为 null 。

在驱动不支持该属性的情况下，
ok 参数的值为 false 。

(\*ColumnType) ScanType 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ci *ColumnType) ScanType() reflect.Type

返回一种 Go 类型，
该类型能够在 Rows.Scan 进行扫描时使用。

如果驱动不支持该属性，
那么返回空接口类型。



DB 类型
-----------------

DB 是一个数据库句柄，
它代表的是包含着零个或多个底层连接的池（pool）。
多个 goroutine 可以安全地、并发地使用这个句柄。

sql 包会自动地创建和释放连接，
并且它也会维护一个由闲置（idle）连接组成的释放池（free pool）。
如果数据库拥有预连接状态（pre-connection state）这一概念，
那么这种状态只会在事务内部可见（observed）。

当 DB.Begin 被调用时，
它返回的 Tx 将与单个连接绑定，
而当事务提交或者回滚时，
这个连接将返回至 DB 的闲置连接池。
闲置连接池的大小可以通过 SetMaxIdleConns 来控制。

::

    type DB struct {
        // contains filtered or unexported fields
    }

Open 函数
^^^^^^^^^^^^^^

::

    func Open(driverName, dataSourceName string) (*DB, error)

Open 函数会根据给定的数据库驱动以及驱动专属的数据源来打开一个数据库，
驱动专属的数据源一般至少会包含数据库的名字以及相关的连接信息。

大多数用户都会通过驱动专属的辅助函数来打开数据库，
这种函数会返回一个指向 DB 结构的指针。
Go 的标准库不包含任何数据库驱动，
所有数据库驱动都是第三方的，
https://golang.org/s/sqldrivers 列出了可用的第三方驱动。

Open 有可能会只对参数进行检查，
但是却并不创建实际的数据库连接。
通过调用 Ping 可以检查给定的数据源是否有效。

Open 函数返回的 DB 可以安全地由多个 goroutine 进行并发使用，
并且 DB 也会维护它自有的闲置连接池。
因此，
Oepn 函数通常只需要调用一次，
并且用户通常不需要手动地关闭一个 DB 。

(\*DB) Begin 方法
^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Begin() (*Tx, error)

开启一个事务，
事务的隔离级别由驱动决定。

(\*DB) BeginTx 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error)

开启一个事务。

给定的上下文会一直使用到事务提交又或者回滚为止。
如果上下文被取消了，
那么 sql 包将会对事务进行回滚。
Tx.Commit 在给定的上下文已被取消时会返回一个错误。

TxOptions 参数是可选的，
传入 nil 则表示使用默认值。
如果用户尝试使用一种驱动不支持的非默认隔离级别，
那么方法将返回一个错误。

(\*DB) Close 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Close() error

关闭数据库并释放所有已打开的资源。

因为 DB 句柄通常会长时间存在，
并且会有多个 goroutine 进行分享，
所以用户一般不需要手动地关闭数据库。

(\*DB) Driver 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Driver() driver.Driver

返回数据库的底层驱动。

(\*DB) Exec 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Exec(query string, args ...interface{}) (Result, error)

执行指定的查询，
但不返回任何数据行。
方法的 arg 部分用于填写查询语句中包含的占位符的实际参数。

(\*DB) ExecContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)

同上。

(\*DB) Ping 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Ping() error

检查数据库连接是否仍然有效，
并在有需要时建立一个连接。

(\*DB) PingContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) PingContext(ctx context.Context) error

同上。

(\*DB) Prepare 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Prepare(query string) (*Stmt, error)

为之后的查询或执行（execution）创建预处理语句，
多个查询或者执行可以并发地使用 Prepare 返回的预处理语句。
当调用者不再需要这个预处理语句时，
它必须调用这个语句的 Close 方法。

(\*DB) PrepareContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error)

为之后的查询或执行（execution）创建预处理语句，
多个查询或者执行可以并发地使用 Prepare 返回的预处理语句。
当调用者不再需要这个预处理语句时，
它必须调用这个语句的 Close 方法。

给定的上下文将在创建预处理语句时使用，
而不是在执行该语句时使用。

(\*DB) Query 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Query(query string, args ...interface{}) (*Rows, error)

执行一个查询并返回多个数据行，
这个查询通常是一个 SELECT 。
方法的 arg 部分用于填写查询语句中包含的占位符的实际参数。

一个单结果示例：

::

    age := 27
    rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var name string
        if err := rows.Scan(&name); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s is %d\n", name, age)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

一个多结果示例：

::

    age := 27
    q := `
    create temp table uid (id bigint); -- Create temp table for queries.
    insert into uid
    select id from users where age < ?; -- Populate temp table.

    -- First result set.
    select
        users.id, name
    from
        users
        join uid on users.id = uid.id
    ;

    -- Second result set.
    select 
        ur.user, ur.role
    from
        user_roles as ur
        join uid on uid.id = ur.user
    ;
    `
    rows, err := db.Query(q, age)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    for rows.Next() {
        var (
            id   int64
            name string
        )
        if err := rows.Scan(&id, &name); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("id %d name is %s\n", id, name)
    }
    if !rows.NextResultSet() {
        log.Fatal("expected more result sets", rows.Err())
    }
    var roleMap = map[int64]string{
        1: "user",
        2: "admin",
        3: "gopher",
    }
    for rows.Next() {
        var (
            id   int64
            role int64
        )
        if err := rows.Scan(&id, &role); err != nil {
            log.Fatal(err)
        }
        fmt.Printf("id %d has role %s\n", id, roleMap[role])
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

(\*DB) QueryContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)

同上。

(\*DB) QueryRow 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) QueryRow(query string, args ...interface{}) *Row

执行一个预期最多只会返回一个数据行的查询。

这个方法总是会返回一个非空的值，
而它引起的错误则会被推延到数据行的 Scan 方法被调用为止。

示例：

::

    id := 123
    var username string
    err := db.QueryRow("SELECT username FROM users WHERE id=?", id).Scan(&username)
    switch {
    case err == sql.ErrNoRows:
        log.Printf("No user with that ID.")
    case err != nil:
        log.Fatal(err)
    default:
        fmt.Printf("Username is %s\n", username)
    }

(\*DB) QueryRowContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row

同上。

(\*DB) SetConnMaxLifetime 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) SetConnMaxLifetime(d time.Duration)

设置可以重用连接的时长。

过期的连接可以在重用之前惰性地进行关闭。

如果 d <= 0 ，
那么连接将一直可用。

(\*DB) SetMaxIdleConns 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) SetMaxIdleConns(n int)

设置闲置连接池里面最多可放置的连接数量。

如果 MaxOpenConns 大于 0 但小于新的 MaxIdleConns ，
那么将 MaxIdleConns 的值设置为与 MaxOpenConns 一样。

如果 n <= 0 ，
那么不存放任何闲置的连接。

(\*DB) SetMaxOpenConns 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) SetMaxOpenConns(n int)

设置最大可创建的数据库连接数量。

如果 MaxIdleConns 大于 0 并且新的 MaxOpenConns 小于 MaxIdleConns ，
那么 MaxIdleConns 的值将设置为与 MaxOpenConns 一样。

如果 n <= 0 ，
那么表示不限制数据库连接的数量。
默认值为 0 （无限制）。

(\*DB) Stats 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (db *DB) Stats() DBStats

返回数据库的统计数据。



DBStats 类型
-----------------

DBStats 包含了数据库的统计数据。

::

    type DBStats struct {
        // OpenConnections is the number of open connections to the database.
        OpenConnections int
    }



IsolationLevel 类型
------------------------

IsolationLevel 是用于 TxOptions 的事务隔离级别：

::

    type IsolationLevel int

不同驱动在 BeginTx 中对隔离级别的支持也是不同的，
如果一个驱动不支持给定的隔离级别，
那么 BeginTx 将返回一个错误。

::

    const (
        LevelDefault IsolationLevel = iota
        LevelReadUncommitted
        LevelReadCommitted
        LevelWriteCommitted
        LevelRepeatableRead
        LevelSnapshot
        LevelSerializable
        LevelLinearizable
    )

关于隔离级别的更多信息请查看：https://en.wikipedia.org/wiki/Isolation_(database_systems)#Isolation_levels



NamedArg 类型
-------------------

一个 NamedArg 就是一个具名参数。
NamedArg 的值可以用作 Query 或者 Exec 的参数，
并与 SQL 语句中相应的具名参数进行绑定。

通过 Named 函数可以更方便地创建 NamedArg 值。

::

    type NamedArg struct {
        // 参数占位符的名字。
        // 如为空，那么根据参数列表中的排列位置进行设置。
        // Name 必须省略所有符号前缀。
        Name string

        // 参数的值。
        // 这个参数可能会被设置为与查询参数具有相同的值类型。
        Value interface{}

        // 其他已过滤字段以及未导出字段
    }

Named 函数
^^^^^^^^^^^^^

::

    func Named(name string, value interface{}) NamedArg

Named 提供了一种更为方便的创建 NamedArg 值的方法。

以下是一个使用示例：

::

    db.ExecContext(ctx, `
        delete from Invoice
        where
            TimeCreated < @end
            and TimeCreated >= @start;`,
        sql.Named("start", startTime),
        sql.Named("end", endTime),
    )



NullBool 类型
-------------------

NullBool 表示一个可能为 null 的布尔值。
NullBool 实现了 Scanner 接口，
因此它可以跟 NullString 一样用作扫描目的地（destination）：

::

    type NullBool struct {
        Bool  bool
        Valid bool // 当 Bool 字段的值不为 NULL 时， Valid 字段的值为 true
    }

(\*NullBool) Scan 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n *NullBool) Scan(value interface{}) error

Scan 实现了 Scanner 接口。

(NullBool) Value 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n NullBool) Value() (driver.Value, error)

Value 实现了驱动的 Valuer 接口。



NullFloat64 类型
---------------------

NullFloat64 用于表示一个可能为 null 的 float64 值。

NullFloat64 实现了 Scanner 接口，
因此它可以跟 NullString 一样用作扫描目的地：

::

    type NullFloat64 struct {
        Float64 float64
        Valid   bool // 当 Float64 不为 NULL 时，Valid 为 true
    }

(\*NullFloat64) Scan 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n *NullFloat64) Scan(value interface{}) error

Scan 实现了 Scanner 接口。

(NullFloat64) Value 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n NullFloat64) Value() (driver.Value, error)

Value 实现了驱动的 Valuer 接口。



NullInt64 类型
----------------------

NullInt64 用于表示一个可能为 null 的 int64 值。
NullInt64 实现了 Scanner 接口，
因此它可以跟 NullString 一样用作扫描目的地：

::

    type NullInt64 struct {
        Int64 int64
        Valid bool // 当 Int64 不为 NULL 时，Valid 为 true
    }

(\*NullInt64) Scan 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n *NullInt64) Scan(value interface{}) error

Scan 实现了 Scanner 接口。

(NullInt64) Value 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (n NullInt64) Value() (driver.Value, error)

Value 实现了驱动的 Valuer 接口。



NullString 类型
---------------------

NullString 表示一个可能为 null 的字符串：

::

    type NullString struct {
        String string
        Valid  bool // 当 String 不为 NULL 时，Valid 为 true
    }

NullString 实现了 Scanner 接口，
因此它可以用作扫描目的地：

::

    var s NullString
    err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
    ...
    if s.Valid {
           // use s.String
    } else {
           // NULL value
    }

(\*NullString) Scan 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ns *NullString) Scan(value interface{}) error

Scan 实现了 Scanner 接口。

(NullString) Value 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (ns NullString) Value() (driver.Value, error)

Value 实现了驱动的 Valuer 接口。



RawBytes 类型
-----------------

RawBytes 是一个字节串，
它持有一个引用，
该引用指向数据库自身拥有的内存。

在 Scan 将结果储存到一个 RawBytes 之后，
该切片会在下一个 Next 、Scan 或者 Close 调用之前一直有效。

::

    type RawBytes []byte



Result 类型
--------------

Result 是对已执行 SQL 命令的总结。

::

    type Result interface {

        LastInsertId() (int64, error)

        RowsAffected() (int64, error)
    }

LastInsertId() 会返回一个由数据库生成的整数，
这个整数是对命令的响应。
在插入一个新的数据行时，
这个整数通常来源于数据表中的自增数据列。
并不是所有数据库都支持这个特性，
并且各个数据库在实现这个特性时使用的语句也会有所不同。

RowsAffected() 返回受到更新、插入或者删除操作影响的行数量，
并不是所有数据库或者所有数据库驱动都支持这个特性。



Row 类型
-----------------

Row 是调用 QueryRow 查询单个数据行所得的结果。

::

    type Row struct {
        // contains filtered or unexported fields
    }

(\*Row) Scan 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (r *Row) Scan(dest ...interface{}) error

Scan 会将被匹配数据行中的各个列复制到 dest 指向的值里面，
更多细节请参考 Rows.Scan 方法的文档。
如果有多个数据行与查询匹配，
那么 Scan 将使用第一个数据行并丢弃其他所有数据行。
如果没有任何数据行与查询匹配，
那么 Scan 将返回 ErrNoRows 。



Rows 类型
-----------------

Rows 是查询的结果：

::

    type Rows struct {
        // contains filtered or unexported fields
    }

它的游标会从结果集的第一个数据行开始，
而用户则可以通过 Next 来遍历结果集中的所有数据行：

::

    rows, err := db.Query("SELECT ...")
    ...
    defer rows.Close()
    for rows.Next() {
        var id int
        var name string
        err = rows.Scan(&id, &name)
        ...
    }
    err = rows.Err() // get any error encountered during iteration
    ...

(\*Rows) Close 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rs *Rows) Close() error

关闭 Rows ，
停止对数据集的迭代。
如果 Next 返回 false ，
那么 Rows 将自动关闭，
并且用户在自动关闭的情况下也同样能够检查 Err 的结果。

(\*Rows) Columns 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rs *Rows) Columns() ([]string, error)

返回各个列的名字。

当数据行已经被关闭时，
Columns 将返回一个错误；
当数据行来自于 QueryRow 时，
Columns 将返回一个推延错误。

(\*Rows) Err 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (rs *Rows) Err() error

返回迭代过程中出现的任何错误，
这个方法在数据行显式或者隐式关闭之后仍然可用。

(\*Rows) Next 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (rs *Rows) Next() bool

Next 可以为 Scan 方法准备好下一个待读取的数据行：
这个方法在执行成功时返回 true ；
返回 false 表示没有下一个数据行可用，
又或者准备期间发生了错误。

通过 Err 方法可以知道 Next 是因为何种原因而执行失败的。

(\*Rows) Scan 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (rs *Rows) Scan(dest ...interface{}) error

将当前被迭代数据行的各个列复制到 dest 指向的值里面，
dest 的值数量必须与数据行中的列数量保持一致。

Scan 会把从数据库里面读取到的各个数据列转换为以下标准 Go 类型，
又或者转换为某些由 sql 包提供的特殊类型：

- \*string
- \*[]byte
- \*int, \*int8, \*int16, \*int32, \*int64
- \*uint, \*uint8, \*uint16, \*uint32, \*uint64
- \*bool
- \*float32, \*float64
- \*interface{}
- \*RawBytes
- 实现了 Scanner 接口的任何类型（具体信息请见 Scanner 接口的文档）

..
    在最简单的情况下，
    如果数据列的值是一个类型为 T 的整数、布尔值或者字符串，
    并且 dest 的类型为 \*T ，
    那么 Scan 只需要将数据列的值赋值给这些指针就可以了。

    In the most simple case, if the type of the value from the source column is an integer, bool or string type T and dest is of type *T, Scan simply assigns the value through the pointer.

    Scan also converts between string and numeric types, as long as no information would be lost. While Scan stringifies all numbers scanned from numeric database columns into *string, scans into numeric types are checked for overflow. For example, a float64 with value 300 or a string with value "300" can scan into a uint16, but not into a uint8, though float64(255) or "255" can scan into a uint8. One exception is that scans of some float64 numbers to strings may lose information when stringifying. In general, scan floating point columns into *float64.

    If a dest argument has type *[]byte, Scan saves in that argument a copy of the corresponding data. The copy is owned by the caller and can be modified and held indefinitely. The copy can be avoided by using an argument of type *RawBytes instead; see the documentation for RawBytes for restrictions on its use.

    If an argument has type *interface{}, Scan copies the value provided by the underlying driver without conversion. When scanning from a source value of type []byte to *interface{}, a copy of the slice is made and the caller owns the result.

    Source values of type time.Time may be scanned into values of type *time.Time, *interface{}, *string, or *[]byte. When converting to the latter two, time.Format3339Nano is used.

    Source values of type bool may be scanned into types *bool, *interface{}, *string, *[]byte, or *RawBytes.

    For scanning into *bool, the source may be true, false, 1, 0, or string inputs parseable by strconv.ParseBool.

.. **



Scanner 类型
---------------

Scanner 是 Scan 使用的一个接口：

::

    type Scanner interface {
        // Scan 会通过数据库驱动获取一个值，
        // 这个值将会是以下类型之一：
        //
        //    int64
        //    float64
        //    bool
        //    []byte
        //    string
        //    time.Time
        //    nil - 用于表示 NULL 值
        //
        // 当一个值无法以不丢失任何信息的情况下储存时，
        // 返回一个错误
        Scan(src interface{}) error
    }



Stmt 类型
--------------

Stmt 用于代表预处理语句，
多个 goroutine 可以安全地以并发的形式使用这种类型。

::

    type Stmt struct {
        // contains filtered or unexported fields
    }

(\*Stmt) Close 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) Close() error

关闭给定的预处理语句。

(\*Stmt) Exec 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) Exec(args ...interface{}) (Result, error)

使用给定的参数执行预处理语句，
并返回一个 Result 值来总结本次执行产生的影响。

(\*Stmt) ExecContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) ExecContext(ctx context.Context, args ...interface{}) (Result, error)

同上。

(\*Stmt) Query 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) Query(args ...interface{}) (*Rows, error)

使用给定的参数执行预处理语句，
并以 \*Rows 形式返回查询结果。

(\*Stmt) QueryContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) QueryContext(ctx context.Context, args ...interface{}) (*Rows, error)

同上。

(\*Stmt) QueryRow 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) QueryRow(args ...interface{}) *Row

使用给定的参数执行预处理语句，
并返回一个 \*Row 作为结果。
如果语句在执行期间出现了错误，
那么这个错误将会在用户对 \*Row 执行 Scan 调用时返回。
因为 Scan 调用总是返回一个非空值，
所以当查询没有获取到任何数据行时，
\*Row 的 Scan 调用将返回 ErrNoRows ；
另一方面，
在正常情况下，
\*Row 的 Scan 调用将返回查询结果中的第一个数据行，
并丢弃可能存在的所有剩余数据行。

使用示例：

::

    var name string
    err := nameByUseridStmt.QueryRow(id).Scan(&name)

(\*Stmt) QueryRowContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (s *Stmt) QueryRowContext(ctx context.Context, args ...interface{}) *Row

与 QueryRow 方法作用相同，
只是参数多了个上下文。

使用示例：

::

    var name string
    err := nameByUseridStmt.QueryRowContext(ctx, id).Scan(&name)



Tx 类型
-------------

Tx 是一个进行中的数据库事务。

一个事务必须以调用 Commit 或者 Rollback 为结束，
在调用这两个方法之中的任何一个之后，
事务的所有操作都会以 ErrTxDone 的方式失效（fail）。

通过 Prepare 方法或者 Stmt 方法放入到事务里面的语句，
将在 Commit 调用或者 Rollback 调用之后关闭。

::

    type Tx struct {
            // contains filtered or unexported fields
    }

(\*Tx) Commit 方法
^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Commit() error

提交事务。

(\*Tx) Exec 方法
^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Exec(query string, args ...interface{}) (Result, error)

执行一个不返回数据行的查询，
比如一个 INSERT 或者一个 UPDATE 。

(\* ExecContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)

同上。

(\*Tx) Prepare 方法
^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Prepare(query string) (*Stmt, error)

创建一个可以在事务里面使用的预备语句。

这个方法返回的语句将在事务中执行，
并且它在事务提交或者回滚之后将不再可用。

如果你想要在事务里面使用一个已经存在的预备语句，
那么请使用 Tx.Stmt 方法。

(\*Tx) PrepareContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) PrepareContext(ctx context.Context, query string) (*Stmt, error)

作用同上。

给定的上下文将用于预备阶段，
而不是事务的执行阶段。

这个方法返回的语句将在事务上下文中执行。

(\*Tx) Query 方法
^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error)

执行一个会返回数据行的查询，
通常是一个 SELECT 。

(\*Tx) QueryContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)

同上。

(\*Tx) QueryRow 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) QueryRow(query string, args ...interface{}) *Row

执行一个预期最多只会返回一个数据行的查询。

这个方法总是返回一个不为 nil 的值。

执行时的错误将推延到数据行的 Scan 方法执行为止。

(\*Tx) QueryRowContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
    
同上。

(\*Tx) Rollback 方法
^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Rollback() error
    
终止事务。

(\*Tx) Stmt 方法
^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) Stmt(stmt *Stmt) *Stmt
    
为已有的语句返回一个事务专用的预备语句。

示例：

::

    updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
    ...
    tx, err := db.Begin()
    ...
    res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
    
这个方法返回的语句将在事务中执行。
在事务提交或者回滚之后，
语句也会被关闭。

(\*Tx) StmtContext 方法
^^^^^^^^^^^^^^^^^^^^^^^^^^^

::

    func (tx *Tx) StmtContext(ctx context.Context, stmt *Stmt) *Stmt

同上。



TxOptions 类型
-------------------

TxOptions 用于持有在 DB.BeginTx 中使用的事务选项：

::

    type TxOptions struct {
        // Isolation 用于设置事务的隔离级别
        // 值为 0 时，使用数据库的默认级别
        Isolation IsolationLevel
        ReadOnly  bool
    }




