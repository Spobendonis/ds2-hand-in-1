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

	go fork(1.5, a1, a2)
	go fork(2.5, b2, b3)
	go fork(3.5, c3, c4)
	go fork(4.5, d4, d5)
	go fork(0.5, e5, e1)
	fmt.Println("Table Initialised")
	time.Sleep(10 * time.Second)
	fmt.Println("Program Terminated")
}

func philo(id int, philosophers [4]twoWayChannel, left chan bool, right chan bool) {
	// The 'id' of the philosopher is only used for debugging purposes, and to clarify who is eating / thinking
	// Does not affect the logic of the code in any way

	isEating := false
	timesEating := 0
	timesThinking := 0

	philosophersReady := [4]bool{false, false, false, false}
	var otherDiceRolls [4]int

	for {

		for i := 0; i < 4; i++ {
			philosophers[i].to <- 0
		}

		for {

			if philosophersReady[0] && philosophersReady[1] && philosophersReady[2] && philosophersReady[3] {
				// fmt.Println("All people ready,", id)
				// time.Sleep(1)
				break
			}

			if !philosophersReady[0] {
				num := <-philosophers[0].from
				if 0 == num {
					philosophersReady[0] = true
				} else {
					fmt.Println("Recieved ", num, " instead of 0")
				}
			}
			if !philosophersReady[1] {
				num := <-philosophers[1].from
				if 0 == num {
					philosophersReady[1] = true
				} else {
					fmt.Println("Recieved ", num, " instead of 0")
				}
			}
			if !philosophersReady[2] {
				num := <-philosophers[2].from
				if 0 == num {
					philosophersReady[2] = true
				} else {
					fmt.Println("Recieved ", num, " instead of 0")
				}
			}
			if !philosophersReady[3] {
				num := <-philosophers[3].from
				if 0 == num {
					philosophersReady[3] = true
				} else {
					fmt.Println("Recieved ", num, " instead of 0")
				}
			}
		}

		philosophersReady = [4]bool{false, false, false, false}

		// wait for confirmation

		// Roll the dice
		diceRoll := rand.Intn(2048) + 10

		// Tell other philosophers about result
		for i := 0; i < 4; i++ {
			philosophers[i].to <- diceRoll
		}

		// fmt.Println(id, "sent dicerolls")

		// Find out what the other philosophers rolled

		for i := 0; i < 4; i++ {
			otherDiceRolls[i] = <-philosophers[i].from
		}

		// fmt.Println(id, "recieved dicerolls")

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

				left <- true
				right <- true

				// Choose who else should eat
				philosophers[0].to <- -1
				philosophers[1].to <- 1
				philosophers[2].to <- -1
				philosophers[3].to <- -1

				if !isEating {
					timesEating++
					fmt.Println("Philosopher ", id, " is eating", timesEating, " times")
					isEating = true
				}

				left <- false
				right <- false

			} else { // If the current philosopher lost

			inner:
				for { // Continuously check for message from other philosophers
					// -1 means the current philosopher geats to eat, -2 means

					select {
					case message := <-philosophers[0].from:
						// fmt.Println(message, id)
						if message == 1 {

							left <- true
							right <- true

							if !isEating {
								timesEating++
								fmt.Println("Philosopher ", id, " is eating", timesEating, " times")
								isEating = true
							}

							left <- false
							right <- false

						} else if message == 0 {
							philosophersReady[0] = true
						} else {

							if isEating {
								timesThinking++
								fmt.Println("Philosopher ", id, " is thinking", timesThinking, " times")
								isEating = false
							}
						}
						break inner
					case message := <-philosophers[1].from:
						// fmt.Println(message, id)
						if message == 1 {

							left <- true
							right <- true

							if !isEating {
								timesEating++
								fmt.Println("Philosopher ", id, " is eating", timesEating, " times")
								isEating = true
							}

							left <- false
							right <- false

						} else if message == 0 {
							philosophersReady[1] = true
						} else {

							if isEating {
								timesThinking++
								fmt.Println("Philosopher ", id, " is thinking", timesThinking, " times")
								isEating = false
							}
						}
						break inner
					case message := <-philosophers[2].from:
						// fmt.Println(message, id)
						if message == 1 {

							left <- true
							right <- true

							if !isEating {
								timesEating++
								fmt.Println("Philosopher ", id, " is eating", timesEating, " times")
								isEating = true
							}

							left <- false
							right <- false

						} else if message == 0 {
							philosophersReady[2] = true
						} else {

							if isEating {
								timesThinking++
								fmt.Println("Philosopher ", id, " is thinking", timesThinking, " times")
								isEating = false
							}
						}
						break inner
					case message := <-philosophers[3].from:
						// fmt.Println(message, id)

						if message == 1 {

							left <- true
							right <- true

							if !isEating {
								timesEating++
								fmt.Println("Philosopher ", id, " is eating", timesEating, " times")
								isEating = true
							}

							left <- false
							right <- false

						} else if message == 0 {
							philosophersReady[3] = true
						} else {

							if isEating {
								timesThinking++
								fmt.Println("Philosopher ", id, " is thinking", timesThinking, " times")
								isEating = false
							}
						}
						break inner
					}

				}

			}
		}
	}

}

func fork(id float32, c1 chan bool, c2 chan bool) {
	beingHeld := false
	// fmt.Println("Fork Created")
	for {
		select {
		case message := <-c1:

			if beingHeld && !message {
				beingHeld = false
			} else if !beingHeld && message {
				beingHeld = true
			} else {
				fmt.Println("ERROR FORK ALREADY HELD: ", id)
				os.Exit(3)
			}

			//fmt.Println("got msg from c1")
			//fmt.Println("waiting to put down fork (c1)")
			// select {
			// case <-c2:
			// 	fmt.Println("ERROR FORK ALREADY HELD")
			// 	os.Exit(3)
			// case <-
			//}
		case message := <-c2:

			if beingHeld && !message {
				beingHeld = false
			} else if !beingHeld && message {
				beingHeld = true
			} else {
				fmt.Println("ERROR FORK ALREADY HELD: ", id)
				os.Exit(3)
			}

			//fmt.Println("got msg from c2")
			//fmt.Println("waiting to put down fork (c2)")
			// select {
			// case <-c1:
			// 	fmt.Println("ERROR FORK ALREADY HELD")
			// 	os.Exit(3)
			// }
			// if false == <-c2 {
			// 	break
			// }
		}
	}
}
