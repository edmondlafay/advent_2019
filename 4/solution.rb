#!/usr/bin/env ruby

require 'set'
input = "134564-585159"

start = Time.now
i = 0
('134566'..'579999').each do |test|
  if test == test.chars.sort.join and
  (test[0] == test[1] or test[1]==test[2] or test[2]==test[3] or test[3]==test[4] or test[4]==test[5])
    i+=1
  end
end
finish = Time.now
puts "solution1: #{i} (in #{finish - start} seconds)"

start = Time.now
i = 0
('134566'..'578899').each do |test|
  if test == test.chars.sort.join and
  ((test[0] == test[1] and test[1]!=test[2]) or 
  (test[1]==test[2] and test[1]!=test[0] and test[3]!=test[2]) or 
  (test[2]==test[3] and test[2]!=test[1] and test[3]!=test[4])or 
  (test[3]==test[4] and test[3]!=test[2] and test[4]!=test[5]) or 
  (test[4]==test[5] and test[4]!=test[3]))
    i+=1
  end
end
finish = Time.now
puts "solution2: #{i} (in #{finish - start} seconds)"
