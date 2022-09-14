package main

import (
	"fmt"
	"math/rand"
	"time"
)

type twoWayChannel struct {
	to   chan int
	from chan int
}

func main() {

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

	firstPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c15, c51}, twoWayChannel{c12, c21}, twoWayChannel{c13, c31}, twoWayChannel{c14, c41}}
	secondPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c25, c52}, twoWayChannel{c21, c12}, twoWayChannel{c23, c32}, twoWayChannel{c24, c42}}
	thirdPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c35, c53}, twoWayChannel{c31, c13}, twoWayChannel{c32, c23}, twoWayChannel{c34, c43}}
	fourthPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c45, c54}, twoWayChannel{c41, c14}, twoWayChannel{c42, c24}, twoWayChannel{c43, c34}}
	fifthPhilosopherChannels := [4]twoWayChannel{twoWayChannel{c51, c15}, twoWayChannel{c52, c25}, twoWayChannel{c53, c35}, twoWayChannel{c54, c45}}

	// cf1 := make(chan int, 1)
	// cf2 := make(chan int, 1)
	// cf3 := make(chan int, 1)
	// cf4 := make(chan int, 1)
	// cf5 := make(chan int, 1)

	fmt.Println("Channels Initialised")
	go philo(1, firstPhilosopherChannels)
	go philo(3, thirdPhilosopherChannels)
	go philo(2, secondPhilosopherChannels)
	go philo(4, fourthPhilosopherChannels)
	go philo(5, fifthPhilosopherChannels)

	// go fork(cf1, cf2)
	// go fork(cf2, cf3)
	// go fork(cf3, cf4)
	// go fork(cf4, cf5)
	// go fork(cf5, cf1)
	fmt.Println("Table Initialised")
	time.Sleep(10 * time.Second)
	fmt.Println("Program Terminated")
}

func philo(id int, philosophers [4]twoWayChannel) {
	isEating := false

	for {

		diceRoll := rand.Intn(2048)
		for i := 0; i < 4; i++ {
			philosophers[i].to <- diceRoll
		}
		// fmt.Println(diceRoll, " ", id)

		var otherDiceRolls [4]int

		for i := 0; i < 4; i++ {
			otherDiceRolls[i] = <-philosophers[i].from
		}

		max := diceRoll

		for i := 0; i < 4; i++ {
			if otherDiceRolls[i] > max {
				max = otherDiceRolls[i]
			}
		}

		totalMaxDiceRolls := 0

		if diceRoll == max {
			totalMaxDiceRolls++
		}

		for i := 0; i < 4; i++ {
			if otherDiceRolls[i] == max {
				totalMaxDiceRolls++
			}
		}

		if totalMaxDiceRolls == 1 {
			if diceRoll == max {
				isEating = true
			}
			break
		}
	}

	fmt.Println(isEating, " ", id)

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
