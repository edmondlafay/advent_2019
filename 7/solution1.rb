#!/usr/bin/env ruby

def getParams(list, pos, mode)
  case mode
    when 0; return list[list[pos+1]]
    when 1; return list[pos+1]
    else; puts "Error: action has an invalid mode (#{action})\n#{list[0..pos]}"
  end
end

def compute(inputs)
  list = File.read("input.txt").split(',').map{|y| y.to_i}
  #list = '3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0'.split(',').map{|y| y.to_i}
  read_intput = 0
  outputs = []
  i=0
  loop do
    action = list[i]%100
    modes = (list[i]/100).to_s.split('').reverse.map{|y| y.to_i}
    params = []
    case action
      when 1
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        list[list[i+3]] = params[1] + params[0]
        i += 4
      when 2
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        list[list[i+3]] = params[1] * params[0]
        i += 4
      when 3
        if read_intput < inputs.length
          input = inputs[read_intput]
          read_intput+=1
          puts "Inputs? #{input.to_s}"
          list[list[i+1]] = input
        else
          print "Inputs?"
          list[list[i+1]] = gets.chomp.to_i
        end
        i += 2
      when 4
        puts "PRINT #{list[list[i+1]]}"
        outputs.push(list[list[i+1]])
        i += 2
      when 5
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        params[0]!=0 ? i = params[1] : i += 3
      when 6
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        params[0]==0 ? i = params[1] : i += 3
      when 7
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        params[0] < params[1] ? list[list[i+3]] = 1 : list[list[i+3]] = 0
        i += 4
      when 8
        2.times{ |j| params.push(getParams(list, i+j, modes[j] || 0))}
        params[0] == params[1] ? list[list[i+3]] = 1 : list[list[i+3]] = 0
        i += 4
      when 99; return outputs
      else; puts "Error: capacity has an invalid value (#{action})\n#{list[0..i]}"; exit
    end
  end
end

amplis = %w(a b c d e)
max = 0
[0,1,2,3,4].permutation.each do |inputs|
  last_output = 0
  puts "inputs #{inputs}"
  inputs.each_with_index do |input, i|
    puts "AMPLI #{amplis[i]}"
    last_output = compute([input, last_output]).first
  end
  puts "output #{last_output}"
  max = [max, last_output].max
end
puts "RESULT: #{max}"

t = [1,2,3]
for _, i := range t {
  j, err := strconv.Atoi(i)
  check(err)
  t2 = append(t2, j)
}











