#!/usr/bin/env ruby
require 'set'
input = File.read("input.txt")

def input_parse(input)
  list_coords = []
  input.split("\n").each do |com_coords|

    if com_coords[0]=='x'
      x_raw, y_raw = com_coords.split(",")
      x = x_raw.delete("^0-9").to_i
      y_raw = y_raw.delete(" y=").split('..').map{|i| i.to_i}
      (y_raw.first..y_raw.last).each{|y| list_coords.push({x: x, y: y})}
    else
      y_raw, x_raw = com_coords.split(",")
      y = y_raw.delete("^0-9").to_i
      x_raw = x_raw.delete(" x=").split('..').map{|i| i.to_i}
      (x_raw.first..x_raw.last).each{|x| list_coords.push({x: x, y: y})}
    end
  end
  list_coords
end

def get_min_max(list_coords)
  xMin, xMax, yMax = [500,0,0]
  list_coords.each do |coords|
    xMin = [xMin, coords[:x]].min
    xMax = [xMax, coords[:x]].max
    yMax = [yMax, coords[:y]].max
  end
  [xMin, xMax, yMax]
end

def update_drawing(drawing, wl={x:500, y:0})
  if wl[:y]+1 == drawing.length
    # flow out of drawing
  elsif drawing[wl[:y]+1][wl[:x]] == '.' # go down
    drawing[wl[:y]+1][wl[:x]] = '|'
  elsif drawing[wl[:y]+1][wl[:x]] == '#' # go sideways
    
    l = wl[:x]
    r = wl[:x]
    # fill left
    until l-1 < 0 or drawing[wl[:y]][l-1] == "#" or drawing[wl[:y]+1][l] = '.'
      l=l-1
      drawing[wl[:y]][l] = '|'
    end
    wl = {x:l, y:wl[:y]}
    # fill right
    until r+1 >= drawing.length or drawing[wl[:y]][r+1] == "#" or drawing[wl[:y]+1][r] = '.'
      r=r+1
      drawing[wl[:y]][r] = '|'
    end
    wl = {x:r, y:wl[:y]}
  end
end

list_coords = input_parse(input)
puts list_coords
xMin, xMax, yMax = get_min_max(list_coords)
width = (xMax) - (xMin)
height = yMax
puts width
drawing = []
yMax.times{|h| drawing.push('.'*(width+3))}
list_coords.each{|coords| drawing[coords[:y]-1][coords[:x] - xMin+1] = '#'} # draw clay
drawing[0][500-xMin+1] = '|' # draw source
puts "#{drawing.join("\n")}"


update_drawing(drawing)
puts "#{drawing.join("\n")}"
