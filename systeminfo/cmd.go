package systeminfo

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"sync"
)

func NewCmdMgr() *Cmd {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Cmd{EchoString: make(map[string][]string), ctx: ctx, exit: cancel}
	return c
}

type Cmd struct {
	m          sync.Mutex
	EchoString map[string][]string
	ctx        context.Context
	exit       context.CancelFunc
}

func Ascii2Utf8(b []byte) (string, error) {
	i := 0
	for ; i < len(b); i++ {
		if b[i] == 0 {
			break
		}
	}

	if 0 == i {
		return "", nil
	}

	var Byte []byte
	var err error
	if Byte, err = simplifiedchinese.GBK.NewDecoder().Bytes(b[:i]); err != nil {
		return "", errors.New(fmt.Sprintf("Ascii2Utf8 convert %v failed", b))
	}

	return string(Byte[:]), nil
}

func (c *Cmd) ExecCommand(guid string, PingAddress string) string {
	c.m.Lock()
	EchoString := make([]string, 0)
	c.EchoString[guid] = EchoString
	c.m.Unlock()
	if _, ok := c.EchoString[guid]; !ok {
		return "error"
	}
	execstr := fmt.Sprintf("ping %s -n 10", PingAddress)
	go c.execCommand(guid, execstr)

	return execstr
}

func (c *Cmd) execCommand(guid string, PingAddress string) error {

	cmd := exec.Command("cmd", "/K", PingAddress)
	//显示运行的命令
	//fmt.Println(cmd.Args)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			c.EchoString[guid] = append(c.EchoString[guid], "EOF")
			break
		}

		utf, err := Ascii2Utf8([]byte(line))
		if err != nil {
			c.EchoString[guid] = append(c.EchoString[guid], "EOF")
			break
		}
		c.m.Lock()
		c.EchoString[guid] = append(c.EchoString[guid], utf)
		c.m.Unlock()
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func (c *Cmd) GetResp(guid string, i int) []string {
	if s, ok := c.EchoString[guid]; ok {
		if len(s) == 0 {
			return s
		}
		if s[len(s)-1] == "EOF" {
			c.m.Lock()
			delete(c.EchoString, guid)
			c.m.Unlock()

		}
		return s
	}
	return nil
}
