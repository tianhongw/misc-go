package util

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func runCmd(name string, arg ...string) (stdOutBytes []byte, stdErrBytes []byte, err error) {
	if _, err = exec.LookPath(name); err != nil {
		log.Printf("%s %s\n%s", name, strings.Join(arg, " "), err.Error())
		return
	}
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	log.Println(cmd.String())
	if err = cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	stdOutBytes, err = ioutil.ReadAll(stdout)
	if err != nil {
		return
	}
	stdErrBytes, err = ioutil.ReadAll(stderr)
	if err != nil {
		return
	}
	if err = cmd.Wait(); err != nil {
		log.Printf("%s", stdErrBytes)
	}
	return
}

func runCmdAndGetStdOutInTime(stdOutMsg chan string, name string, arg ...string) (stdErrBytes []byte, err error) {
	if _, err = exec.LookPath(name); err != nil {
		log.Printf("%s %s\n%s", name, strings.Join(arg, " "), err.Error())
		return
	}
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	log.Println(cmd.String())
	if err = cmd.Start(); err != nil {
		log.Println(err)
		return
	}
	reader := bufio.NewReader(stdout)
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err)
			return
		}
		stdOutMsg <- line
	}
	stdErrBytes, err = ioutil.ReadAll(stderr)
	if err != nil {
		return
	}
	if err = cmd.Wait(); err != nil {
		log.Printf("%s", stdErrBytes)
	}
	return
}
