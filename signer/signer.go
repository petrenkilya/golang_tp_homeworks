package main

import (
	"fmt"
	"sort"
	"strconv"
	strings2 "strings"
	"sync"
)

func ExecutePipeline(jobs ...job) []interface{} {
	numberOfChannels := len(jobs)
	var channels []chan interface{}
	for i := 0; i < numberOfChannels; i++ {
		channels = append(channels, make(chan interface{}, 1))
	}

	preparedWorker := func(worker job, in, out chan interface{}) {
		worker(in, out)
		close(out)
	}

	for index, worker := range jobs {
		if index == 0 {
			go preparedWorker(worker, nil, channels[index])
			continue
		}
		go preparedWorker(worker, channels[index-1], channels[index])
	}

	var result []interface{}

	for item := range channels[len(channels)-1] {
		result = append(result, item)
	}
	return result
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
			md5OfDataChan := make(chan string, 1)
			crcOfDataChan := make(chan string, 1)
			crcOfmd5DataChan := make(chan string, 1)

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
		go func(data string, oldDataSendedMutex *sync.Mutex, newDataSendedMutex *sync.Mutex) {
			var channels []chan string
			for i := 0; i < 6; i++ {
				currentChannel := make(chan string, 1)

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
		}(data, oldDataSendedMutex, newDataSendedMutex)
		oldDataSendedMutex = newDataSendedMutex
	}
	waiter.Wait()
}

func CombineResults(in, out chan interface{}) {
	var strings []string
	for input := range in {
		data, ok := input.(string)
		if !ok {
			return
		}
		strings = append(strings, data)
	}
	sort.Strings(strings)
	result := strings2.Join(strings, "_")

	out <- result
	fmt.Printf("CombineResults %s\n", result)
}
