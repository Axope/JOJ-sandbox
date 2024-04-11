package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type Config struct {
	MemLimit  int64    `json:"memLimit"`
	TimeLimit int64    `json:"timeLimit"`
	Solution  string   `json:"solution"`
	TestCases []string `json:"testCases,omitempty"`
}

const (
	compileCmd = "compile"
	runCmd     = "run"

	compileJsonPath = "./config/compile.json"
	runJsonPath     = "./config/run.json"
)

func parseJson(tp string) Config {
	filePath := compileJsonPath
	if tp == runCmd {
		filePath = runJsonPath
	}

	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("Error opening config file, err = %s", err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	err = decoder.Decode(&config)
	if err != nil {
		panic(fmt.Sprintf("Error decoding config file, err = %s", err))
	}

	fmt.Printf("parse config: %+v\n", config)
	return config
}

func kill(pid int) {
	pgid, err := syscall.Getpgid(pid)
	if err != nil {
		fmt.Println("kill: Getpgid error:", err)
		return
	}

	killCmd := exec.Command("kill", "-TERM", fmt.Sprintf("-%d", pgid))
	err = killCmd.Run()
	if err != nil {
		fmt.Println("kill cmd error:", err)
		return
	}
}

func compile(config Config) {
	cmd := exec.Command("/bin/bash", "-c",
		fmt.Sprintf("sh ./script/compile.sh %d %s", config.MemLimit, config.Solution))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("Failed to start process, err = %v", err))
	}

	waitCh := make(chan error)
	go func() {
		waitCh <- cmd.Wait()
	}()

	select {
	case <-time.After(time.Duration(config.TimeLimit)):
		fmt.Println("timeout")
		kill(cmd.Process.Pid)
		return

	case err := <-waitCh:
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func readWords(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		line := scanner.Text()
		lineWords := strings.Fields(line)
		for _, word := range lineWords {
			if word != "" {
				words = append(words, word)
			}
		}
	}
	return words, nil
}

func compare(v string) bool {
	s1, err := readWords(fmt.Sprintf("./output/%s.out", v))
	if err != nil {
		return false
	}
	s2, err := readWords(fmt.Sprintf("./data/%s.ans", v))
	if err != nil {
		return false
	}

	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func runCases(config Config) {
	testCases := config.TestCases

	for _, v := range testCases {
		cmd := exec.Command("/bin/bash", "-c",
			fmt.Sprintf("sh ./script/run.sh %d %s %s", config.MemLimit, config.Solution, v))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			panic(fmt.Sprintf("Failed to start process, err = %v", err))
		}

		waitCh := make(chan error)
		go func() {
			waitCh <- cmd.Wait()
		}()

		select {
		case <-time.After(time.Duration(config.TimeLimit)):
			fmt.Println("timeout")
			// kill(cmd.Process.Pid)
			os.Exit(4)

		case err := <-waitCh:
			if err != nil {
				fmt.Println("waitCh:", err)
				// fmt.Println("exitCode:", cmd.ProcessState.ExitCode())
				os.Exit(cmd.ProcessState.ExitCode())
			}
		}

		if compare(v) {
			fmt.Printf("%v ok\n", v)
		} else {
			fmt.Printf("%v wrong answer\n", v)
			os.Exit(2)
		}
	}
}

func main() {
	var tp string
	flag.StringVar(&tp, "type", compileCmd, "Usage: -type compile/run")
	flag.Parse()

	config := parseJson(tp)

	switch tp {
	case compileCmd:
		compile(config)
	case runCmd:
		runCases(config)
	}
}
