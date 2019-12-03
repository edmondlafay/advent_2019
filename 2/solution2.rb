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

res = 0
j = 0
loop do
  i = 0
  loop do
    res = compute(i, j)
    puts("#{i},#{j} : #{res}") 
    if res==19690720 or i>98
      break
    end
    i=i+1
  end
  if res==19690720
    puts("RESULT #{i}#{j}") 
    break
  end
  j=j+1
  break if j>=99
end
