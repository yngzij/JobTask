package main

import (
	"fmt"
	"time"

	"JobTask/job"
)

type Order struct {
}

func (o Order) Prepare() {
	time.Sleep(4 * time.Second)
	fmt.Println("Order Prepare")
}

type OrderGoods struct {
}

func (o OrderGoods) Prepare() {
	//TODO implement me
	time.Sleep(4 * time.Second)
	fmt.Println("OrderGoods Prepare")
}

func (o OrderGoods) Do() {
	//TODO implement me
	time.Sleep(3 * time.Second)
	fmt.Println("OrderGoods")
}

func (o Order) Do() {
	time.Sleep(3 * time.Second)
	fmt.Println("Order")
}

func main() {
	j := job.NewJob()
	j.PushJob(Order{}, OrderGoods{})
	j.PushJob(Order{}, OrderGoods{})
	j.PushJob(Order{}, OrderGoods{})
	j.PushJob(Order{}, OrderGoods{})

	//j.PushJob(job.AsyncJob, Order{}, OrderGoods{})
	j.Run()
}
