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

func calculate_fuel_for_mass(mass int) int {
  return (mass/3)-2
}

func calculate_fuel_for_mass_and_fuel(mass int) int {
  var fuel = calculate_fuel_for_mass(mass)
  if fuel > 0 {
    return fuel + calculate_fuel_for_mass_and_fuel(fuel)
  }
  return 0
}

func main() {
  input, err := ioutil.ReadFile("input.txt")
  check(err)
  list_of_masses := ArrayStoArrayI(strings.Split(string(input), "\n"))
  var fuel_for_mass, fuel_for_mass_and_fuel = 0, 0
  for _, mass := range list_of_masses {
    fuel_for_mass = fuel_for_mass + calculate_fuel_for_mass(mass)
    fuel_for_mass_and_fuel = fuel_for_mass_and_fuel + calculate_fuel_for_mass_and_fuel(mass)
  }
  fmt.Printf("solution 1 is %d\n", fuel_for_mass)
  fmt.Printf("solution 2 is %d\n", fuel_for_mass_and_fuel)
}
