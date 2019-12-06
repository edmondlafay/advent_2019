#!/usr/bin/env ruby

def getParams(list, pos, mode)
  case mode
    when 0; return list[list[pos+1]]
    when 1; return list[pos+1]
    else; puts "Error: action has an invalid mode (#{action})\n#{list[0..pos]}"
  end
end

def compute()
  list = File.read("input.txt").split(',').map{|y| y.to_i}
  #list = '3,9,8,9,10,9,4,9,99,-1,8'.split(',').map{|y| y.to_i}
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
        print "Inputs? "; list[list[i+1]] = gets.chomp.to_i
        i += 2
      when 4
        puts "PRINT #{list[list[i+1]]}"
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
      when 99; exit
      else; puts "Error: capacity has an invalid value (#{action})\n#{list[0..i]}"; exit
    end
  end
end

compute()