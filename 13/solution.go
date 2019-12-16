package main

import (
  "fmt"
  "math"
  "time"
  "strings"
  "strconv"
  "io/ioutil"
)

const ADD, MULTIPLY, GET, PUT, JUMPIF, JUMPUNLESS, LOWERTHAN, EQUALTO, CHANGEREL, HALT int = 1, 2, 3, 4, 5, 6, 7, 8, 9, 99

func check(e error) {
  if e != nil { panic(e) }
}

func IntAbs(n int) int {
  return int(math.Abs(float64(n)))
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
    //compute_debug_prints(memory, action, modes, params, relative_base)
  }
}

func isBetween(ax, ay, bx, by, cx, cy int) bool {
  // aligned?
  if IntAbs((cy - ay) * (bx - ax) - (cx - ax) * (by - ay)) != 0 { return false }
  //contained?
  dotproduct := (cx - ax) * (bx - ax) + (cy - ay)*(by - ay)
  if dotproduct < 0 { return false }
  if dotproduct > (bx - ax)*(bx - ax) + (by - ay)*(by - ay) { return false }
  return true
}

func detectRicochet (screen [][]byte, ballX, ballY, paddleY, ballXDirection, ballYDirection int, doneBlocks map[[2]int]bool) int {
  projectedX := (ballX + (paddleY-1 - ballY) * ballXDirection)

  for x, y := ballX, ballY ; y<paddleY && x + ballXDirection<len(screen[0]); y += ballYDirection {
    if screen[y][x] == ' ' { screen[y][x] = '.' }
    if (screen[y][x + ballXDirection] == '@' || screen[y][x + ballXDirection] == '#') && !doneBlocks[[2]int{x + ballXDirection, y}] { // change x direction
      if screen[y][x + ballXDirection] == '@' { doneBlocks[[2]int{x + ballXDirection, y}] = true }
      if ballXDirection>0 {ballXDirection=-1} else {ballXDirection=1}
      return detectRicochet(screen, x, y, paddleY, ballXDirection, ballYDirection, doneBlocks)
    }
    if (screen[y + ballYDirection][x] == '@' || screen[y + ballYDirection][x] == '#') && !doneBlocks[[2]int{x, y + ballYDirection}] { // change y direction
      if screen[y + ballYDirection][x] == '@' {doneBlocks[[2]int{x, y + ballYDirection}] = true}
      if ballYDirection>0 {ballYDirection=-1} else {ballYDirection=1}
      return detectRicochet(screen, x, y, paddleY, ballXDirection, ballYDirection, doneBlocks)
    }
    if (screen[y + ballYDirection][x + ballXDirection] == '@' || screen[y + ballYDirection][x + ballXDirection] == '#') && !doneBlocks[[2]int{x + ballXDirection, y + ballYDirection}] { // change y+x direction
      if screen[y + ballYDirection][x + ballXDirection] == '@' { doneBlocks[[2]int{x + ballXDirection, y + ballYDirection}] = true }
      if ballYDirection>0 {ballYDirection=-1} else {ballYDirection=1}
      if ballXDirection>0 {ballXDirection=-1} else {ballXDirection=1}
      return detectRicochet(screen, x, y, paddleY, ballXDirection, ballYDirection, doneBlocks)
    }
    x+=ballXDirection
  }
  return projectedX
}

