package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting Simulation")
	c12 := make(chan bool)
	cf1 := make(chan int, 1)
	c23 := make(chan bool)
	cf2 := make(chan int, 1)
	c34 := make(chan bool)
	cf3 := make(chan int, 1)
	c45 := make(chan bool)
	cf4 := make(chan int, 1)
	c51 := make(chan bool)
	cf5 := make(chan int, 1)

	fmt.Println("Channels Initialised")
	go fork(cf1, cf2)
	go philo(1, c51, c12, cf1)
	go fork(cf2, cf3)
	go philo(2, c12, c23, cf2)
	go fork(cf3, cf4)
	go philo(3, c23, c34, cf3)
	go fork(cf4, cf5)
	go philo(4, c34, c45, cf4)
	go fork(cf5, cf1)
	go philo(5, c45, c51, cf5)
	fmt.Println("Table Initialised")
	time.Sleep(10 * time.Second)
	fmt.Println("Program Terminated")
}

func philo(id int, plchan chan bool, prchan chan bool, forkchan chan int) {
	fmt.Println("Philosopher ", id, " created")
	var pl bool
	var pr bool
	var f1 int
	var f2 int
	for {
		if !pl && !pr {
			plchan <- true
			prchan <- true
			forkchan <- id
			forkchan <- id
			f1 = <-forkchan
			f2 = <-forkchan
			if f1 == f2 {
				fmt.Println("Philosopher ", id, " is eating")
				time.Sleep(time.Second)
				plchan <- false
				prchan <- false
				forkchan <- 0
				forkchan <- 0
				fmt.Println("Philosopher ", id, " is thinking")
			} else {
				fmt.Println("Error: Could not get both forks")
				plchan <- false
				prchan <- false
				forkchan <- 0
				forkchan <- 0
			}

		}
	}
}

func fork(c1 chan int, c2 chan int) {
	fmt.Println("Fork Created")
	var holder int
	var c1res int
	var c2res int
	for {
		c1res = <-c1
		c2res = <-c2
		if holder == 0 && c1res != 0 {
			holder = c1res
			c1res = 0
			c1 <- holder
		} else if holder == 0 && c2res != 0 {
			holder = c2res
			c2res = 0
			c2 <- holder
		}
	}
}
