package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const MultiHashTh = 6
const MaxInputLen = 1

func ExecutePipeline(jobs ...job) {
	waiter := &sync.WaitGroup{}

	preparedWorker := func(worker job, in, out chan interface{}) {
		worker(in, out)
		close(out)
		waiter.Done()
	}

	var prevChannel chan interface{} = nil
	var nextChannel chan interface{} = nil
	for i := 0; i < len(jobs); i++ {
		nextChannel = make(chan interface{}, MaxInputLen)
		waiter.Add(1)
		go preparedWorker(jobs[i], prevChannel, nextChannel)
		prevChannel = nextChannel
	}

	waiter.Wait()
}

func SingleHash(in, out chan interface{}) {
	md5Mutex := &sync.Mutex{}
	oldDataSendedMutex := &sync.Mutex{}
	waiter := &sync.WaitGroup{}

	for item := range in {
		number, ok := item.(int)
		if !ok {
			return
		}
		data := strconv.Itoa(number)

		newDataSendedMutex := &sync.Mutex{}
		newDataSendedMutex.Lock()

		waiter.Add(1)
		go func(data string, oldDataSendedMutex *sync.Mutex, newDataSendedMutex *sync.Mutex) {
			md5OfDataChan := make(chan string, MaxInputLen)
			crcOfDataChan := make(chan string, MaxInputLen)
			crcOfmd5DataChan := make(chan string, MaxInputLen)

			go func() {
				md5Mutex.Lock()
				md5OfDataChan <- DataSignerMd5(data)
				md5Mutex.Unlock()
			}()
			go func() {
				crcOfDataChan <- DataSignerCrc32(data)
			}()

			md5OfData := <-md5OfDataChan

			go func() {
				crcOfmd5DataChan <- DataSignerCrc32(md5OfData)
			}()
			crcOfmd5Data := <-crcOfmd5DataChan
			crcOfData := <-crcOfDataChan

			result := crcOfData + "~" + crcOfmd5Data

			close(crcOfDataChan)
			close(crcOfmd5DataChan)
			close(md5OfDataChan)

			oldDataSendedMutex.Lock()
			out <- result
			oldDataSendedMutex.Unlock()
			newDataSendedMutex.Unlock()

			fmt.Printf("%s SingleHash data %s\n", data, data)
			fmt.Printf("%s SingleHash md5(data) %s\n", data, md5OfData)
			fmt.Printf("%s SingleHash crc32(md5(data)) %s\n", data, crcOfmd5Data)
			fmt.Printf("%s SingleHash crc32(data) %s\n", data, crcOfData)
			fmt.Printf("%s SingleHash result %s\n", data, result)
			waiter.Done()
		}(data, oldDataSendedMutex, newDataSendedMutex)
		oldDataSendedMutex = newDataSendedMutex
	}
	waiter.Wait()
}

func multiHashWorker(data string, oldDataSendedMutex *sync.Mutex, newDataSendedMutex *sync.Mutex, out chan interface{},
	waiter *sync.WaitGroup) {
	var channels []chan string
	for i := 0; i < MultiHashTh; i++ {
		currentChannel := make(chan string, MaxInputLen)

		go func(data string, index int) {
			th := strconv.Itoa(index)

			currentChannel <- DataSignerCrc32(th + data)
		}(data, i)

		channels = append(channels, currentChannel)
	}
	result := ""

	for index, ch := range channels {
		buffer := <-ch
		result += buffer
		fmt.Printf("%s MultiHash: crc32(th+step1) %d %s\n", data, index, buffer)
	}
	oldDataSendedMutex.Lock()
	out <- result
	oldDataSendedMutex.Unlock()
	newDataSendedMutex.Unlock()

	fmt.Printf("%s MultiHash result: %s\n", data, result)
	waiter.Done()
}

func MultiHash(in, out chan interface{}) {
	oldDataSendedMutex := &sync.Mutex{}
	waiter := &sync.WaitGroup{}

	for input := range in {
		data, ok := input.(string)
		if !ok {
			return
		}

		newDataSendedMutex := &sync.Mutex{}
		newDataSendedMutex.Lock()
		waiter.Add(1)
		go multiHashWorker(data, oldDataSendedMutex, newDataSendedMutex, out, waiter)
		oldDataSendedMutex = newDataSendedMutex
	}
	waiter.Wait()
}

func CombineResults(in, out chan interface{}) {
	var outStrings []string
	for input := range in {
		data, ok := input.(string)
		if !ok {
			return
		}
		outStrings = append(outStrings, data)
	}
	sort.Strings(outStrings)
	result := strings.Join(outStrings, "_")

	out <- result
	fmt.Printf("CombineResults %s\n", result)
}
