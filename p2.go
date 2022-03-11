package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

// Consumer task to operate on queue
func consumer_task(task_num int) {
	
	fmt.Printf("I'm consumer task #%v ", task_num)
	fmt.Println("Line being popped off queue: " + queue[0])

	queue = queue[1:]


}

// Initialize queue
var queue = make([]string, 0)


func main() {

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
	for scanner.Scan() {
        line := scanner.Text()  
		queue = append(queue, line)
    }

	// Start specified # of consumer tasks
	for i := 1; i <= numof_tasks; i++ {
		wg.Add(1)
		go func(i int) { 
			consumer_task(i) 
			wg.Done()
		}(i)
	}


	wg.Wait()
	fmt.Println("All done")
	fmt.Println(queue)
}
