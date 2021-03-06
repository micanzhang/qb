# qb [![Build Status](https://travis-ci.org/micanzhang/qb.svg?branch=master)](https://travis-ci.org/micanzhang/qb) [![Go Report Card](https://goreportcard.com/badge/github.com/micanzhang/qb)](https://goreportcard.com/report/github.com/micanzhang/qb)
cli tool for backup files to qiniu storage system.

[chinese version](README_CN.md)

## Dependency:

1. go 1.7 

## Install 

``` sh 
go get -u github.com/micanzhang/qb
```

## basic usages

```sh 
qb --help 
```
### put files

```sh
$qb put file1 file2 ...filen -name name1,name2,...namen
```

default name is file's own name;

use can customrize your file's name by pass `-name` parameters, split by `,` for multiple files.

### get file 

```sh
qb get name1  -dir ~/Downloads
```

`-dir`: specific directory where files downloaded into, and default dir is directory where cmd your are run at.

### info files

```sh
$qb info name1,name2....namen
```

### remove files

```sh
$qb remove name1,name2....namen
```

## TODO 
1. refact list command like linux command ls does
2. auto sync specfic directory


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


