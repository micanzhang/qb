# qb [![Build Status](https://travis-ci.org/micanzhang/qb.svg?branch=master)](https://travis-ci.org/micanzhang/qb) [![Go Report Card](https://goreportcard.com/badge/github.com/micanzhang/qb)](https://goreportcard.com/report/github.com/micanzhang/qb)
qb是一个备份文件到七牛云存储的命令行工具.

## 依赖:

1. go 1.7 

## 安装

``` sh 
go get -u github.com/micanzhang/qb
```

## 使用

```sh 
qb --help 
```
### 文件上传

```sh
$qb put file1 file2 ...filen -name name1,name2,...namen
```

name默认为文件名，可以通过 `-name` 参数自定义文件名, 多个文件以 `,` 分隔.

### 文件下载

```sh
qb get name1  -dir ~/Downloads
```

`-dir`: 指定文件下载目录，默认为当前文件夹.

### 获取文件信息

```sh
$qb info name1,name2....namen
```

### 文件删除

```sh
$qb remove name1,name2....namen
```

## TODO
1. 实现类似linux ls命令.
2. 文件夹文件自动同步.


## License 
MIT License

Copyright (c) 2016 micanzhang

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.


