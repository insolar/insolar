package main

import (
	"fmt"
	"sync"
)

type QueueItem struct {
	itemType uint
	signal   uint
	index    uint
	payload  interface{}
	next     *QueueItem
}

var emptyQueueItem QueueItem

func init() {
	emptyQueueItem = QueueItem{
		itemType: 0,
		signal:   0,
		index:    0,
		payload:  nil,
		next:     nil,
	}
}

type Queue struct {
	head *QueueItem
}

func NewQueue() *Queue {
	queue := &Queue{
		head: &emptyQueueItem,
	}

	return queue
}

func (q *Queue) sinkPush(data interface{}) bool {
	fmt.Println("Pushing: ", data)
	newNode := &QueueItem{
		payload: data,
		next:    q.head,
	}

	q.head = newNode

	return true
}

func (q *Queue) Next() (*QueueItem, error) {
	if *q.head == emptyQueueItem {
		return nil, fmt.Errorf("Empty queue")
	}
	retElement := q.head
	q.head = q.head.next

	return retElement, nil
}

func main() {
	queue := NewQueue()

	parallel := 200
	wg := sync.WaitGroup{}
	wg.Add(parallel)

	numIterations := 20

	for i := 0; i < parallel; i++ {
		go func(wg *sync.WaitGroup, q *Queue) {
			fmt.Println("START")
			for i := 0; i < numIterations; i++ {
				q.sinkPush(i)
			}
			wg.Done()
		}(&wg, queue)
	}

	wg.Wait()

	numElement := 0
	for true {
		element, err := queue.Next()
		if err != nil {
			fmt.Printf("Got error %s. Stop.\n", err)
			break
		}
		numElement++

		fmt.Println("Next element: ", element.payload)
	}

	fmt.Printf("Num Elements: %d . Must be: %d\n", numElement, parallel*numIterations)

}
