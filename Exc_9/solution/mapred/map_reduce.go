package mapred

import (
	"regexp"
	"runtime"
	"strings"
	"sync"
)

type MapReduce struct {
}

// todo implement mapreduce
var nonAZ = regexp.MustCompile(`[^a-z]+`) // remove numbers and special characters

func (mr MapReduce) Run(input []string) map[string]int {
	//MAP in parallel using goroutines
	n := runtime.NumCPU()

	jobs := make(chan string, len(input)) //channel for distributing lines of text
	mapped := make(chan []KeyValue)       //channel to store the mapped results (key-value)
	//waitgroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(n)

	//launch the goroutines to process the lines in parallel
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for line := range jobs {
				mapped <- mr.wordCountMapper(line) //send the mapped results to de mapped channel
			}
		}()
	}

	//distribute the input lines to the jobs channel
	go func() {
		for _, line := range input {
			jobs <- line
		}
		close(jobs)
	}()

	//wait for goroutines to finish and then close the mapped channel
	go func() {
		wg.Wait()
		close(mapped)
	}()

	//SHUFFLE
	shuffled := make(map[string][]int) //map to store the shuffled results
	//iterate the mapped results
	for kvs := range mapped {
		for _, kv := range kvs {
			shuffled[kv.Key] = append(shuffled[kv.Key], kv.Value)
		}
	}

	//REDUCE: sum the counts for each word
	result := make(map[string]int, len(shuffled))
	for k, vals := range shuffled {
		kv := mr.wordCountReducer(k, vals)
		result[kv.Key] = kv.Value
	}

	return result
}

// convert a line of text to a list of key-value pairs
func (mr MapReduce) wordCountMapper(text string) []KeyValue {
	s := strings.ToLower(text)
	s = nonAZ.ReplaceAllString(s, " ") //replace the non-alphabetic characters with space
	words := strings.Fields(s)         //split in words

	var out []KeyValue
	for _, w := range words {
		out = append(out, KeyValue{Key: w, Value: 1}) //each word has a count of 1
	}
	return out
}

// sum the occurrences of a word
func (mr MapReduce) wordCountReducer(key string, values []int) KeyValue {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return KeyValue{Key: key, Value: sum}
}
