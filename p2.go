package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

// Counts the number of words in a string
func count_words(s string) int {
	return len(strings.Fields(s))
}

// Consumer task to operate on queue
func consumer_task(task_num int, ch <-chan string, ch2 chan<- int) {

	// Loop through the channel and output desired information and send word count of line back through the channel
	// in order to communicate with other threads
	for item := range ch {
		word_count := count_words(item) 
		fmt.Printf("\nTask %d consuming line: %v\nNumber of words in line: %d", task_num, item, word_count)
		ch2 <- word_count
	}

}

func main() {

	// Initialize channels
	queue := make(chan string)
	word_counts := make(chan int)

	// Variable for keeping track of the total number of words
	total_words := 0

	// Initialize wait group
	var wg sync.WaitGroup

	// Get num of tasks to run as well as name of the file from user
	fmt.Printf("Enter number of consumer tasks to run followed by name of text file (ex. 5 text.txt): ")
	var numof_tasks int = 0
	var file_name string
	fmt.Scanf("%d %s", &numof_tasks, &file_name)


	// Read the specified file in
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Scanner to scan the file
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Start consumer tasks and add them to the wait group
	for i := 1; i <= numof_tasks; i++ {
		wg.Add(1)
		go func(i int) {
			consumer_task(i, queue, word_counts)
			wg.Done()
		}(i)
	}

	// Loop through each line in the file and append it to the queue, when all lines have been read
	// signal the queue channel to close
	go func() {
	
		for scanner.Scan() {
			queue <- scanner.Text()		
		}
		close(queue) 
	}()


	// Accumlate the word counts, when all word counts have been accumulated signal the word_counts channel to close
	go func() {
		for words := range word_counts {
			total_words += words
		}
		close(word_counts)
	}()

	// Wait for threads to finish executing and then output the total word count
	wg.Wait()
	fmt.Printf("\nDone")
	fmt.Printf("\nTotal number of words: %d\n", total_words)

}
