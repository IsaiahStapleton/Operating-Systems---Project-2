package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

var total_num_word int

func count_words(s string) int {
	return len(strings.Fields(s))
}

// Consumer task to operate on queue
func consumer_task(task_num int, ch <-chan string) {

	
	for item := range ch {
		fmt.Printf("Task %d consuming line: %v\n", task_num, item)
		word_count := count_words(item)
		fmt.Printf("	Number of words in line: %d\n", word_count)
		total_num_word += word_count
	}
	// each thread will drop out of their loop when channel is closed
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
			consumer_task(i, queue) 
			wg.Done()
		}(i)
	}

	wg.Wait()
	fmt.Println("All done")
	fmt.Println(total_num_word)

}
