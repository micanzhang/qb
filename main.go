package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/howeyc/gopass"
	"github.com/micanzhang/qb/backup"
)

var (
	// default config location
	defaultConf string
	// backup action
	action string
	// customize file name
	name string
	// directory where file downloaded to
	dirpath string
	// page size of list files
	size int
	// marker for pagination
	marker string
	// search prefix
	prefix string
	// args
	args []string
	// backup provider
	provider backup.BackupProvider
)

func init() {
	flag.StringVar(&name, "name", "", "specific name for files")
	flag.StringVar(&dirpath, "dir", "./", "directory where files downloaded")
	flag.IntVar(&size, "size", 10, "page size of list files")
	flag.StringVar(&marker, "marker", "", "marker for pagination")
	flag.StringVar(&prefix, "search", "", "search file name which only supports prefix search")

	// set command usage
	flag.Usage = Usage

	// init configuration
	initConf()

	// check default folder exists or not, if not, create new one
	if !backup.FileExists(backup.Conf.Path) {
		if _, err := os.Create(backup.Conf.Path); err != nil {
			panic(err)
		}
	}

	// pg action [flags] [args]
	temp := os.Args
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(2)
	}

	action = os.Args[1]
	os.Args = os.Args[1:]
	flag.Parse()
	args = flag.Args()
	os.Args = temp
}

func main() {
	switch action {
	case "put":
		files := args
		if len(files) == 0 {
			flag.Usage()
			return
		}

		var keys []string
		if name != "" {
			keys = strings.Split(name, ",")
			if len(keys) != len(files) {
				log.Fatal("invalid length of names")
			}
		}

		for i, file := range files {
			key := ""
			if name != "" && keys[i] != "" {
				key = keys[i]
			}
			err := putFile(file, key)
			if err != nil {
				// handle error
				log.Printf("put file %s, key %s failed: %s", file, key, err)
			}

		}
	case "info":
		files := args
		if len(files) == 0 {
			flag.Usage()
			return
		}

		for _, file := range files {
			if entry, err := provider.Info(file); err != nil {
				panic(err)
			} else {
				fmt.Printf("info : %+v", entry)
			}
		}
	case "get":
		if len(args) == 0 {
			flag.Usage()
			return
		}

		if _, err := os.Stat(dirpath); err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("directory of %s not exists", dirpath)
			}

			if os.IsPermission(err) {
				log.Fatalf("permission denied")
			}

			log.Fatalf("%s", err)
		}

		for _, name := range args {
			if err := provider.Get(name, dirpath); err != nil {
				log.Printf("get file %s failed: %s", name, err)
			}
		}
	case "remove":
		files := args
		if len(files) == 0 {
			flag.Usage()
			return
		}

		for _, file := range files {
			if err := provider.Remove(file); err != nil {
				panic(err)
			}
		}
	case "list":
		if err := provider.List(prefix, marker, size); err != nil {
			panic(err)
		}
	default:
		flag.Usage()
	}

}

func Usage() {
	fmt.Fprintf(os.Stdout, "%s is a cli tools for files backup to cloud storage services, like qiniu.\n", os.Args[0])
	fmt.Fprintf(os.Stdout, "\nUsage:\n\n")
	fmt.Fprintf(os.Stdout, "\t%s command [flags] [arguments]\n", os.Args[0])
	fmt.Fprintf(os.Stdout, `
the commands are:

    put        put files to cloud
    get        get files from cloud
    info       get files's info
    remove     remove files
    list       list files

`)
	fmt.Fprintf(os.Stdout, "the flags are:\n\n")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stdout)

}

func putFile(filepath string, key string) (err error) {
	if key == "" {
		key, err = backup.FileKey(filepath)
		if err != nil {
			return err
		}
	}

	// check out exists or not
	entry, err := provider.Info(key)
	if err == nil {
		// calculate etag
		etag, err := backup.QEtag(filepath)
		if err != nil {
			return err
		}
		// same file
		if entry.Hash == etag {
			return backup.ErrDuplicated
		}

		// same name
		reader := bufio.NewReader(os.Stdin)
	READPUTACTION:
		fmt.Fprintf(os.Stdout, "file: %s already exists!, type [A]bort, [O]veride or [R]ename this file?\nOption: ", filepath)
		text, _ := reader.ReadString('\n')
		switch strings.TrimSpace(strings.ToLower(text)) {
		case "a":
			return nil
		case "o":
			return provider.Put(filepath, key)
		case "r":
			fmt.Fprint(os.Stdout, "Please type new file name: ")
			name, _ := reader.ReadString('\n')
			return putFile(filepath, strings.TrimSpace(name))
		default:
			goto READPUTACTION
		}
	} else if err != backup.ErrNotFound {
		return err
	} else {
		return provider.Put(filepath, key)
	}
}

func initConf() {
	defaultDir := fmt.Sprintf("%s/.config/qb", os.Getenv("HOME"))
	defaultConf = fmt.Sprintf("%s/config.json", defaultDir)
	backup.Conf = backup.NewConfig()

	err := backup.Conf.Restore(defaultConf)
	if err != nil {
		if os.IsNotExist(err) {
			if _, err := os.Stat(defaultDir); os.IsNotExist(err) {
				if err := os.MkdirAll(defaultDir, 0700); err != nil {
					panic(err)
				}
			}
			if _, err := os.Create(defaultConf); err != nil {
				fmt.Println(defaultConf)
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	validated := backup.Conf.Validate()
	if validated == false {
		reader := bufio.NewReader(os.Stdin)
		if backup.Conf.AccessKey == "" {
			fmt.Fprintf(os.Stdout, "Qiniu Access Key:")
			ak, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			backup.Conf.AccessKey = strings.TrimSpace(ak)
		}

		if backup.Conf.Secretkey == "" {
			fmt.Fprintf(os.Stdout, "Qiniu Secret Key:")
			sk, err := gopass.GetPasswdMasked()
			if err != nil {
				panic(err)
			}
			backup.Conf.Secretkey = string(sk)
		}

		if backup.Conf.Domain == "" {
			fmt.Fprintf(os.Stdout, "Qiniu Domain for file download:")
			domain, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			backup.Conf.Domain = strings.TrimSpace(domain)
		}

		if backup.Conf.Bucket == "" {
			fmt.Fprintf(os.Stdout, "Qiniu bucket:")
			bucket, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			backup.Conf.Bucket = strings.TrimSpace(bucket)
		}
	}

	provider = backup.NewQBackup(backup.Conf.AccessKey, backup.Conf.Secretkey, backup.Conf.Domain, backup.Conf.Bucket)
	if validated == false {
		err = provider.List("", "", 1)
		if err != nil {
			panic(err)
		}

		err = backup.Conf.Save(defaultConf)
		if err != nil {
			panic(err)
		}
	}
}
