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

func main() {
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  data := ArrayStoArrayI(strings.Split(string(file), ""))
  width, height, length := 25, 6, len(data)
  // data = ArrayStoArrayI(strings.Split("0222112222120000", ""))
  // width, height, length = 2, 2, len(data)

  layer_length := width * height
  count_layer := length/layer_length
  Mzeros, Mones, Mtwos, Mlayer := layer_length, 0, 0, 0
  picture := make([]string, width * height, width * height)
  for i:=0; i<count_layer; i++ {
    layer := data[layer_length*i:i*layer_length+layer_length]
    zeros, ones, twos := 0 ,0, 0
    for x, number := range layer {
      switch number {
        case 0:
          if picture[x]!="#" {picture[x]=" "}
          zeros++
        case 1:
          if picture[x]!=" " {picture[x]="#"}
          ones++
        case 2:
          twos++
      }
    }
    fmt.Printf("layer %d, %d zeroes\n", i, zeros)
    for h:=0; h<height; h++ {
      fmt.Printf("%v\n", layer[h*width:h*width+width])
    }
    fmt.Printf("\n")
    if zeros<Mzeros {
      Mlayer, Mzeros, Mones, Mtwos=i, zeros, ones, twos
    }
  }
  fmt.Printf("Solution 1: %d (layer%d)\n\n", Mones*Mtwos, Mlayer)
  fmt.Printf("Solution 2:\n")
  for h:=0; h<height; h++ {
    fmt.Printf("%v\n", picture[h*width:h*width+width])
  }
  fmt.Printf("\n")
}
