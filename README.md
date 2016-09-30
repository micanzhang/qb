# qb
cli tool for backup files  to qiniu storage system.


## basic usages

### put files

```sh
$qb put file1 file2 ...filen -name name1,name2,...namen
```

default name is file's own name;

use can customrize your file's name by pass `-name` parameters, split by `,` for multiple files;

### info files

```sh
$qb info name1,name2....namen
```

### remove files

```sh
$qb remove name1,name2....namen
```

