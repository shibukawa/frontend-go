package frontend

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/shlex"
)

func parseCmd(cmdStr string) (cmd string, args []string, err error) {
	l := shlex.NewLexer(strings.NewReader(cmdStr))
	cmd, err = l.Next()
	if err != nil {
		return
	}
	for {
		token, err := l.Next()
		if err != nil {
			break
		}
		args = append(args, token)
	}
	return
}

type devServer struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func startDevServer(ctx context.Context, folder, cmdStr string) (d *devServer, host string, err error) {
	ctx, cancel := context.WithCancel(ctx)
	d = &devServer{
		ctx:    ctx,
		cancel: cancel,
	}

	cmdName, args, err := parseCmd(cmdStr)
	if err != nil {
		return nil, "", err
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	cmd.Dir = folder
	stdout, _ := cmd.StdoutPipe()
	ch := make(chan string)
	go func() {
		re := regexp.MustCompile(`(http://[0-9A-Za-z.]+:\d+)`)
		foundPort := false
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if !foundPort {
				m := re.FindStringSubmatch(scanner.Text())
				if len(m) > 0 {
					ch <- m[0]
					foundPort = true
				}
			}
			fmt.Println(scanner.Text())
		}
	}()
	err = cmd.Start()
	if err != nil {
		return
	}
	host = <-ch
	go func() {
		<-ctx.Done()
		d.Stop()
	}()
	return
}

func (d *devServer) Stop() {
	if d.cancel == nil {
		return
	}
	d.cancel()
	d.cancel = nil
	d.Wait()
}

func (d *devServer) Wait() {
	<-d.ctx.Done()
}
