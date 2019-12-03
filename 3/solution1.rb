require 'set'
test_input = "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7"
test output = 159

def buildCoords(input)
  paths = []
  input.split("\n").each do |path_raw|
    coords = Set[]
    x, y = [0, 0]
    path_raw.split(',').each do |dir_dist|
      order = dir_dist[0]
      distance = dir_dist[1..-1].to_i
      distance.times do |i|
        case order
          when 'R' then x = x + 1
          when 'L' then x = x - 1
          when 'U' then y = y + 1
          when 'D' then y = y - 1
        end
        coords.add([x, y])
      end
    end; 0
    paths.push(coords); 0
  end; 0
  return paths
end

def pathsIntersections(pathA,pathB)
  results = []
  pathA.each do |coordsA|
    if pathB.include? coordsA
      res = coordsA[0].abs + coordsA[1].abs
      results.push(res) if res > 0
    end
  end; 0
  return results.min
end

def print_cables(paths)
  maxX, maxY, minX, minY = [0,0,0,0]
  paths.each do |coords|
    coords.each do |coord|
      maxX = coord[0] if coord[0]>maxX
      maxY = coord[1] if coord[1]>maxY
      minX = coord[0] if coord[0]<minX
      minY = coord[1] if coord[1]<minY
    end; 0
  end; 0
  lines = []
  (maxY+minY.abs+1).times do |i|
    line = []
    (maxX+minX.abs+1).times{|j| line.push('.')}
    lines.push(line)
  end
  paths.each_with_index do |coords, i|
    coords.each do |coord|
      if lines[coord[1]+minY.abs][coord[0]+minX.abs] == '.'
        lines[coord[1]+minY.abs][coord[0]+minX.abs] = i
      else
        lines[coord[1]+minY.abs][coord[0]+minX.abs] = 'X'
      end
      newpage = ''
      (maxX+minX.abs+10).times{|j| newpage << "\n"}
      sleep(0.1)
      puts newpage
      puts "#{lines.map{|line| line.join(' ')}.join("\n")}"
    end; 0
  end; 0
end

paths = buildCoords(test_input)#File.read("input.txt"))
puts "RESULT #{pathsIntersections(paths[0],paths[1])}"
print_cables(paths)
