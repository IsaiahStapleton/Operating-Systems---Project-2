package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

// Function to count the number of words in a string
func count_words(value string) int {
    // Match non-space character sequences.
    re := regexp.MustCompile(`[\S]+`)
    
    // Find all matches and return count.
    results := re.FindAllString(value, -1)
    return len(results)
}

// Consumer task to operate on queue
func consumer_task(task_num int, ch <-chan string, ch2 chan<- int) {



	// Loop through the channel and output desired information and send word count of line back through the channel
	// in order to communicate with other threads
	for item := range ch {
		word_count := count_words(item)
		fmt.Printf("\nTask %d consuming line: %v\nNumber of words: %d", task_num, item, word_count)
		ch2 <- word_count
	}

}

func main() {

	// Initialize queue
	queue := make(chan string)
	total_words := make(chan int)
	var totalWords = 0

	// Initialize wait group
	var wg sync.WaitGroup

	// Get num of tasks to run as well as name of the file from user
	fmt.Printf("Enter number of tasks to run followed by name of text file (ex. 5 text.txt): ")
	var numof_tasks int = 0
	var file_name string
	fmt.Scanf("%d %s", &numof_tasks, &file_name)


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

	for i := 1; i <= numof_tasks; i++ {
		wg.Add(1)
		go func(i int) {
			consumer_task(i, queue, total_words)
			wg.Done()
		}(i)
	}

	// Loop through each line in the file and append it to the queue
	go func() {
		// Loop through each line in the file and append it to the queue
		for scanner.Scan() {
			queue <- scanner.Text()
			
		}
		// Signals consumer threads that there is no more lines to be read
		close(queue) 
	}()


	go func() {
		for wordCount := range total_words {
			totalWords += wordCount
		}
		close(total_words)
	}()

	wg.Wait()

	fmt.Printf("\nDone")

	fmt.Printf("\nTotal number of words: %d\n", totalWords)
}
