package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func GetFile() {
	log.Println("缓存文件读取...")

	fi, err := os.OpenFile("keys.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewReader(fi) // 创建 Reader
	for {
		lineB, err := r.ReadBytes('\n')
		if len(lineB) > 3 {
			KeysFile = append(KeysFile, strings.TrimSpace(string(lineB)))
		}
		if err != nil {
			break
		}
	}

	fi2, err2 := os.OpenFile("urls.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err2 != nil {
		panic(err2)
	}
	defer fi2.Close()
	r2 := bufio.NewReader(fi2) // 创建 Reader
	for {
		lineB, err := r2.ReadBytes('\n')
		if len(lineB) > 3 {
			UrlsFile = append(UrlsFile, strings.TrimSpace(string(lineB)))
		}
		if err != nil {
			break
		}
	}

	fi3, err3 := os.OpenFile("ids.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err3 != nil {
		panic(err2)
	}
	defer fi3.Close()
	r3 := bufio.NewReader(fi3) // 创建 Reader
	for {
		lineB, err := r3.ReadBytes('\n')
		if len(lineB) > 3 {
			IdsFile = append(IdsFile, strings.TrimSpace(string(lineB)))
		}
		if err != nil {
			break
		}
	}

	fi4, err := os.OpenFile("keyss.txt", os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi4.Close()
	r4 := bufio.NewReader(fi4) // 创建 Reader
	for {
		lineB, err := r4.ReadBytes('\n')
		if len(lineB) > 3 {
			KeyssFile = append(KeyssFile, strings.TrimSpace(string(lineB)))
		}
		if err != nil {
			break
		}
	}
}

func OutFileKeys(keys []string) {
	MuxKey.Lock()
	defer MuxKey.Unlock()

	fi, err := os.OpenFile("keys.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewWriter(fi) // 创建 Reader
	for i := range keys {
		r.WriteString(keys[i] + "\n")
		KeysFile = append(KeysFile, keys[i])
	}
	r.Flush()
}

func OutFileUrls(keys []string) {
	MuxUrl.Lock()
	defer MuxUrl.Unlock()

	fi, err := os.OpenFile("urls.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	r := bufio.NewWriter(fi) // 创建 Reader
	for i := range keys {
		r.WriteString(keys[i] + "\n")
		KeysFile = append(KeysFile, keys[i])
	}
	r.Flush()
}