func joystick (screen [][]byte, input_chan chan int, coordsChan chan [4]int, still_running *bool) {
  var ballX, ballY, paddleX, paddleY, ballXDirection, ballYDirection, projectedX, prevBallX, prevBallY, nextProjectedX int
  ballXDirection, ballYDirection = 1, 1
  for *still_running {
    projectedX = nextProjectedX
    coords := <- coordsChan
    ballX, ballY, paddleX, paddleY = coords[0], coords[1], coords[2], coords[3]

    if ballY == paddleY-1 {
      ballYDirection = -1
      if ballX<paddleX {
        ballXDirection=-1
      } else if ballX>paddleX {
        ballXDirection=1
      } else if (prevBallX > len(screen[0])-3 && ballX > len(screen[0])-4) || (prevBallX < 2 && ballX < 3) {
        ballXDirection=ballXDirection*-1
      }
    }
    if ballY == paddleY-1 && ballYDirection==-1 { 
      nextProjectedX = detectRicochet(screen, ballX, ballY, paddleY, ballXDirection, ballYDirection, make(map[[2]int]bool))
    }
    if ballY < paddleY-1 {
      clear := true && ballY-3>0 && ballX-3>0 && ballX+3<len(screen[0])
      for i:=-3;i<=3 && clear;i++{
        for j:=-3;j<=3 && clear;j++{
          clear = (clear && screen[ballY+i][ballX+j]!='@' && screen[ballY+i][ballX+j]!='#' && screen[ballY+i][ballX+j]!='=')
        }
      }
      if clear {
        projectedX = detectRicochet(screen, ballX, ballY, paddleY, ballXDirection, ballYDirection, make(map[[2]int]bool))
        nextProjectedX = projectedX
      }
    }

    if prevBallX<ballX {ballXDirection=1} else if prevBallX>ballX {ballXDirection=-1}
    if prevBallY<ballY {ballYDirection=1} else if prevBallY>ballY {ballYDirection=-1}
    if ballXDirection==1 {screen[0][1] = '+'} else if ballXDirection==-1 {screen[0][1] = '-'}
    if ballYDirection==1 {screen[1][0] = '+'} else if ballYDirection==-1 {screen[1][0] = '-'}
    prevBallX, prevBallY = ballX, ballY

    for i:=1;i<37;i++ { if screen[paddleY][i]=='?' { screen[paddleY][i] = ' ' } }
    if projectedX>0 && projectedX<len(screen[0]) && screen[paddleY][projectedX]==' ' { screen[paddleY][projectedX] = '?' }
    if projectedX>0 && projectedX < paddleX {
      input_chan <- -1
    } else if projectedX>0 && projectedX > paddleX {
      input_chan <- 1
    } else {
      input_chan <- 0
    }
  }
}

func main() {
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  input_chan, output_chan := make(chan int), make(chan int)
  coordsChan := make(chan [4]int)
  
  screen := make([][]byte, 0)
  for l:=0;l<22;l++{
    screen = append(screen, []byte("                                      "))
  }
  canDisplay, still_running := false, false
  paddleOn, ballOn := true, true
  var x, y, tile_id, score, step, ballX, ballY, paddleX, paddleY int

  go compute(ArrayStoArrayI(strings.Split(string(file), ",")), input_chan, output_chan, 0)
  x, still_running = <- output_chan
  y, still_running = <- output_chan
  tile_id, still_running = <- output_chan

  go joystick(screen, input_chan, coordsChan, &still_running)

  for still_running {
    if x == -1 && y==0 {
      score = tile_id
      canDisplay = true
    } else {
      switch tile_id {
        case 0: 
          screen[y][x] = ' '
          if x==paddleX && y==paddleY { paddleOn = false }
          if x==ballX && y==ballY { ballOn = false }
        case 1: screen[y][x] = '#'
        case 2: screen[y][x] = '@'
        case 3: 
          screen[y][x] = '='
          paddleOn = true
          paddleX, paddleY = x, y
        case 4: 
          screen[y][x] = 'O'
          ballOn = true
          ballX, ballY = x, y
          coordsChan <- [4]int{ballX,ballY,paddleX,paddleY}
      }
    }
    if canDisplay && paddleOn && ballOn {
      fmt.Printf("\n\n\n\n\n\nscore: %d (%d)\n", score, step)
      for line:=0; line < len(screen); line++ { 
        fmt.Println(string(screen[line]))
      }
      
      time.Sleep(10 * time.Millisecond)
    }
    x, still_running = <- output_chan
    y, still_running = <- output_chan
    tile_id, still_running = <- output_chan
    step = step+1
  }
  fmt.Printf("\nfinal score: %d (%d)\n", score, step)
}
