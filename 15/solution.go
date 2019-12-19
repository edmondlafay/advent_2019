package main

import (
  "fmt"
  "math"
  "strings"
  "strconv"
  "io/ioutil"
)

const ADD, MULTIPLY, GET, PUT, JUMPIF, JUMPUNLESS, LOWERTHAN, EQUALTO, CHANGEREL, HALT int = 1, 2, 3, 4, 5, 6, 7, 8, 9, 99

func check(e error) {
  if e != nil { panic(e) }
}

func ArrayStoArrayI(t []string) []int {
  var t2 = []int{}
  for _, i := range t {
    if len(i)>0 {
      j, err := strconv.Atoi(i)
      check(err)
      t2 = append(t2, j)
    }
  }
  return t2
}

func replaceAtIndex(in string, r rune, i int) string {
  out := []rune(in)
  out[i] = r
  return string(out)
}

func get_params(memory map[int]int, i int, length int, modes int, relative_base int) []int {
  params := make([]int, 0)
  for j:=0; j < length-1; j++ {
    mode := (modes/int(math.Pow10(j)))%10
    absolute, direct, relative := 0, 1, 2
    switch mode {
      case absolute: params = append(params, memory[i+j+1])
      case direct:   params = append(params, i+j+1)
      case relative: params = append(params, relative_base+memory[i+j+1])
      default: panic(fmt.Sprintf("Error: compute has an invalid mode %d (modes: %d).\n", mode, modes))
    }
  }
  return params
}

func compute(input_data []int, input chan int, output chan int, number int) {
  memory := make(map[int]int)
  for memory_position, memory_order := range input_data {
    memory[memory_position] = memory_order
  }
  relative_base, i, length := 0, 0, 0
  for i < len(memory) {
    action, modes := memory[i]%100 , memory[i]/100
    var params []int
    switch action {
      case ADD:
        length = 4
        params = get_params(memory, i, length, modes, relative_base)
        memory[params[2]] = memory[params[1]] + memory[params[0]]
      case MULTIPLY:
        length = 4
        params = get_params(memory, i, length, modes, relative_base)
        memory[params[2]] = memory[params[1]] * memory[params[0]]
      case GET:
        length = 2
        params = get_params(memory, i, length, modes, relative_base)
        memory[params[0]] = <-input
      case PUT:
        length = 2
        params = get_params(memory, i, length, modes, relative_base)
        output <- memory[params[0]]
      case JUMPIF:
        length = 3
        params = get_params(memory, i, length, modes, relative_base)
        if memory[params[0]]!=0 { 
          i=memory[params[1]]
          continue
        }
      case JUMPUNLESS:
        length = 3
        params = get_params(memory, i, length, modes, relative_base)
        if memory[params[0]]==0 {
          i=memory[params[1]]
          continue
        }
      case LOWERTHAN:
        length = 4
        params = get_params(memory, i, length, modes, relative_base)
        if memory[params[0]]<memory[params[1]] { memory[params[2]]=1 } else { memory[params[2]]=0 }
      case EQUALTO:
        length = 4
        params = get_params(memory, i, length, modes, relative_base)
        if memory[params[0]]==memory[params[1]] { memory[params[2]]=1 } else { memory[params[2]]=0 }
      case CHANGEREL:
        length = 2
        params = get_params(memory, i, length, modes, relative_base)
        relative_base += memory[params[0]]
      case HALT:
        close(output)
        return
      default:
        panic(fmt.Sprintf("Error: compute has an invalid action %d.\n", action))
    }
    i += length
  }
}

func printCanvas(canvas map[[2]int]string, start_x, start_y, end_x, end_y, minX, minY, maxX, maxY, step int) {
  fmt.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n%d\n", step)
  for y:=minY; y< maxY; y++ {
    for x:=minX; x< maxX; x++ {
      if x==start_x && y==start_y {
        fmt.Printf("£")
      } else if x==end_x && y==end_y {
        fmt.Printf("$")
      } else if known_position, ok := canvas[[2]int{x,y}]; ok {
        fmt.Printf("%s", known_position)
      } else {
        fmt.Printf(" ")
      }
    }
    fmt.Printf("\n")
  }
}

