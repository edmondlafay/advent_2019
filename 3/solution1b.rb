#!/usr/bin/env ruby
require 'ruby2d'
require 'set'
test_input = "R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51\nU98,R91,D20,R16,D67,R40,U7,R15,U6,R7"

def findDuplicate(array)
  map = {}
  array.each do |value|
    map[value] = (map[value] || 0 ) + 1
    return value if map[value] > 1
  end
  nil
end

def rludToSegments(input)
  paths = []
  input.split("\n").each do |path_raw|
    segments = []
    x1, y1, x2, y2 = [0, 0, 0, 0]
    path_raw.split(',').each do |dir_dist|
      order = dir_dist[0]
      distance = dir_dist[1..-1].to_i
      case order
        when 'R' then x2 = x1 + distance
        when 'L' then x2 = x1 - distance
        when 'U' then y2 = y1 + distance
        when 'D' then y2 = y1 - distance
      end
      segments.push({x1:x1, y1:y1, x2:x2, y2:y2})
      x1, y1 = [x2, y2]
    end
    paths.push(segments)
  end
  return paths
end

def segmentsIntersectSegment?(segmentA, segmentB)
  if ([segmentA[:x1], segmentA[:x2]].max > [segmentB[:x1], segmentB[:x2]].max and 
      [segmentA[:x1], segmentA[:x2]].min < [segmentB[:x1], segmentB[:x2]].min and 
      [segmentA[:y1], segmentA[:y2]].max < [segmentB[:y1], segmentB[:y2]].max and 
      [segmentA[:y1], segmentA[:y2]].min > [segmentB[:y1], segmentB[:y2]].min) or (
      [segmentB[:x1], segmentB[:x2]].max > [segmentA[:x1], segmentA[:x2]].max and 
      [segmentB[:x1], segmentB[:x2]].min < [segmentA[:x1], segmentA[:x2]].min and 
      [segmentB[:y1], segmentB[:y2]].max < [segmentA[:y1], segmentA[:y2]].max and 
      [segmentB[:y1], segmentB[:y2]].min > [segmentA[:y1], segmentA[:y2]].min)
    return true
  end
  false
end

def segmentsIntersectPath?(segmentA, pathB)
  intersections = []
  pathB.each do |segmentB|
    if segmentsIntersectSegment?(segmentA, segmentB)
      x = findDuplicate([segmentA[:x1], segmentA[:x2], segmentB[:x1], segmentB[:x2]])
      y = findDuplicate([segmentA[:y1], segmentA[:y2], segmentB[:y1], segmentB[:y2]])
      intersections.push({x: x, y: y})
    end
  end
  intersections
end

def pathIntersectPath?(pathA, pathB)
  intersections = []
  pathA.each do |segmentA|
    intersections += segmentsIntersectPath?(segmentA, pathB)
  end
  intersections
end

def print_cables_segments(paths)
  maxX, maxY, minX, minY = [0,0,0,0]
  paths.each do |coords|
    coords.each do |coord|
      maxX, maxY, minX, minY = [ 
        [maxX, coord[:x1], coord[:x2]].max,  [maxY, coord[:y1], coord[:y2]].max,
        [minX, coord[:x1], coord[:x2]].min, [minY, coord[:y1], coord[:y2]].min
      ]
    end
  end
  minX, minY = [minX.abs, minY.abs]
  width, height = [maxX+minX, maxY+minY]

  coef = 1
  if width>1080 or height>750
    coef = [1080.0/width, 750.0/height].min
  elsif width<1080 or height<750
    coef = [1080/width, 750/height].max
  end

  puts "frame size : #{(width*coef).to_i} #{(height*coef).to_i} (coef #{coef})"
  set title: "Central Port", width: ((width+1)*coef).to_i+1, height: ((height+1)*coef).to_i+1
  advence = 0
  colors = ['blue', 'red', 'green']

  intersectShapes = {
    999999999 => Circle.new(x: ((paths[0][0][:x1]+minX)*coef).to_i, y: ((paths[0][0][:y1]+minY)*coef).to_i, z: 100,radius: coef + 5, color: 'yellow')
  }
  update do
    if get(:frames)%6==0
      paths.each_with_index do |path, i|
        if advence < path.length
          x1, y1, x2, y2  = [path[advence][:x1], path[advence][:y1], path[advence][:x2], path[advence][:y2]]
          Line.new(x1: ((x1+minX)*coef).to_i, y1: ((y1+minY)*coef).to_i, x2: ((x2+minX)*coef).to_i, y2: ((y2+minY)*coef).to_i, width: [1, coef].max, color: colors[i], z: i)
          intersections = segmentsIntersectPath?(path[advence], paths[(i+1)%paths.length][0..[path.length,advence].min])
          intersections.each do |intersection|
            intersectShapes[intersection[:x].abs + intersection[:y].abs] = Sprite.new(
              'boom.png', clip_width: 127, time: 75, loop: true, width: 2 * (coef + 5), height: 2 * (coef + 5), z: 200,
              x: ((intersection[:x]+minX)*coef).to_i - (coef + 5), y: ((intersection[:y]+minY)*coef).to_i - (coef + 5)
            )
            intersectShapes[intersection[:x].abs + intersection[:y].abs].play
          end
        end
      end
      advence += 1
      if advence == [paths[0].length, paths[1].length].max
        closest = intersectShapes[intersectShapes.keys.min]
        
        coin = Sprite.new(
          'coin.png', clip_width: 84, time: 150, loop: true, width: 2 * (coef + 5), height: 2 * (coef + 5), z: 200, x: closest.x, y: closest.y
        )
        coin.play
      end
    end
  end
  show
end

input = test_input

start = Time.now
input = File.read("input.txt")
paths = rludToSegments(input)
intersections = pathIntersectPath?(paths.first, paths.last)
res = intersections.map{|int|int[:x].abs+int[:y].abs}.min
finish = Time.now
puts "RESULT1 #{res} (in #{finish - start} seconds)"

print_cables_segments(paths)
