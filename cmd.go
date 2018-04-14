package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

type HateClass struct {
	Text     string `json:"text"`
	TopClass string `json:"top_class"`
	Classes  []struct {
		ClassName  string  `json:"class_name"`
		Confidence float64 `json:"confidence"`
	} `json:"classes"`
}

func runHateDetectionCmd(message string) (HateClass, error) {
	var hc HateClass

	cmd := exec.Command("python", "hate.py", "check", fmt.Sprintf("--txt=%s", message))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	var x string
	go func() {
		x = copyOutput(stdout)
	}()

	cmd.Wait()

	hc, err = marshalClass(x)
	if err != nil {
		fmt.Println(err)
		return hc, err
	}

	return hc, err
}

func copyOutput(r io.Reader) string {
	var str string
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		str = str + scanner.Text()
	}
	return str
}

func marshalClass(in string) (HateClass, error) {
	var hc HateClass
	fmt.Println(in)

	err := json.Unmarshal([]byte(in), &hc)
	if err != nil {
		fmt.Println(err)
		return hc, err
	}

	return hc, err
}
