package main

import (
  "fmt"
  "math"
  "strings"
  "strconv"
  "io/ioutil"
)

const ADD int = 1

func check(e error) {
  if e != nil {
      panic(e)
  }
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

func intPermutations(array []int) [][]int {
  var res [][]int
  if len(array)==1 {
    return append(res, []int{array[0]})
  } else {
    sub_permutation := intPermutations(array[1:])
    for j:=0; j<len(sub_permutation); j++ {
      for i:=0; i<len(sub_permutation[j]); i++ {
        tmp := make([]int, len(array), len(array))
        copy(tmp, sub_permutation[j])
        copy(tmp[i+1:], tmp[i:])
        tmp[i] = array[0]
        res = append(res, tmp)
      }
      res = append(res, append(sub_permutation[j], array[0]))
    }
  }
  return res
}

func get_params(memory map[int]int, i int, j int, modes int, relative int) int {
  mode := (modes/int(math.Pow(10, float64(j))))%10
  switch mode {
    case 0: return memory[i+j+1]
    case 1: return i+j+1
    case 2: return relative+memory[i+j+1]
    default: panic(fmt.Sprintf("Error: compute has an invalid mode %d (modes: %d).\n", mode, modes))
  }
}

func compute_debug_prints(memory map[int]int, action int, modes int, params []int) {
  fmt.Printf("action %d,  mode %d, params: ", action, modes)
  for _, param := range params { fmt.Printf("%d(%d) ", param, memory[param]) }
  fmt.Println()
  for k:=0;k<len(memory); k++ { fmt.Printf("%d ",memory[k]) }
  fmt.Printf("\n\n")
}

func compute(input_data []int, input chan int, output chan int, number int) {
  memory := make(map[int]int)
  for memory_position, memory_order := range input_data {
    memory[memory_position] = memory_order
  }
  relative_base, i := 0, 0
  for i < len(memory) {
    action, modes := memory[i]%100 , memory[i]/100
    var params []int
    switch action {
      case ADD:
        for j:=0; j < 3; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        memory[params[2]] = memory[params[1]] + memory[params[0]]
        i += 4
      case 2:
        for j:=0; j < 3; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        memory[params[2]] = memory[params[1]] * memory[params[0]]
        i += 4
      case 3:
        params = append(params, get_params(memory, i, 0, modes, relative_base))
        memory[params[0]] = <-input
        i += 2
      case 4:
        params = append(params, get_params(memory, i, 0, modes, relative_base))
        output <- memory[params[0]]
        i += 2
      case 5:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        if memory[params[0]]!=0 { i=memory[params[1]] } else { i+=3 }
      case 6:
        for j:=0; j < 2; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        if memory[params[0]]==0 { i=memory[params[1]] } else { i+=3 }
      case 7:
        for j:=0; j < 3; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        if memory[params[0]]<memory[params[1]] { memory[params[2]]=1 } else { memory[params[2]]=0 }
        i += 4
      case 8:
        for j:=0; j < 3; j++ {
          params = append(params, get_params(memory, i, j, modes, relative_base))
        }
        if memory[params[0]]==memory[params[1]] { memory[params[2]]=1 } else { memory[params[2]]=0 }
        i += 4
      case 9:
        params = append(params, get_params(memory, i, 0, modes, relative_base))
        relative_base += memory[params[0]]
        i += 2
      case 99:
        close(output)
        return
      default:
        panic(fmt.Sprintf("Error: compute has an invalid action %d.\n", action))
    }
    //compute_debug_prints(memory map[int]int, action int, modes int, params []int)
  }
}

func tests() {
  string_inputs := [][3]string{
    {"3,3,1108,-1,8,3,4,3,99", "8", "1"},
    {"3,3,1108,-1,8,3,4,3,99", "9", "0"},
    {"3,3,1107,-1,8,3,4,3,99", "0", "1"},
    {"3,3,1107,-1,8,3,4,3,99", "9", "0"},
    {"3,9,8,9,10,9,4,9,99,-1,8", "8", "1"},
    {"3,9,8,9,10,9,4,9,99,-1,8", "0", "0"},
    {"3,9,7,9,10,9,4,9,99,-1,8", "0", "1"},
    {"3,9,7,9,10,9,4,9,99,-1,8", "9", "0"},
    {"3,3,1105,-1,9,1101,0,0,12,4,12,99,1", "0", "0"},
    {"3,3,1105,-1,9,1101,0,0,12,4,12,99,1", "5", "1"},
    {"104,1125899906842624,99", "", "1125899906842624"},
    {"3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", "0", "0"},
    {"3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9", "5", "1"},
    {"1102,34915192,34915192,7,4,7,99,0", "", "1219070632396864"},
    {"109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99", "", "109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99"},
    {"3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", "3", "999"},
    {"3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", "8", "1000"},
    {"3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99", "9", "1001"},
    {"3,225,1,225,6,6,1100,1,238,225,104,0,1102,91,92,225,1102,85,13,225,1,47,17,224,101,-176,224,224,4,224,1002,223,8,223,1001,224,7,224,1,223,224,223,1102,79,43,225,1102,91,79,225,1101,94,61,225,1002,99,42,224,1001,224,-1890,224,4,224,1002,223,8,223,1001,224,6,224,1,224,223,223,102,77,52,224,1001,224,-4697,224,4,224,102,8,223,223,1001,224,7,224,1,224,223,223,1101,45,47,225,1001,43,93,224,1001,224,-172,224,4,224,102,8,223,223,1001,224,1,224,1,224,223,223,1102,53,88,225,1101,64,75,225,2,14,129,224,101,-5888,224,224,4,224,102,8,223,223,101,6,224,224,1,223,224,223,101,60,126,224,101,-148,224,224,4,224,1002,223,8,223,1001,224,2,224,1,224,223,223,1102,82,56,224,1001,224,-4592,224,4,224,1002,223,8,223,101,4,224,224,1,224,223,223,1101,22,82,224,1001,224,-104,224,4,224,1002,223,8,223,101,4,224,224,1,223,224,223,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,8,226,677,224,102,2,223,223,1005,224,329,1001,223,1,223,1007,226,226,224,1002,223,2,223,1006,224,344,101,1,223,223,108,226,226,224,1002,223,2,223,1006,224,359,1001,223,1,223,107,226,677,224,102,2,223,223,1006,224,374,101,1,223,223,8,677,677,224,102,2,223,223,1006,224,389,1001,223,1,223,1008,226,677,224,1002,223,2,223,1006,224,404,101,1,223,223,7,677,677,224,1002,223,2,223,1005,224,419,101,1,223,223,1108,226,677,224,1002,223,2,223,1005,224,434,101,1,223,223,1108,226,226,224,102,2,223,223,1005,224,449,1001,223,1,223,107,226,226,224,102,2,223,223,1005,224,464,101,1,223,223,1007,677,677,224,102,2,223,223,1006,224,479,101,1,223,223,1007,226,677,224,102,2,223,223,1005,224,494,1001,223,1,223,1008,226,226,224,1002,223,2,223,1005,224,509,1001,223,1,223,1108,677,226,224,1002,223,2,223,1006,224,524,1001,223,1,223,108,677,677,224,1002,223,2,223,1005,224,539,101,1,223,223,108,226,677,224,1002,223,2,223,1005,224,554,101,1,223,223,1008,677,677,224,1002,223,2,223,1006,224,569,1001,223,1,223,1107,677,677,224,102,2,223,223,1005,224,584,1001,223,1,223,7,677,226,224,102,2,223,223,1005,224,599,1001,223,1,223,8,677,226,224,1002,223,2,223,1005,224,614,1001,223,1,223,7,226,677,224,1002,223,2,223,1006,224,629,101,1,223,223,1107,677,226,224,1002,223,2,223,1005,224,644,1001,223,1,223,1107,226,677,224,102,2,223,223,1006,224,659,1001,223,1,223,107,677,677,224,1002,223,2,223,1005,224,674,101,1,223,223,4,223,99,226", "5", "9386583"},
    {"3,225,1,225,6,6,1100,1,238,225,104,0,1102,91,92,225,1102,85,13,225,1,47,17,224,101,-176,224,224,4,224,1002,223,8,223,1001,224,7,224,1,223,224,223,1102,79,43,225,1102,91,79,225,1101,94,61,225,1002,99,42,224,1001,224,-1890,224,4,224,1002,223,8,223,1001,224,6,224,1,224,223,223,102,77,52,224,1001,224,-4697,224,4,224,102,8,223,223,1001,224,7,224,1,224,223,223,1101,45,47,225,1001,43,93,224,1001,224,-172,224,4,224,102,8,223,223,1001,224,1,224,1,224,223,223,1102,53,88,225,1101,64,75,225,2,14,129,224,101,-5888,224,224,4,224,102,8,223,223,101,6,224,224,1,223,224,223,101,60,126,224,101,-148,224,224,4,224,1002,223,8,223,1001,224,2,224,1,224,223,223,1102,82,56,224,1001,224,-4592,224,4,224,1002,223,8,223,101,4,224,224,1,224,223,223,1101,22,82,224,1001,224,-104,224,4,224,1002,223,8,223,101,4,224,224,1,223,224,223,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,8,226,677,224,102,2,223,223,1005,224,329,1001,223,1,223,1007,226,226,224,1002,223,2,223,1006,224,344,101,1,223,223,108,226,226,224,1002,223,2,223,1006,224,359,1001,223,1,223,107,226,677,224,102,2,223,223,1006,224,374,101,1,223,223,8,677,677,224,102,2,223,223,1006,224,389,1001,223,1,223,1008,226,677,224,1002,223,2,223,1006,224,404,101,1,223,223,7,677,677,224,1002,223,2,223,1005,224,419,101,1,223,223,1108,226,677,224,1002,223,2,223,1005,224,434,101,1,223,223,1108,226,226,224,102,2,223,223,1005,224,449,1001,223,1,223,107,226,226,224,102,2,223,223,1005,224,464,101,1,223,223,1007,677,677,224,102,2,223,223,1006,224,479,101,1,223,223,1007,226,677,224,102,2,223,223,1005,224,494,1001,223,1,223,1008,226,226,224,1002,223,2,223,1005,224,509,1001,223,1,223,1108,677,226,224,1002,223,2,223,1006,224,524,1001,223,1,223,108,677,677,224,1002,223,2,223,1005,224,539,101,1,223,223,108,226,677,224,1002,223,2,223,1005,224,554,101,1,223,223,1008,677,677,224,1002,223,2,223,1006,224,569,1001,223,1,223,1107,677,677,224,102,2,223,223,1005,224,584,1001,223,1,223,7,677,226,224,102,2,223,223,1005,224,599,1001,223,1,223,8,677,226,224,1002,223,2,223,1005,224,614,1001,223,1,223,7,226,677,224,1002,223,2,223,1006,224,629,101,1,223,223,1107,677,226,224,1002,223,2,223,1005,224,644,1001,223,1,223,1107,226,677,224,102,2,223,223,1006,224,659,1001,223,1,223,107,677,677,224,1002,223,2,223,1005,224,674,101,1,223,223,4,223,99,226", "1", "0,0,0,0,0,0,0,0,0,16489636"},
    {"1102,34463338,34463338,63,1007,63,34463338,63,1005,63,53,1101,0,3,1000,109,988,209,12,9,1000,209,6,209,3,203,0,1008,1000,1,63,1005,63,65,1008,1000,2,63,1005,63,902,1008,1000,0,63,1005,63,58,4,25,104,0,99,4,0,104,0,99,4,17,104,0,99,0,0,1102,32,1,1019,1101,0,500,1023,1101,0,636,1025,1102,36,1,1010,1101,0,29,1013,1102,864,1,1029,1102,21,1,1000,1102,1,507,1022,1102,1,28,1011,1102,38,1,1008,1101,0,35,1004,1101,25,0,1018,1102,24,1,1005,1102,30,1,1009,1102,1,869,1028,1101,0,37,1007,1102,1,23,1017,1102,1,20,1015,1102,1,22,1003,1101,0,39,1001,1102,1,31,1012,1101,701,0,1026,1101,0,641,1024,1101,0,34,1016,1102,1,0,1020,1102,698,1,1027,1102,33,1,1002,1102,26,1,1006,1101,0,1,1021,1101,0,27,1014,109,12,21101,40,0,0,1008,1012,40,63,1005,63,203,4,187,1105,1,207,1001,64,1,64,1002,64,2,64,109,-11,1207,7,37,63,1005,63,223,1105,1,229,4,213,1001,64,1,64,1002,64,2,64,109,14,1206,5,247,4,235,1001,64,1,64,1105,1,247,1002,64,2,64,109,-2,1207,-4,31,63,1005,63,269,4,253,1001,64,1,64,1105,1,269,1002,64,2,64,109,-6,1208,-5,35,63,1005,63,289,1001,64,1,64,1106,0,291,4,275,1002,64,2,64,109,9,21108,41,39,-1,1005,1015,311,1001,64,1,64,1105,1,313,4,297,1002,64,2,64,109,-5,2101,0,-9,63,1008,63,33,63,1005,63,339,4,319,1001,64,1,64,1106,0,339,1002,64,2,64,1205,10,351,4,343,1106,0,355,1001,64,1,64,1002,64,2,64,109,-18,2108,35,9,63,1005,63,375,1001,64,1,64,1105,1,377,4,361,1002,64,2,64,109,18,1205,9,389,1105,1,395,4,383,1001,64,1,64,1002,64,2,64,109,7,21107,42,41,-8,1005,1010,415,1001,64,1,64,1106,0,417,4,401,1002,64,2,64,109,-12,2102,1,0,63,1008,63,29,63,1005,63,437,1106,0,443,4,423,1001,64,1,64,1002,64,2,64,109,3,1208,0,30,63,1005,63,461,4,449,1105,1,465,1001,64,1,64,1002,64,2,64,109,5,1202,-5,1,63,1008,63,31,63,1005,63,489,1001,64,1,64,1106,0,491,4,471,1002,64,2,64,109,15,2105,1,-6,1001,64,1,64,1106,0,509,4,497,1002,64,2,64,109,-10,1206,2,525,1001,64,1,64,1106,0,527,4,515,1002,64,2,64,109,-18,1202,0,1,63,1008,63,39,63,1005,63,553,4,533,1001,64,1,64,1106,0,553,1002,64,2,64,109,1,2107,21,1,63,1005,63,571,4,559,1105,1,575,1001,64,1,64,1002,64,2,64,109,7,2102,1,-8,63,1008,63,39,63,1005,63,601,4,581,1001,64,1,64,1105,1,601,1002,64,2,64,109,2,1201,-7,0,63,1008,63,35,63,1005,63,623,4,607,1106,0,627,1001,64,1,64,1002,64,2,64,109,20,2105,1,-7,4,633,1106,0,645,1001,64,1,64,1002,64,2,64,109,-16,21107,43,44,-4,1005,1011,663,4,651,1105,1,667,1001,64,1,64,1002,64,2,64,109,-11,2107,36,0,63,1005,63,687,1001,64,1,64,1106,0,689,4,673,1002,64,2,64,109,19,2106,0,4,1106,0,707,4,695,1001,64,1,64,1002,64,2,64,109,-14,21108,44,44,6,1005,1015,725,4,713,1105,1,729,1001,64,1,64,1002,64,2,64,109,1,1201,-6,0,63,1008,63,36,63,1005,63,749,1106,0,755,4,735,1001,64,1,64,1002,64,2,64,109,-1,21101,45,0,10,1008,1019,42,63,1005,63,775,1105,1,781,4,761,1001,64,1,64,1002,64,2,64,109,16,21102,46,1,-7,1008,1018,44,63,1005,63,801,1105,1,807,4,787,1001,64,1,64,1002,64,2,64,109,-3,21102,47,1,-4,1008,1018,47,63,1005,63,833,4,813,1001,64,1,64,1105,1,833,1002,64,2,64,109,-14,2108,38,0,63,1005,63,851,4,839,1105,1,855,1001,64,1,64,1002,64,2,64,109,17,2106,0,3,4,861,1106,0,873,1001,64,1,64,1002,64,2,64,109,-31,2101,0,10,63,1008,63,36,63,1005,63,897,1001,64,1,64,1106,0,899,4,879,4,64,99,21101,0,27,1,21101,0,913,0,1106,0,920,21201,1,53612,1,204,1,99,109,3,1207,-2,3,63,1005,63,962,21201,-2,-1,1,21102,940,1,0,1106,0,920,21202,1,1,-1,21201,-2,-3,1,21101,955,0,0,1106,0,920,22201,1,-1,-2,1105,1,966,21201,-2,0,-2,109,-3,2106,0,0","2","74917"},
    {"1102,34463338,34463338,63,1007,63,34463338,63,1005,63,53,1101,0,3,1000,109,988,209,12,9,1000,209,6,209,3,203,0,1008,1000,1,63,1005,63,65,1008,1000,2,63,1005,63,902,1008,1000,0,63,1005,63,58,4,25,104,0,99,4,0,104,0,99,4,17,104,0,99,0,0,1102,32,1,1019,1101,0,500,1023,1101,0,636,1025,1102,36,1,1010,1101,0,29,1013,1102,864,1,1029,1102,21,1,1000,1102,1,507,1022,1102,1,28,1011,1102,38,1,1008,1101,0,35,1004,1101,25,0,1018,1102,24,1,1005,1102,30,1,1009,1102,1,869,1028,1101,0,37,1007,1102,1,23,1017,1102,1,20,1015,1102,1,22,1003,1101,0,39,1001,1102,1,31,1012,1101,701,0,1026,1101,0,641,1024,1101,0,34,1016,1102,1,0,1020,1102,698,1,1027,1102,33,1,1002,1102,26,1,1006,1101,0,1,1021,1101,0,27,1014,109,12,21101,40,0,0,1008,1012,40,63,1005,63,203,4,187,1105,1,207,1001,64,1,64,1002,64,2,64,109,-11,1207,7,37,63,1005,63,223,1105,1,229,4,213,1001,64,1,64,1002,64,2,64,109,14,1206,5,247,4,235,1001,64,1,64,1105,1,247,1002,64,2,64,109,-2,1207,-4,31,63,1005,63,269,4,253,1001,64,1,64,1105,1,269,1002,64,2,64,109,-6,1208,-5,35,63,1005,63,289,1001,64,1,64,1106,0,291,4,275,1002,64,2,64,109,9,21108,41,39,-1,1005,1015,311,1001,64,1,64,1105,1,313,4,297,1002,64,2,64,109,-5,2101,0,-9,63,1008,63,33,63,1005,63,339,4,319,1001,64,1,64,1106,0,339,1002,64,2,64,1205,10,351,4,343,1106,0,355,1001,64,1,64,1002,64,2,64,109,-18,2108,35,9,63,1005,63,375,1001,64,1,64,1105,1,377,4,361,1002,64,2,64,109,18,1205,9,389,1105,1,395,4,383,1001,64,1,64,1002,64,2,64,109,7,21107,42,41,-8,1005,1010,415,1001,64,1,64,1106,0,417,4,401,1002,64,2,64,109,-12,2102,1,0,63,1008,63,29,63,1005,63,437,1106,0,443,4,423,1001,64,1,64,1002,64,2,64,109,3,1208,0,30,63,1005,63,461,4,449,1105,1,465,1001,64,1,64,1002,64,2,64,109,5,1202,-5,1,63,1008,63,31,63,1005,63,489,1001,64,1,64,1106,0,491,4,471,1002,64,2,64,109,15,2105,1,-6,1001,64,1,64,1106,0,509,4,497,1002,64,2,64,109,-10,1206,2,525,1001,64,1,64,1106,0,527,4,515,1002,64,2,64,109,-18,1202,0,1,63,1008,63,39,63,1005,63,553,4,533,1001,64,1,64,1106,0,553,1002,64,2,64,109,1,2107,21,1,63,1005,63,571,4,559,1105,1,575,1001,64,1,64,1002,64,2,64,109,7,2102,1,-8,63,1008,63,39,63,1005,63,601,4,581,1001,64,1,64,1105,1,601,1002,64,2,64,109,2,1201,-7,0,63,1008,63,35,63,1005,63,623,4,607,1106,0,627,1001,64,1,64,1002,64,2,64,109,20,2105,1,-7,4,633,1106,0,645,1001,64,1,64,1002,64,2,64,109,-16,21107,43,44,-4,1005,1011,663,4,651,1105,1,667,1001,64,1,64,1002,64,2,64,109,-11,2107,36,0,63,1005,63,687,1001,64,1,64,1106,0,689,4,673,1002,64,2,64,109,19,2106,0,4,1106,0,707,4,695,1001,64,1,64,1002,64,2,64,109,-14,21108,44,44,6,1005,1015,725,4,713,1105,1,729,1001,64,1,64,1002,64,2,64,109,1,1201,-6,0,63,1008,63,36,63,1005,63,749,1106,0,755,4,735,1001,64,1,64,1002,64,2,64,109,-1,21101,45,0,10,1008,1019,42,63,1005,63,775,1105,1,781,4,761,1001,64,1,64,1002,64,2,64,109,16,21102,46,1,-7,1008,1018,44,63,1005,63,801,1105,1,807,4,787,1001,64,1,64,1002,64,2,64,109,-3,21102,47,1,-4,1008,1018,47,63,1005,63,833,4,813,1001,64,1,64,1105,1,833,1002,64,2,64,109,-14,2108,38,0,63,1005,63,851,4,839,1105,1,855,1001,64,1,64,1002,64,2,64,109,17,2106,0,3,4,861,1106,0,873,1001,64,1,64,1002,64,2,64,109,-31,2101,0,10,63,1008,63,36,63,1005,63,897,1001,64,1,64,1106,0,899,4,879,4,64,99,21101,0,27,1,21101,0,913,0,1106,0,920,21201,1,53612,1,204,1,99,109,3,1207,-2,3,63,1005,63,962,21201,-2,-1,1,21102,940,1,0,1106,0,920,21202,1,1,-1,21201,-2,-3,1,21101,955,0,0,1106,0,920,22201,1,-1,-2,1105,1,966,21201,-2,0,-2,109,-3,2106,0,0","1","2377080455"},
  }
  for i, string_input := range string_inputs {
    test_data := ArrayStoArrayI(strings.Split(string_input[0], ","))
    test_inputs := ArrayStoArrayI(strings.Split(string_input[1], ","))
    test_outputs := ArrayStoArrayI(strings.Split(string_input[2], ","))
    input_chan := make(chan int, len(test_inputs))
    output_chan := make(chan int)
    go compute(test_data, input_chan, output_chan, i)
    for _, input := range test_inputs { input_chan <- input }
    var more_outputs bool
    var output int
    for _, test_output := range test_outputs {
      output, more_outputs = <-output_chan
      if output != test_output { panic(fmt.Sprintf("Error: test %d, expecting %d got %d\n", i, test_output, output)) }
      if !more_outputs { panic(fmt.Sprintf("Error: test %d, expecting %d got nothing\n", i, test_output)) }
    }
    if more_outputs {
      _, more_outputs = <-output_chan
      if more_outputs {
        panic(fmt.Sprintf("Error: test %d, expecting nothing got %d\n", i, <-output_chan)) 
      }
    }
  }
}

func main() {
  tests()
  file, err := ioutil.ReadFile("input.txt")
  check(err)
  input_chan := make(chan int)
  output_chan := make(chan int)
  go compute(ArrayStoArrayI(strings.Split(string(file), ",")), input_chan, output_chan, 0)
  input_chan <- 2
  var result int
  for output := range output_chan {
    result = output
  }
  fmt.Printf("RESULT: %d\n", result)
}