package getconfig

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"os"
	"regexp"
	"strings"
)

func GetCVS(cvsFile string) (csvall [][]string, err error) {
	f, err := os.Open(cvsFile)
	if err != nil {
		return
	}
	defer f.Close()
	r := csv.NewReader(f)
	csvall, err = r.ReadAll()
	if err != nil {
		return
	}
	return
}

func GetIDList(listFile string) (idList []string, err error) {
	//打开文件
	f, err := os.Open(listFile)
	if err != nil {
		return
	}
	defer f.Close()
	//读取文件到buffer里边
	buf := bufio.NewReader(f)
	for {
		//按照换行读取每一行
		l, err := buf.ReadString('\n')
		//bom头处理
		lb := []byte(l)
		if len(lb) > 2 {
			bom := []byte{0xef, 0xbb, 0xbf}
			if bytes.Compare(lb[0:3], bom) == 0 {
				l = string(lb[3:])
			}
		}
		//相当于PHP的trim
		line := strings.TrimSpace(l)
		//判断退出循环
		if err != nil {
			if err != io.EOF {
				//return err
				return nil, err
			}
			if len(line) == 0 {
				break
			}
		}

		r := regexp.MustCompile("^\\d*$")
		rs := r.FindStringSubmatch(line)
		if len(rs) > 0 {
			idList = append(idList, line)
		}

	}
	return
}

func Getconfigini(conffile string) (cfg map[string]map[string]string, err error) {
	//实例化这个map
	cfg = map[string]map[string]string{}
	//打开这个ini文件
	f, err := os.Open(conffile)
	if err != nil {
		return
	}
	defer f.Close()
	//读取文件到buffer里边
	buf := bufio.NewReader(f)
	section := "default"
	for {
		//按照换行读取每一行
		l, err := buf.ReadString('\n')
		//bom头处理
		lb := []byte(l)
		if len(lb) > 2 {
			bom := []byte{0xef, 0xbb, 0xbf}
			if bytes.Compare(lb[0:3], bom) == 0 {
				l = string(lb[3:])
			}
		}
		//相当于PHP的trim
		line := strings.TrimSpace(l)
		//判断退出循环
		if err != nil {
			if err != io.EOF {
				//return err
				return nil, err
			}
			if len(line) == 0 {
				break
			}
		}

		switch {
		case len(line) == 0:
		case line[0] == '#':
		//匹配[]然后存储
		case line[0] == '[' && line[len(line)-1] == ']':
			section = strings.TrimSpace(line[1 : len(line)-1])
		default:
			//xxx = yyy 这种的可以匹配存储
			i := strings.IndexAny(line, "=")

			cfgn, ok := cfg[section]
			if !ok {
				cfgn = make(map[string]string)
				cfg[section] = cfgn
			}

			cfgn[strings.TrimSpace(line[0:i])] = strings.TrimSpace(line[i+1:])

		}
	}
	return
}
