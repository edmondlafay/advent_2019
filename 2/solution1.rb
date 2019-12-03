def compute(a,b)
  list = File.read("input.txt").split(',').map{|y| y.to_i}
  list[1] = a
  list[2] = b

  i=0
  loop do
    action = list[i]
    case action
    when 1
      list[list[i+3]] = list[list[i+1]] + list[list[i+2]]
      i += 4
    when 2
      list[list[i+3]] = list[list[i+1]] * list[list[i+2]]
      i += 4
    when 99
      return list[0]
    else
      puts "Error: capacity has an invalid value (#{action})"
      return -1
    end
  end
end

puts compute(12,2)
