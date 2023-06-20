package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// Dir 日志文件存放目录
var Dir = "../../tmp/"

// 获取第几层函数的名称
// 1-当前层
// 2-上一层
// 3-再上一层
// 4-再再上一层
func getFunName() string {
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	split := strings.Split(name, ".")
	//fmt.Printf("第%d层函数,函数名称是:%s\n", l, name)
	return strings.ToLower(split[len(split)-1])
}

func Error(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println(errType, msg, GetStacktrace(5))
}

func Info(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println(errType, msg, GetStacktrace(5))
}

func Println(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println(errType, msg, GetStacktrace(5))
}

func Debug(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println(errType, msg, GetStacktrace(5))
}

func Fatalln(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Fatalln(errType, msg, string(debug.Stack()))
}

func Fatal(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Fatalln(errType, msg, string(debug.Stack()))
}

func Fatalf(msg ...any) {
	errType := getFunName()
	filename := Dir + time.Now().Format("2006-01-02_") + errType + ".log"
	err := os.MkdirAll(Dir, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	// 组合一下即可，os.Stdout代表标准输出流
	multiWriter := io.MultiWriter(os.Stdout, f)
	log.SetOutput(multiWriter)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Fatalln(errType, msg, string(debug.Stack()))
}

// GetStacktrace 获取当前 goroutine 的调用栈信息
// skip 参数用于控制要跳过的栈帧数，一般传入1即可，表示跳过当前函数的栈帧
func GetStacktrace(skip int) string {
	return strings.Join(strings.Split(string(debug.Stack()), "\n")[skip+1:], "\n")
}
