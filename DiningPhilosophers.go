package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type twoWayChannel struct {
	to   chan int
	from chan int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Starting Simulation")
	c12 := make(chan int, 2)
	c13 := make(chan int, 2)
	c14 := make(chan int, 2)
	c15 := make(chan int, 2)

	c21 := make(chan int, 2)
	c23 := make(chan int, 2)
	c24 := make(chan int, 2)
	c25 := make(chan int, 2)

	c31 := make(chan int, 2)
	c32 := make(chan int, 2)
	c34 := make(chan int, 2)
	c35 := make(chan int, 2)

	c41 := make(chan int, 2)
	c42 := make(chan int, 2)
	c43 := make(chan int, 2)
	c45 := make(chan int, 2)

	c51 := make(chan int, 2)
	c52 := make(chan int, 2)
	c53 := make(chan int, 2)
	c54 := make(chan int, 2)

	firstPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c15, c51}, twoWayChannel{c13, c31}, twoWayChannel{c14, c41}, twoWayChannel{c12, c21}}
	secondPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c21, c12}, twoWayChannel{c24, c42}, twoWayChannel{c25, c52}, twoWayChannel{c23, c32}}
	thirdPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c32, c23}, twoWayChannel{c35, c53}, twoWayChannel{c31, c13}, twoWayChannel{c34, c43}}
	fourthPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c43, c34}, twoWayChannel{c41, c14}, twoWayChannel{c42, c24}, twoWayChannel{c45, c54}}
	fifthPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c54, c45}, twoWayChannel{c52, c25}, twoWayChannel{c53, c35}, twoWayChannel{c51, c15}}

	a1 := make(chan bool, 1)
	a2 := make(chan bool, 1)
	b2 := make(chan bool, 1)
	b3 := make(chan bool, 1)
	c3 := make(chan bool, 1)
	c4 := make(chan bool, 1)
	d4 := make(chan bool, 1)
	d5 := make(chan bool, 1)
	e5 := make(chan bool, 1)
	e1 := make(chan bool, 1)

	fmt.Println("Channels Initialised")
	go philo(1, firstPhilosopherChannels, a1, e1)
	go philo(2, secondPhilosopherChannels, a2, b2)
	go philo(3, thirdPhilosopherChannels, b3, c3)
	go philo(4, fourthPhilosopherChannels, c4, d4)
	go philo(5, fifthPhilosopherChannels, d5, e5)

	go fork(a1, a2)
	go fork(b2, b3)
	go fork(c3, c4)
	go fork(d4, d5)
	go fork(e5, e1)
	fmt.Println("Table Initialised")
	time.Sleep(10 * time.Second)
	fmt.Println("Program Terminated")
}

func philo(id int, philosophers [4]twoWayChannel, left chan bool, right chan bool) {
	// The 'id' of the philosopher is only used for debugging purposes, and to clarify who is eating / thinking
	// Does not affect the logic of the code in any way

	for {

		// Roll the dice
		diceRoll := rand.Intn(2048)

		// Tell other philosophers about result
		for i := 0; i < 4; i++ {
			philosophers[i].to <- diceRoll
		}

		// Find out what the other philosophers rolled
		var otherDiceRolls [4]int

		for i := 0; i < 4; i++ {
			otherDiceRolls[i] = <-philosophers[i].from
		}

		// Looks for the highest diceroll
		max := diceRoll

		for i := 0; i < 4; i++ {
			if otherDiceRolls[i] > max {
				max = otherDiceRolls[i]
			}
		}

		// Checks how many philosophers got the highest dice roll

		totalMaxDiceRolls := 0

		if diceRoll == max {
			totalMaxDiceRolls++
		}

		for i := 0; i < 4; i++ {
			if otherDiceRolls[i] == max {
				totalMaxDiceRolls++
			}
		}

		if totalMaxDiceRolls == 1 { // Case where there is found a winner. If no winner is found, the for-loop repeats
			if diceRoll == max { // If the current philosopher is the winner of the dice roll

				// Pick up forks

				fmt.Println("Philosopher ", id, " is eating.")

				// Choose who else should eat
				philosophers[0].to <- -1
				philosophers[1].to <- 1
				philosophers[2].to <- -1
				philosophers[3].to <- -1

				// Put down forks

			} else { // If the current philosopher lost

				for { // Continuously check for message from other philosophers
					// -1 means the current philosopher geats to eat, -2 means
					select {
					case message := <-philosophers[0].from:
						if message == 1 {

							// Pick up forks

							fmt.Println("Philosopher ", id, " is eating.")

							// Put down forks

						} else if message == -1 {
							fmt.Println("Philosopher ", id, " is thinking.")
						}
						break
					case message := <-philosophers[1].from:
						if message == 1 {

							// Pick up forks

							fmt.Println("Philosopher ", id, " is eating.")

							// Put down forks

						} else if message == -1 {
							fmt.Println("Philosopher ", id, " is thinking.")
						}
						break
					case message := <-philosophers[2].from:
						if message == 1 {

							// Pick up forks

							fmt.Println("Philosopher ", id, " is eating.")

							// Put down forks

						} else if message == -1 {
							fmt.Println("Philosopher ", id, " is thinking.")
						}
						break
					case message := <-philosophers[3].from:
						if message == 1 {

							// Pick up forks

							fmt.Println("Philosopher ", id, " is eating.")

							// Put down forks

						} else if message == -1 {
							fmt.Println("Philosopher ", id, " is thinking.")
						}
						break
					}
				}
			}
		}
	}
}

func fork(c1 chan bool, c2 chan bool) {
	fmt.Println("Fork Created")
	select {
	case <-c1:
		fmt.Println("got msg from c1")
		fmt.Println("waiting to put down fork (c1)")
		select {
		case <-c2:
			fmt.Println("ERROR FORK ALREADY HELD")
			os.Exit(3)
		}
		if false == <-c1 {
			break
		}
	case <-c2:
		fmt.Println("got msg from c2")
		fmt.Println("waiting to put down fork (c2)")
		select {
		case <-c1:
			fmt.Println("ERROR FORK ALREADY HELD")
			os.Exit(3)
		}
		if false == <-c2 {
			break
		}
	}
}
