package main

import (
  "fmt"
  "math"
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

func intPermutations(array []int) [][]int {
  var res [][]int
  if len(array)<=1 {
    return append(res, array)
  } else {
    sub_permutation := intPermutations(array[1:])
    for j:=0; j<len(sub_permutation); j++ {
      for i:=0; i<len(sub_permutation[j]); i++ {
        tmp := append(sub_permutation[j], 0)
        copy(tmp[i+1:], tmp[i:])
        tmp[i] = array[0]
        res = append(res, tmp)
      }
      res = append(res, append(sub_permutation[j], array[0]))
    }
  }
  return res
}

func Perm(a []rune, f func([]rune)) {
  perm(a, f, 0)
}

func perm(a []rune, f func([]rune), i int) {
  if i > len(a) {
      f(a)
      return
  }
  perm(a, f, i+1)
  for j := i + 1; j < len(a); j++ {
      a[i], a[j] = a[j], a[i]
      perm(a, f, i+1)
      a[i], a[j] = a[j], a[i]
  }
}

func reverse_int(n int) int {
  new_int := 0
  for n > 0 {
      remainder := n % 10
      new_int *= 10
      new_int += remainder 
      n /= 10
  }
  return new_int 
}

func get_params(list []int, i int, j int, modes int) int {
  mode := modes/int(math.Pow(10, float64(j)))%int(math.Pow(10, float64(j+1)))
  switch mode {
    case 0:
      return list[list[i+j+1]]
    case 1:
      return list[i+j+1]
    default:
      fmt.Printf("Error: compute has an invalid mode %d.\n", mode)
      return -1
  }
}

func compute(input chan int, output chan int, result chan int) {
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  memory := ArrayStoArrayI(strings.Split(string(file), ","))
  var last_output int
  var i int = 0
  for i < len(memory) {
    action := memory[i]%100
    modes := memory[i]/100
    var params []int
    switch action {
      case 1:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        memory[memory[i+3]] = params[1] + params[0]
        i += 4
      case 2:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        memory[memory[i+3]] = params[1] * params[0]
        i += 4
      case 3:
        memory[memory[i+1]] = <-input
        i += 2
      case 4:
        output <- memory[memory[i+1]]
        last_output = memory[memory[i+1]]
        i += 2
      case 5:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        if params[0]!=0 { i=params[1] } else { i+=3 }
      case 6:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        if params[0]==0 { i=params[1] } else { i+=3 }
      case 7:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        if params[0]<params[1] { memory[memory[i+3]]=1 } else { memory[memory[i+3]]=0 }
        i += 4
      case 8:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes))
        }
        if params[0]==params[1] { memory[memory[i+3]]=1 } else { memory[memory[i+3]]=0 }
        i += 4
      case 99:
        result <- last_output
      default:
        fmt.Printf("Error: compute has an invalid action %d.\n", action)
        output <- -1
    }
  }
}

func main() {
  var chans = [3]chan int {
    make(chan int),
    make(chan int),
    make(chan int),
  }
  go compute(chans[0], chans[1], chans[2])
  chans[0] <- 5
  fmt.Printf("solution is %d\n", <-chans[1])
}
