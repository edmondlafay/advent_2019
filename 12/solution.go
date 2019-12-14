package main

import (
  "fmt"
  "sync"
  "math"
  "reflect"
)

func check(e error) {if e != nil { panic(e) }}

type Safe struct { lock sync.Mutex }
func (t *Safe) Lock() { t.lock.Lock() }
func (t *Safe) Unlock() { t.lock.Unlock() }

type Vertex struct {
  x int
  y int
  z int
}

type Moon struct {
  position Vertex
  velocity Vertex
  name string
}

func (v *Vertex) getField(field string) int {
  r := reflect.ValueOf(v)
  f := reflect.Indirect(r).FieldByName(field)
  return int(f.Int())
}

func IntAbs(n int) int {
  return int(math.Abs(float64(n)))
}

func GCD(a, b int) int {
  for b != 0 { b, a = a % b, b }
  return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
  result := a * b / GCD(a, b)
  for i := 0; i < len(integers); i++ {
    result = LCM(result, integers[i])
  }
  return result
}

func updateGravityDimension(moon1, moon2 *Moon, dimension string) (int, int){
  if moon1.position.getField(dimension)>moon2.position.getField(dimension) {
    return moon1.velocity.getField(dimension)-1, moon2.velocity.getField(dimension)+1
  }
  if moon1.position.getField(dimension)<moon2.position.getField(dimension) {
    return moon1.velocity.getField(dimension)+1, moon2.velocity.getField(dimension)-1
  }
  return moon1.velocity.getField(dimension), moon2.velocity.getField(dimension)
}

func gravity(moon1, moon2 *Moon) {
  moon1.velocity.x, moon2.velocity.x = updateGravityDimension(moon1, moon2, "x")
  moon1.velocity.y, moon2.velocity.y = updateGravityDimension(moon1, moon2, "y")
  moon1.velocity.z, moon2.velocity.z = updateGravityDimension(moon1, moon2, "z")
}

func moonEnergy (moon Moon) int{
  return (IntAbs(moon.position.x) + IntAbs(moon.position.y) + IntAbs(moon.position.z)) * (IntAbs(moon.velocity.x) + IntAbs(moon.velocity.y) + IntAbs(moon.velocity.z))
}

func systemEnergy (moons []Moon) int {
  total_energy := 0
  for _, moon := range moons {
    total_energy += moonEnergy(moon)
  }
  return total_energy
}

func printMoons(moons []Moon) {
  for _, moon := range moons {
    fmt.Printf("%10s position %5d,%5d,%5d velocity %3d,%3d,%3d\n", moon.name, 
      moon.position.x, moon.position.y, moon.position.z, 
      moon.velocity.x, moon.velocity.y, moon.velocity.z)
  }
  fmt.Printf("\n")
}

func main() {
  moons := [4]Moon{
    Moon{name: "Io",       position: Vertex{x:  5, y: 13, z:-3}},
    Moon{name: "Europa",   position: Vertex{x: 18, y: -7, z:13}},
    Moon{name: "Ganymede", position: Vertex{x: 16, y:  3, z: 4}},
    Moon{name: "Callisto", position: Vertex{x:  0, y:  8, z: 8}},
  }
  // moons = [4]Moon{
  //   Moon{name: "Io",       position: Vertex{x: -1, y:   0, z: 2}},
  //   Moon{name: "Europa",   position: Vertex{x:  2, y: -10, z:-7}},
  //   Moon{name: "Ganymede", position: Vertex{x:  4, y:  -8, z: 8}},
  //   Moon{name: "Callisto", position: Vertex{x:  3, y:   5, z:-1}},
  // }
  // moons = [4]Moon{
  //   Moon{name: "Io",       position: Vertex{x: -8, y:  10, z: 0}},
  //   Moon{name: "Europa",   position: Vertex{x:  5, y:   5, z:10}},
  //   Moon{name: "Ganymede", position: Vertex{x:  2, y:  -7, z: 3}},
  //   Moon{name: "Callisto", position: Vertex{x:  9, y:  -8, z:-3}},
  // }

  step := 0 
  max := 4767828*2
  fmt.Printf("step: %d/%d\n", step, max)
  printMoons(moons[:])

  x_init,y_init,z_init := 0,0,0
  for step<max {
    step++
    veloXzero,veloYzero,veloZzero := true, true, true
    for moon1:=0;moon1<len(moons);moon1++ {
      for moon2:=moon1+1;moon2<len(moons);moon2++ {
        gravity(&moons[moon1], &moons[moon2])
      }
      moons[moon1].position.x += moons[moon1].velocity.x
      moons[moon1].position.y += moons[moon1].velocity.y
      moons[moon1].position.z += moons[moon1].velocity.z
      veloXzero= moons[moon1].velocity.x==0 && veloXzero
      veloYzero= moons[moon1].velocity.y==0 && veloYzero
      veloZzero= moons[moon1].velocity.z==0 && veloZzero
    }
    
    if step==1000 { fmt.Printf("Solution 1: %d\n", systemEnergy(moons[:]))}
    if veloXzero && x_init==0 { x_init = step }
    if veloYzero && y_init==0 { y_init = step }
    if veloZzero && z_init==0 { z_init = step }

    if x_init*y_init*z_init>0 && step>=1000 {
      fmt.Printf("%d %d %d\n", x_init,y_init,z_init)
      break
    }
  }

  fmt.Printf("step: %d/%d\n", step, max)
  printMoons(moons[:])

  if step<max {
    fmt.Printf("result2: %d\n", LCM(x_init*2,y_init*2,z_init*2))
  }
}
