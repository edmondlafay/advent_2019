require 'set'
test_input = "R75,D30,R83,U83,L12,D49,R71,U7,L72\nU62,R66,U55,R34,D71,R55,D58,R83"
test output = 159

def buildCoords(input, dict)
  paths = []
  list = -1
  input.split("\n").each do |path_raw|
    list = list + 1
    coords = Set[]
    distanceTot = 0
    x1 = 0
    y1 = 0
    path_raw.split(',').each do |dir_dist|
      order = dir_dist[0]
      distance = dir_dist[1..-1].to_i
      y2 = y1
      x2 = x1
      case order
      when 'R'
        distance.times do |i| 
          x2 = x2 + 1
          coords.add([x2, y2])
          distanceTot += 1
          dict[list][x2+10000*y2] = distanceTot unless dict[list].key? x2+10000*y2
        end
      when 'L'
        distance.times do |i| 
          x2 = x2 - 1
          coords.add([x2, y2])
          distanceTot += 1
          dict[list][x2+10000*y2] = distanceTot unless dict[list].key? x2+10000*y2
        end
      when 'U'
        distance.times do |i| 
          y2 = y2 + 1
          coords.add([x2, y2])
          distanceTot += 1
          dict[list][x2+10000*y2] = distanceTot unless dict[list].key? x2+10000*y2
        end
      when 'D'
        distance.times do |i| 
          y2 = y2 - 1
          coords.add([x2, y2])
          distanceTot += 1
          dict[list][x2+10000*y2] = distanceTot unless dict[list].key? x2+10000*y2
        end
      end
      x1 = x2
      y1 = y2
    end
    paths.push(coords)
  end
  return paths
end

def pathsIntersections(pathA, pathB, dict)
  pathA.each do |coordsA|
    if pathB.include? coordsA
      return dict[0][coordsA[0]+10000*coordsA[1]] + dict[1][coordsA[0]+10000*coordsA[1]]
    end
  end
  return -1
end

dist_hashA = {}
dist_hashB = {}
dict = [dist_hashA, dist_hashB]

paths = buildCoords(File.read("input.txt"), dict)
puts 'RESULT', pathsIntersections(paths[0], paths[1], dict)
