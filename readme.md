# Isaiah Stapleton & Emma Steuer

## This purpose of this project is to build understanding of threading and process synchronization and to implement a single producer/multiple consumer queue

### This program accepts two parameters: number of threads, name of text file. The program will then start up the specified number of threads (consumers) and start a routine that adds each line to a queue line by line (producer). Each consumer thread will pop a line off of the queue one by one and print their task number, the line, and number of words in that line. When all lines of input have been read, the consumer tasks will terminate and the word counts will be accumlated and printed out for the user.

### In order to run this program, navigate to the directory where the program is located on your host machine. Type "go run p2.go". When prompted, enter number of consumer threads you would like to start, followed by the name of a text file that you would like to be processed. Make sure that the text file is located within the same directory as the program.