func chooseDirection (canvas map[[2]int]string, input_chan chan int, x, y int) (next_x, next_y, next_movement int) {
  switch canvas[[2]int{x, y}] {
    case "<":
      if canvas[[2]int{x, y+1}] != "#" {
        next_x, next_y = x, y+1
        next_movement = 2
        return
      } else {
        canvas[[2]int{x, y}] = "^"
        return chooseDirection(canvas, input_chan, x, y)
      }
    case "^":
      if canvas[[2]int{x-1, y}] != "#" {
        next_x, next_y = x-1, y
        next_movement = 3
        return
      } else {
        canvas[[2]int{x, y}] = ">"
        return chooseDirection(canvas, input_chan, x, y)
      }
    case ">":
      if canvas[[2]int{x, y-1}] != "#" {
        next_x, next_y = x, y-1
        next_movement = 1
        return
      } else {
        canvas[[2]int{x, y}] = "v"
        return chooseDirection(canvas, input_chan, x, y)
      }
    case "v":
      if canvas[[2]int{x+1, y}] != "#" {
        next_x, next_y = x+1, y
        next_movement = 4
        return
      } else {
        canvas[[2]int{x, y}] = "<"
        return chooseDirection(canvas, input_chan, x, y)
      }
  }
  return x, y, 1
}

func propagate(canvas map[[2]int]string, x, y int) int{
  time := 0
  propagated := false
  if canvas[[2]int{x-1,y}] == "." {
    canvas[[2]int{x-1,y}] = "O"
    propagated = true
    tmp := propagate(canvas, x-1,y)
    if tmp > time { time = tmp }
  }
  if canvas[[2]int{x+1,y}] == "." {
    canvas[[2]int{x+1,y}] = "O"
    propagated = true
    tmp := propagate(canvas, x+1,y)
    if tmp > time { time = tmp }
  }
  if canvas[[2]int{x,y-1}] == "." {
    canvas[[2]int{x,y-1}] = "O"
    propagated = true
    tmp := propagate(canvas, x,y-1)
    if tmp > time { time = tmp }
  }
  if canvas[[2]int{x,y+1}] == "." {
    canvas[[2]int{x,y+1}] = "O"
    propagated = true
    tmp := propagate(canvas, x, y+1)
    if tmp > time { time = tmp } 
  }
  if propagated {
    time++
  }
  return time
}

func main() {
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  canvas := make(map[[2]int]string)

  input_chan := make(chan int)
  output_chan := make(chan int)
  go compute(ArrayStoArrayI(strings.Split(string(file), ",")), input_chan, output_chan, 0)

  start_x,start_y,minX,minY,maxX,maxY,step:=40,25,0,0,80,50,0
  x,y := start_x,start_y
  movements := "^v<>"
  canvas[[2]int{x,y}] = "<"
  // go down
  next_x, next_y := x+1, y
  next_movement := 4

  keep_going := true
  for keep_going {
    input_chan <- next_movement
    if step%100==0 { printCanvas(canvas, start_x, start_y, start_x-20, start_y-12, minX, minY, maxX, maxY, step) }
    output := <- output_chan
    switch output {
      case 0:
        canvas[[2]int{next_x, next_y}] = "#"
      case 1:
        canvas[[2]int{next_x, next_y}] = string(movements[next_movement-1])
        canvas[[2]int{x, y}] = "."
        x, y=next_x, next_y
      case 2:
        canvas[[2]int{next_x, next_y}] = string(movements[next_movement-1])
        canvas[[2]int{x, y}] = "."
        x, y=next_x, next_y
        // canvas[[2]int{next_x, next_y}] = "$"
        // keep_going = false
    }
    if step==3000 { keep_going = false }
    
    if keep_going {
      next_x, next_y, next_movement = chooseDirection (canvas, input_chan, x, y)
      step++
    }
  }

  printCanvas(canvas, start_x, start_y, start_x-20, start_y-12, minX, minY, maxX, maxY, step)
  canvas[[2]int{x, y}] = "."
  canvas[[2]int{start_x, start_y}] = "."
  canvas[[2]int{start_x-20, start_y-12}] = "O"

  res2 := propagate(canvas, start_x-20, start_y-12)
  printCanvas(canvas, start_x, start_y, start_x-20, start_y-12, minX, minY, maxX, maxY, step)
  
  fmt.Printf("RESULT 2: %d\n\n", res2)
}
