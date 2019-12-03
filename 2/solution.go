package main

import (
  "fmt"
  "strings"
  "strconv"
  "io/ioutil"
)

func check(e error) {
  if e != nil {
      panic(e)
  }
}

func ArrayStoArrayI(t []string) []int {
  var t2 = []int{}
  for _, i := range t {
    j, err := strconv.Atoi(i)
    check(err)
    t2 = append(t2, j)
  }
  return t2
}

func compute(a int, b int) int {
  input, err := ioutil.ReadFile("input.txt")
  check(err)
  memory := ArrayStoArrayI(strings.Split(string(input), ","))
  memory[1] = a
  memory[2] = b

  var i int = 0
  for i < len(memory) {
    var instruction int = memory[i]
    switch instruction {
    case 1:
      memory[memory[i+3]] = memory[memory[i+1]] + memory[memory[i+2]]
      i=i+4
    case 2:
      memory[memory[i+3]] = memory[memory[i+1]] * memory[memory[i+2]]
      i=i+4
    case 99:
      return memory[0]
    default:
      fmt.Printf("Error: compute has an invalid instruction %d.\n", instruction)
      return -1
    }
  }
  return -1
}

func reverseCompute(desired_result int, result_chanel chan int) {
  for verb:=0; verb<100 ; verb++ {
    for noun:=0;noun<100;noun++ {
      go computeResultInputOutput(noun, verb, desired_result, result_chanel)
    }
  }
}

func computeResultInputOutput(noun, verb, expected_result int, result_chanel chan int) {
  if compute(noun, verb)==expected_result {
    result_chanel <- noun
    result_chanel <- verb
  }
}

func main() {
  fmt.Printf("solution 1 is %d\n", compute(12, 2))
  c := make(chan int)
  reverseCompute(19690720, c)
  fmt.Printf("solution 2 is %d%d\n", <-c, <-c)
}
