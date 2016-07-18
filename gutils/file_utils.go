package gutils

import (
	"fmt"
	"os"
	"bufio"
	"errors"
	"io"
	"path/filepath"
	"strings"
)




//遍历文件夹  查找文件
func ListPath(dir, fileName string) string {

	targetPath := ""
	err_ := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if ( f == nil ) {
			return err
		}
		if f.IsDir() {
			//fmt.Println("目录:" + f.Name())
		}
		if strings.HasSuffix(path, fileName) {
			targetPath = path
			//println(path)
			return errors.New("FOUND")
		}
		return nil
	})
	if err_ != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err_)
	}
	return targetPath
}

type FileInfo struct {
	FileNames []string
	fileSize  int64
}



//查找文件重复文件
func FoundRepeatFile(dir string) map[string]FileInfo {
	//定义一个字典存放所有文件信息
	allFile := make(map[string]FileInfo)
	//定义一个字典存放重名文件
	repeatFile := make(map[string]FileInfo)
	//遍历目标路径
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() && ".DS_Store" != f.Name() {
			//获得文件名
			fileName := f.Name()
			//判断文件名是否已经存在
			if v, ok := allFile[fileName]; ok {
				if v1, ok1 := repeatFile[fileName]; ok1 {
					//添加重复的文件到 重复字典里
					v1.FileNames = append(v1.FileNames, path)
				} else {
					//讲重名的2个文件 存放到 重复的字典里
					v.FileNames = append(v.FileNames, path)
					repeatFile[fileName] = v
				}
			}
			fp := []string{path}
			allFile[f.Name()] = FileInfo{fp, f.Size()}
		}
		fmt.Println(path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return repeatFile
}



//获得文件名,忽略掉路径,把路径按"/"切割,取最后一截
func GetFileName(path string) string {
	name := strings.Split(path, "/")
	//fmt.Printf("%v", name)
	return name[len(name) - 1]
}





//写字符串到文件里,拼到文件的末尾
func WriteTextToFile(path, text string) error {
	//获得一个文件指针
	//os.O_APPEND:是否拼到文件末尾
	file, err := os.OpenFile(path, os.O_WRONLY | os.O_APPEND, 0666)
	if err != nil {
		//抛出一个异常
		panic(errors.New("打开文件错误!"))
	}
	// 执行完函数后,关闭文件
	defer func() {
		file.Close()
	}()
	//获得一个文件写入器
	writer := bufio.NewWriter(file)
	//写入text
	writer.WriteString(text)
	return writer.Flush()
}


//读取文件内所有字符串
func ReadFileString(path string) error {

	//获得一个文件指针
	file, erro := os.Open(path)
	if erro != nil {
		//抛出一个异常
		panic(errors.New("打开文件错误!"))
	}
	// 执行完函数后,关闭文件
	defer func() {
		file.Close()
	}()
	//获得一个Reader指针
	reader := bufio.NewReader(file)
	var str string
	var err error
	//如果文件已经读完, 会返回一个erro: EOF
	for err == nil {
		//循环读取文件内的字符串
		str, err = reader.ReadString('\n');
		fmt.Printf(str)
	}
	return err
}



//读取文件内所有字符串2,通过切片实现
func ReadFileString2(path string) {

	//获得一个文件指针
	file, erro := os.Open(path)
	if erro != nil {
		//抛出一个异常
		panic(errors.New("打开文件错误!"))
	}
	// 执行完函数后,关闭文件
	defer func() {
		file.Close()
	}()
	var buf [1021]byte
	for {
		switch num, err := file.Read(buf[:]); true {
		case num > 0:
			n, e := os.Stdout.Write(buf[:])
			//如果读出来的数量和写进去的数量不同
			if ( num != n) {
				panic(e)
			}
		case num == 0:
			break
		case num < 0:
			//如果读到的数量为负
			panic(err)
		}
	}
}


//复制文件
func CopyFile(from, to string) (written int64, err error) {
	//获得一个文件指针
	file_from, err := os.Open(from)
	if err != nil {
		//抛出一个异常
		panic(errors.New("打开文件错误!"))
	}
	// 执行完函数后,关闭文件
	defer func() {
		fmt.Println("defer:关闭文件!")
		file_from.Close()
	}()
	//获得一个文件指针
	file_to, err := os.OpenFile(to, os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil {
		//抛出一个异常
		panic(errors.New("打开目标文件错误!"))
	}
	// 执行完函数后,关闭文件
	defer func() {
		fmt.Println("defer:关闭目标文件!")
		file_to.Close()
	}()
	return io.Copy(file_to, file_from)

}
