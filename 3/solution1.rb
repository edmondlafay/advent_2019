require 'ruby2d'
require 'set'
test_input = "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7"
test output = 159
start = Time.now

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
  width = maxX+minX.abs+1
  height = maxY+minY.abs+1

  arrayOfCoordsA = []
  arrayOfCoordsB = []
  paths.each_with_index do |coords, i|
    coords.each do |coord|
      arrayOfCoordsA.push [i, coord[0], coord[1]] if i == 0
      arrayOfCoordsB.push [i, coord[0], coord[1]] if i == 1
    end; 0
  end; 0

  coef = 1
  if width>1080 or height>800
    coef = [1080.0/width, 800.0/height].min
    puts [width/1080.0, height/800.0]
  elsif width<1080 or height<800
    coef = [1080/width, 800/height].max
  end

  puts "frame size : #{(width*coef).to_i} #{(height*coef).to_i} (coef #{coef})"
  set title: "Central Port", width: ((width+3)*coef).to_i, height: ((height+3)*coef).to_i
  tick = 0
  advence = 0
  colors = ['blue', 'red']
  update do
    if (coef<1 or tick%3==0)
      if advence < arrayOfCoordsA.length
        path, x, y = arrayOfCoordsA[advence]
        Square.new(x: ((x+minX.abs+1)*coef).to_i, y: ((y+minY.abs+1)*coef).to_i, z: path,
          size: [1, coef].max, color: colors[path % colors.length])
      end
      if advence < arrayOfCoordsB.length
        path, x, y = arrayOfCoordsB[advence]
        Square.new(x: ((x+minX.abs+1)*coef).to_i, y: ((y+minY.abs+1)*coef).to_i, z: path,
          size: [1, coef].max, color: colors[path % colors.length])
      end
      advence += 1
    end
    tick += 1
  end
  show
end

#paths = buildCoords(test_input); 0
paths = buildCoords(File.read("input.txt")); 0
res = pathsIntersections(paths[0],paths[1])
finish = Time.now
puts "RESULT1 #{res} (in #{finish - start} seconds)"
print_cables(paths)
