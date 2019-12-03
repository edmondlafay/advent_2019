package main

import (
  "fmt"
  "strings"
  "strconv"
  "io/ioutil"
  "math"
)

func check(e error) {
  if e != nil {
      panic(e)
  }
}

type Vertex struct {
	x int
	y int
}

func buildCoords(input_raw string) [2]map[Vertex]bool {
  var result [2]map[Vertex]bool
  paths_raw := strings.Split(string(input_raw), "\n")
  for paths_number, path_raw := range paths_raw {
    path := strings.Split(string(path_raw), ",")
    set := map[Vertex]bool {}
    x, y := 0, 0
    for _, order_distance := range path {
      order := order_distance[0]
      distance, err := strconv.Atoi(order_distance[1:])
      check(err)
      for i:=0; i<distance; i++ {
        switch order {
          case 'R': x++;
          case 'L': x--;
          case 'U': y++;
          case 'D': y--;
        }
        set[Vertex{x, y}] = true
      }
    }
    result[paths_number] = set
  }
  return result
}

func pathsIntersections(pathA, pathB map[Vertex]bool) int {
  const UintSize = 32 << (^uint(0) >> 32 & 1)
  var results int = 1<<(UintSize-1) - 1
  for coordsA, _ := range pathA {
    _, present := pathB[coordsA]
    distance := int(math.Abs(float64(coordsA.x))) + int(math.Abs(float64(coordsA.y)))
    if present && distance < results {
      results = distance
    }
  }
  return results
}

func main() {
  input, err := ioutil.ReadFile("input.txt")
  check(err)
  paths := buildCoords(string(input))
  fmt.Printf("solution 1 is %d\n", pathsIntersections(paths[0], paths[1]))
}
