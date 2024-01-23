/*
 * CMPU-377 - Parallel Programming
 * Fall 2022
 * 
 * sort-pump.go
 * <your name>
 *
 * This sort pump is a concurrent implementation of bubble sort.
 */

package main

import "fmt"
import "math/rand"
import "time"

/* 
 * Generate and send sequence of random numbers [0..99]  to channel 'ch'.
 * (to be run as a goroutine/process)
 *
 * Parameters:
 * - ch: output channel for random numbers
 * - kill: output channel to indicate end of random number stream
 */
func randNums(ch chan<- int, kill chan<- bool) {
  s1:= rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
  for i := 0; i<100;i++{
    num := r1.Intn(100)
    ch <- num
    }
  kill <- true
    }

/*
 * Pump the values from channel 'src' to channel 'dst',
 * holding back the larger of the last two numbers read.
 * Terminate when kill signal received / pass along to next cell
 * (run as concurrent goroutines / as many as numbers being sorted)
 *
 * Parameters:
 * - src:     read numbers from this channel
 * - killIn:  check this channel for kill signal
 * - dst:     write smaller number to this channel
 * - killOut: pass along kill signal on this channel
 */
func cell(src <-chan int, killIn <-chan bool, dst chan<- int, killOut chan<- bool){
  num2 := <-src
  for{
    select{
    case num := <-src:{
      if num < num2 {
        dst <- num
      } else {
        dst <- num2
        num2 = num
      }
    }
  case boo := <- killIn:{
    dst <- num2
    killOut <- boo
    break
  }
    }
  }
}

/*
 * print out numbers in sorted order
 * (to be run as a goroutine/process)
 * 
 * Parameters:
 * - src:     read numbers to be printed from this channel
 * - killIn:  check this channel for kill signal
 * - killOut: pass along kill signal on this channel
 */
func printSorted(src <-chan int, killIn <-chan bool, killO chan<- bool){
  for{
  select{
  case num := <-src:{
    fmt.Print(num, "\n")
  }
  case boo := <-killIn:{
    killO <- boo
    break
    }
  }
  }
}

/*
 * Create the sort pump pipeline consisting of:
 * - randNums() process: to generate random numbers into pipeline
 * - 100 cell() processes: to pump numbers through from one cell to next
 * - printSorted(): last process in pipeline connected to last cell
 */
 func sortPump() {
  src := make(chan int)
  kill := make(chan bool)
  go randNums(src, kill)

  for i:=1;i<100;i++{
    nextSrc := make(chan int)
    nextKill := make(chan bool)
    go cell(src, kill, nextSrc, nextKill)
		src = nextSrc
		kill = nextKill
  }
  finalKill := make(chan bool)
	go printSorted(src, kill, finalKill)
 
	<-finalKill
  }


func main() {
  sortPump()
}
