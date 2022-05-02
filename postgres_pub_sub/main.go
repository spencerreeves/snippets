package main

import (
	"fmt"
	"github.com/lib/pq"
	"time"
)

const connStr = "postgres://postgres@127.0.0.1:5432/spencerreeves?sslmode=disable"
const psqlInfo = "dbname=spencerreeves user=postgres sslmode=disable"

func main() {
	fmt.Println("Start")
	listener := pq.NewListener(psqlInfo, 10*time.Second, time.Minute, ReportProblem)
	if err := listener.Listen("job_channel"); err != nil {
		panic(err)
	}

	Simple(listener)
	fmt.Println("End")
}

func Simple(listener *pq.Listener) {
	go func() {
		fmt.Println("inside")
		for {
			select {
			case notification := <-listener.Notify:
				fmt.Printf("Recieved postgres notification: %v\n", notification)
			}
		}
	}()
}

func Complex() {
	//db, _ := sqlx.Connect("postgres", connStr)
	//const workerCount = 5
	//// Create a pool of workers
	//var wg sync.WaitGroup
	//wg.Add(workerCount)
	//
	//fmt.Println("Starting listening service...")
	//for i := 0; i < workerCount; i++ {
	//	go func(i int) {
	//		defer wg.Done()
	//		fmt.Printf("Worker %v ")
	//		val := slice[i]
	//		fmt.Printf("i: %v, val: %v\n", i, val)
	//	}(i)
	//}
	//wg.Wait()
	//fmt.Println("Finished for loop")

}

func ReportProblem(event pq.ListenerEventType, err error) {
	fmt.Printf("problem listening to `job_channel`: %v\n", err)
}
