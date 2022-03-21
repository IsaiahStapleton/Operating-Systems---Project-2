package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

// Consumer task to operate on queue
func consumer_task(task_num int, ch <-chan string) {

	fmt.Printf("I'm consumer task #%v\n", task_num)
	for item := range ch {
		fmt.Printf("task %d consuming: Line item: %v\n", task_num, item)
	}
	// each worker will drop out of their loop when channel is closed
}

func main() {
	// Initialize queue
	queue := make(chan string)
	// Initialize wait group
	var wg sync.WaitGroup
	// Get number of tasks to run from user
	var numof_tasks int
	fmt.Print("Enter number of tasks to run: ")
	fmt.Scan(&numof_tasks)

	// Open file
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Scanner to scan the file
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Loop through each line in the file and append it to the queue
	go func() {
		// Loop through each line in the file and append it to the queue
		for scanner.Scan() {
			queue <- scanner.Text()
		}
		close(queue) // signal to workers that there is no more items
	}()

	// Start specified # of consumer tasks
	for i := 1; i <= numof_tasks; i++ {
		wg.Add(1)
		go func(i int) {
			consumer_task(i, queue) // add the queue parameter
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("All done")
	fmt.Println(queue)
}
