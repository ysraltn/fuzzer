package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type fuzzTask struct {
	url     string
	payload string
}

func fuzzerWorker(id int, tasks <-chan fuzzTask, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	for task := range tasks {
		fuzz(task.url, task.payload, id)
	}
}

func fuzz(url string, payload string, workerId int) {
	target := url + "/" + payload
	response, err := http.Get(target)
	if err != nil {
		fmt.Println(err)
	}
	if response.StatusCode == 200 {
		fmt.Println(target)
		fmt.Println(response.StatusCode)
		fmt.Println("ok\n")
	}

}

func help() {
	fmt.Println("\nYavuzlar Fuzzer Tool\n")
	fmt.Println("usage:")
	fmt.Println("go run main.go [target url] [payload (text file)]\n")
}

func main() {

	args := os.Args
	if len(args) != 3 {
		help()
		return
	}

	url := args[1]
	filePath := args[2]
	//filePathToWrite := os.Args[2]
	extension := filepath.Ext(filePath)
	if extension != ".txt" {
		fmt.Println("payload file must have .txt extension")
		help()
		return
	}
	numberOfWorkers := 3

	taskChannel := make(chan fuzzTask, numberOfWorkers)
	var waitGroup sync.WaitGroup

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for i := 1; i <= numberOfWorkers; i++ {
		waitGroup.Add(1)
		go fuzzerWorker(i, taskChannel, &waitGroup)
	}

	for scanner.Scan() {
		payload := scanner.Text()
		taskChannel <- fuzzTask{url: url, payload: payload}
	}
	close(taskChannel)
	waitGroup.Wait()

}
