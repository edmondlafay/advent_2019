package main

import (
  "fmt"
  "sync"
  "math"
  "sort"
  "time"
  "strings"
  "io/ioutil"
)

func check(e error) {if e != nil { panic(e) }}

type Safe struct { lock sync.Mutex }
func (t *Safe) Lock() { t.lock.Lock() }
func (t *Safe) Unlock() { t.lock.Unlock() }

type Asteroid struct {
  x int
  y int
  asteroid_seen map[float64]Asteroid
  angleScore float64
}

func angleScore(xa, ya, xb, yb int) float64 {
  if xb==xa {
    if yb-ya>0 { return float64(10000) }
    return float64(-10000)
  }
  if xb-xa>0 { return float64(yb-ya)/float64(xb-xa) - 1000 }
  return float64(yb-ya)/float64(xb-xa)
}

func (safeLock *Safe) checkAsteroidBetween(lines []string, a *Asteroid, b Asteroid, angleScore float64, output_chan chan int) {
  incx, incy := func(i int)int { return i+1 }, func(i int)int { return i+1 }
  checkx, checky := func(i, j int) bool { return i<=j }, func(i, j int) bool { return i<=j }
  if b.y<a.y { incy, checky = func(i int) int { return i-1 }, func(i, j int) bool { return i>=j } }
  if b.x<a.x { incx, checkx = func(i int) int { return i-1 }, func(i, j int) bool { return i>=j } }

  // vertical alignment
  if a.x==b.x {
    xc, yc := a.x, incy(a.y)
    for checky(yc,b.y) {
      if lines[yc][xc]=='#'{
        safeLock.Lock()
        a.asteroid_seen[angleScore] = Asteroid{x: xc, y: yc, angleScore: angleScore}
        safeLock.Unlock()
        output_chan <- 1
        return
      }
      yc = incy(yc)
    }
  }

  xc := incx(a.x)
  x1,x2,y1,y2:=a.x,b.x,a.y,b.y
  if b.x>a.x {x1,x2,y1,y2=b.x,a.x,b.y,a.y}
  angle := float64(y2-y1)/float64(x2-x1)
  miss := float64(x2*y1-x1*y2)/float64(x2-x1)
  for checkx(xc, b.x) {
    yc := angle * float64(xc) + miss // y=ax+b
    yc = math.Round(yc*1000)/1000 // round because of float
    if yc==float64(int(yc)) && checky(int(yc), b.y){
      if lines[int(yc)][xc]=='#' {
        safeLock.Lock()
        a.asteroid_seen[angleScore] = Asteroid{x: xc, y: int(yc), angleScore: angleScore}
        safeLock.Unlock()
        output_chan <- 1
        return
      }
    }
    xc = incx(xc)
  }
}

func (safeLock *Safe) checkVisibility(lines []string, asteroids []Asteroid, a *Asteroid, ch chan Asteroid) {
  routines := 0
  output_chan := make(chan int)
  for _, asteroid := range asteroids {
    if a.x!=asteroid.x || a.y!=asteroid.y {
      angleScore := angleScore(a.x, a.y, asteroid.x, asteroid.y)
      safeLock.Lock()
      _, alreadyDone := a.asteroid_seen[angleScore]
      safeLock.Unlock()
      if !alreadyDone {
        go safeLock.checkAsteroidBetween(lines, a, asteroid, angleScore, output_chan)
        routines++
      }
    }
  }
  for i:=0;i<routines;i++{ <- output_chan }
  ch <- *a
}

func main() {
  start := time.Now()
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  lines := strings.Split(string(file), "\n")
  //lines := strings.Split(".#..##.###...#######\n##.############..##.\n.#.######.########.#\n.###.#######.####.#.\n#####.##.#.##.###.##\n..#####..#.#########\n####################\n#.####....###.#.#.##\n##.#################\n#####.##.###..####..\n..######..##.#######\n####.##.####...##..#\n.#####..#.######.###\n##...#.##########...\n#.##########.#######\n.####.#.###.###.#.##\n....##.##.###..#####\n.#.#.###########.###\n#.#.#.#####.####.###\n###.##.####.##.#..##", "\n")

  // build asteroids array
  asteroids := []Asteroid{}
  for ya, line := range lines {
    for xa, data_point := range line {
      if data_point=='#' {
        asteroids = append(asteroids, Asteroid{x: xa, y: ya, asteroid_seen: make(map[float64]Asteroid)})
      }
    }
  }

  routines := 0
  output_chan := make(chan Asteroid)
  for i, _ := range asteroids {
    safeLock := Safe{lock: sync.Mutex{}}
    go safeLock.checkVisibility(lines, asteroids, &asteroids[i], output_chan)
    routines++
  }

  max := 0
  var base Asteroid
  for i:=0;i<routines;i++{ 
    asteroid := <-output_chan
    if max < len(asteroid.asteroid_seen) {
      max = len(asteroid.asteroid_seen)
      base = asteroid
    }
  } 
  fmt.Printf("Solution 1: asteroid %d,%d sees %d other\n", base.x, base.y, max)

  // sort asteroid by angle
  lineOfsight := []Asteroid{}
  for _, asteroid := range base.asteroid_seen {
    lineOfsight = append(lineOfsight, asteroid)
  }
  sort.Slice(lineOfsight, func(i, j int) bool {
		return lineOfsight[i].angleScore < lineOfsight[j].angleScore
  })

  fmt.Printf("Solution 2: %d\n", lineOfsight[198].x*100+lineOfsight[198].y)
  time := time.Now()
  fmt.Printf("time: %d\n", time.Sub(start))

  // build visuals
  t := make([][]string, 0)
  for _, line := range lines { t = append(t, strings.Split(line, "")) }
  t[base.y][base.x] = "@"
  letters := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890abcdefghijklmnopqrstuvwxyz"
  for i, asteroid := range lineOfsight {
    if i<len(letters) { t[asteroid.y][asteroid.x] = string(letters[i]) }
  }
  for i, line := range t { fmt.Printf("%d%v\n", i%10, line) }
}
