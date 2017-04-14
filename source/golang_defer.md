title: "golang defer"
date: 2017-04-13 23:16:00 +0800
update: 2017-04-13 23:16:00 +0800
author: mutemaniac
#cover: "-/images/example.png"
tags:
    - golang
    - defer
hide: false #隐藏文章，只可通过链接访问，可选

---

go语言defer声明调用的函数会在函数结束后才会执行。对于回收资源和释放锁这种容易忘记但又必须做的操作，defer很方便，又不容易忘记。最近使用遇见一些注意项mark一下：  
###  LIFO 
栈是的后进先出，对于有顺序的资源释放要注意。
### 先判断返回错误，在defer释放资源  
对于有些资源的申请，一般会返回错误，在用defer释放资源是一定要先判断是否有错误产生。如果有错误产生这时候是不用释放资源的，因为资源就没有分配，我就遇到过一次发送网络请求后，没有判断错误就```defer resp.Body.Close()```,导致在网络请求失败时释放资源出错。
### 一个defer做一件事  
一个defer做一件事，当有多个资源需要释放是，放到不同的defer里，因为万一defer中有代码执行出错，后面的资源就不会释放了。直接看例子：
```go  
package main

import "fmt"

func main() {
    func() {
        defer func() {
            fmt.Println("1")
            fmt.Println("2")
            panic("after 2")
            fmt.Println("3")
        }()
    }()
    fmt.Println("Hello, 世界")
}
```
运行结果:
```
1
2
panic: after 2

goroutine 1 [running]:
main.main.func1.1()
    /tmp/sandbox790552561/main.go:10 +0x180
main.main.func1()
    /tmp/sandbox790552561/main.go:13 +0x7b
main.main()
    /tmp/sandbox790552561/main.go:13 +0x20
```
以上程序的假设打印2后程序错误，会导致3没有打印出来，如果把每件事分开：
```go
package main

import "fmt"

func main() {
    func() {
       defer func() {
           fmt.Println("1")
        }()
        defer func() {
            fmt.Println("2")
            panic("after 2")
        }()
        defer func() {
            fmt.Println("3")
        }()
    }()
    fmt.Println("Hello, 世界")
}
```
运行结果
```
3
2
1
panic: after 2

goroutine 1 [running]:
main.main.func1.2()
    /tmp/sandbox153708588/main.go:14 +0x100
main.main.func1()
    /tmp/sandbox153708588/main.go:20 +0xbb
main.main()
    /tmp/sandbox153708588/main.go:20 +0x20
```
这样的话123就会全部执行了。

### defer调用函数时，参数在声明的一刻就已经确定了
defer调用函数时，参数在声明的一刻就已经确定了，即使参数变量在以后会有变化，也不会影响defer调用。（个人猜想：defer声明时会把所需参数变量当前值入栈，函数结束后变量已经不再存在），
```go
package main

import "fmt"

func main() {
    func() {
        i := 0
        defer func(j int) {
            fmt.Println("i = ", i)
            fmt.Println("j = ", j)
        }(i + 1)
        i = 3
    }()
    fmt.Println("Hello, 世界")
}
```
result:
```
i =  3
j =  1
Hello, 世界
```
```i+1```表达式在defer声明的一刻就已经计算出来了，不会再受后面程序的影响 （***注：闭包的特性没有受影响***）

还有一个golang官方的例子，贴出来：
```go
func trace(s string) string {
    fmt.Println("entering:", s)
    return s
}

func un(s string) {
    fmt.Println("leaving:", s)
}

func a() {
    defer un(trace("a"))
    fmt.Println("in a")
}

func b() {
    defer un(trace("b"))
    fmt.Println("in b")
    a()
}

func main() {
    b()
}
```
prints
```
entering: b
in b
entering: a
in a
leaving: a
leaving: b
```