package main

import (
  "fmt"
  "strings"
  "strconv"
  "io/ioutil"
)

func check(e error) {if e != nil { panic(e) }}

type Ressource struct {
  name string
  needs map[string]int
  yield int
  have int
  made int
}

func produce(elements map[string]Ressource, ressourceName string, needs int) {
  ressource := elements[ressourceName]

  if ressourceName == "ORE" {
    ressource.made += needs
    elements[ressourceName] = ressource
    return
  }

  //fmt.Printf("%5s needs:%2d, have: %2d", ressourceName, needs, ressource.have)

  if ressource.have >= needs {
    ressource.have -= needs
    elements[ressourceName] = ressource
    //fmt.Printf("\n")
    return
  }
  needs -= ressource.have
  ressource.have = 0
  production := 1
  if ressource.yield >= needs {
    ressource.have = ressource.yield - needs
  } else {
    production = (needs / ressource.yield)
    if needs % ressource.yield > 0 { production += 1 }
    ressource.have = ressource.yield * production - needs
  }
  //fmt.Printf(", yields:%2d => launch:%2d cycles, %2d remaining\n", ressource.yield, production, ressource.have)
  ressource.made += ressource.yield * production
  elements[ressourceName] = ressource

  for sub_resource_name, sub_resource_needs := range ressource.needs {
    produce(elements, sub_resource_name, sub_resource_needs*production)
  }
  return
}

func main() {
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  lines := strings.Split(string(file), "\n")
  // lines = strings.Split("171 ORE => 8 CNZTR\n7 ZLQW, 3 BMBT, 9 XCVML, 26 XMNCP, 1 WPTQ, 2 MZWV, 1 RJRHP => 4 PLWSL\n114 ORE => 4 BHXH\n14 VRPVC => 6 BMBT\n6 BHXH, 18 KTJDG, 12 WPTQ, 7 PLWSL, 31 FHTLT, 37 ZDVW => 1 FUEL\n6 WPTQ, 2 BMBT, 8 ZLQW, 18 KTJDG, 1 XMNCP, 6 MZWV, 1 RJRHP => 6 FHTLT\n15 XDBXC, 2 LTCX, 1 VRPVC => 6 ZLQW\n13 WPTQ, 10 LTCX, 3 RJRHP, 14 XMNCP, 2 MZWV, 1 ZLQW => 1 ZDVW\n5 BMBT => 4 WPTQ\n189 ORE => 9 KTJDG\n1 MZWV, 17 XDBXC, 3 XCVML => 2 XMNCP\n12 VRPVC, 27 CNZTR => 2 XDBXC\n15 KTJDG, 12 BHXH => 5 XCVML\n3 BHXH, 2 VRPVC => 7 MZWV\n121 ORE => 7 VRPVC\n7 XCVML => 6 RJRHP\n5 BHXH, 4 VRPVC => 5 LTCX", "\n") // 2210736
  // lines = strings.Split("2 VPVL, 7 FWMGM, 2 CXFTF, 11 MNCFX => 1 STKFG\n17 NVRVD, 3 JNWZP => 8 VPVL\n53 STKFG, 6 MNCFX, 46 VJHF, 81 HVMC, 68 CXFTF, 25 GNMV => 1 FUEL\n22 VJHF, 37 MNCFX => 5 FWMGM\n139 ORE => 4 NVRVD\n144 ORE => 7 JNWZP\n5 MNCFX, 7 RFSQX, 2 FWMGM, 2 VPVL, 19 CXFTF => 3 HVMC\n5 VJHF, 7 MNCFX, 9 VPVL, 37 CXFTF => 6 GNMV\n145 ORE => 6 MNCFX\n1 NVRVD => 8 CXFTF\n1 VJHF, 6 MNCFX => 4 RFSQX\n176 ORE => 6 VJHF", "\n") // 180697
  // lines = strings.Split("157 ORE => 5 NZVS\n165 ORE => 6 DCFZ\n44 XJWVT, 5 KHKGT, 1 QDVJ, 29 NZVS, 9 GPVTF, 48 HKGWZ => 1 FUEL\n12 HKGWZ, 1 GPVTF, 8 PSHF => 9 QDVJ\n179 ORE => 7 PSHF\n177 ORE => 5 HKGWZ\n7 DCFZ, 7 PSHF => 2 XJWVT\n165 ORE => 2 GPVTF\n3 DCFZ, 7 NZVS, 5 HKGWZ, 10 PSHF => 8 KHKGT", "\n") // 13312
  // lines = strings.Split("9 ORE => 2 A\n8 ORE => 3 B\n7 ORE => 5 C\n3 A, 4 B => 1 AB\n5 B, 7 C => 1 BC\n4 C, 1 A => 1 CA\n2 AB, 3 BC, 4 CA => 1 FUEL", "\n") // 165
  // lines = strings.Split("10 ORE => 10 A\n1 ORE => 1 B\n7 A, 1 B => 1 C\n7 A, 1 C => 1 D\n7 A, 1 D => 1 E\n7 A, 1 E => 1 FUEL", "\n") // 31
  // lines = strings.Split("10 ORE => 1 AAAA\n10 ORE => 1 BBBB\n10 AAAA, 10 BBBB => 1 FUEL", "\n") // 200
  // lines = strings.Split("1 ORE => 10 AAAA\n1 ORE => 10 BBBB\n7 AAAA, 7 BBBB => 1 FUEL", "\n") // 2
  // lines = strings.Split("10 ORE => 1 FUEL", "\n") // 10
  elements := make(map[string]Ressource)
  for _, line := range lines {
    in_out := strings.Split(line, " => ")
    q_out := strings.Split(in_out[1], " ")
    name, yield_s := q_out[1], q_out[0]
    yield, err := strconv.Atoi(yield_s)
    check(err)
    ressource := Ressource{name: name, yield: yield, needs: make(map[string]int)}
    ins := strings.Split(in_out[0], ", ")
    for _, in := range ins {
      q_int := strings.Split(in, " ")
      name, needed_s := q_int[1], q_int[0]
      needed, err := strconv.Atoi(needed_s)
      check(err)
      ressource.needs[name] = needed
    }
    elements[ressource.name] = ressource
  }
  elements["ORE"] = Ressource{name: "ORE"}
  
  minOreForFuelElements := make(map[string]Ressource)
  for element_name, element_desc := range elements {
    element := Ressource{name: element_name, yield: element_desc.yield, needs: make(map[string]int)}
    for sub_element_name, sub_elements_needed := range element_desc.needs {
      element.needs[sub_element_name] = sub_elements_needed
    }
    minOreForFuelElements[element_name] = element
  }
  produce(minOreForFuelElements, "FUEL", 1)
  result1 := minOreForFuelElements["ORE"].made
  fmt.Printf("RESULT1: %d\n", result1)

  result2 := 4436981
  produce(elements, "FUEL", result2)
  oreMade := elements["ORE"].made
  fmt.Printf("RESULT2: %d (%d)\n", result2, oreMade)
}
