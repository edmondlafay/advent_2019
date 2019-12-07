#!/usr/bin/env ruby
input = File.read("input.txt")
#input = "x=495, y=2..7\ny=7, x=495..501\nx=501, y=3..7\nx=498, y=2..4\nx=506, y=1..2\nx=498, y=10..13\nx=504, y=10..13\ny=13, x=498..504"

def input_parse(input)
  list_coords = []
  height = 0
  xMin, xMax = [500, 0]
  input.split("\n").each do |com_coords|
    if com_coords[0]=='x'
      x_raw, y_raw = com_coords.split(",")
      x = x_raw.delete("^0-9").to_i
      xMin, xMax = [[xMin, x].min, [xMax, x].max]
      y_min, y_max = y_raw.delete(" y=").split('..').map{|i| i.to_i}
      height = [height, y_max].max
      (y_min..y_max).each do |y| 
        list_coords.push({x: x, y: y})
      end
    else
      y_raw, x_raw = com_coords.split(",")
      y = y_raw.delete("^0-9").to_i
      height = [height, y].max
      x_min, x_max = x_raw.delete(" x=").split('..').map{|i| i.to_i}
      xMin, xMax = [[xMin, x_min].min, [xMax, x_max].max]
      (x_min..x_max).each do |x| 
        list_coords.push({x: x, y: y})
      end
    end
  end

  drawing = []
  height.times{|h| drawing.push('.'*((xMax - xMin)+3))}
  list_coords.each{|coords| drawing[coords[:y]-1][coords[:x] - xMin+1] = '#'} # draw clay
  drawing[0][500-xMin+1] = '|'
  drawing
end

def update_drawing(drawing)
  step = 0
  heighest_update = 0
  count_water = 1
  loop do
    step+=1
    puts "#{"\n"*100}"
    puts "step: #{step}, heighest_update: #{heighest_update}/#{drawing.length}, , count_water: #{count_water}"
    puts "#{drawing[[0,heighest_update-35].max..[heighest_update+35,drawing.length-1].min].join("\n")}"
    was_updated = false
    (heighest_update..drawing.length-1).each do |y|
      drawing[y].split('').each_with_index do |character, x|
        if drawing[y][x]=='|'
          if y==drawing.length-1 or drawing[y+1][x] == '|'
            # this water has already been treated, ignore
          elsif drawing[y+1][x] == '.'
            drawing[y+1][x] = '|'
            count_water+=1
            unless was_updated
              heighest_update = y+1
              was_updated = true
            end
          else
            if ['#', '~'].include? drawing[y+1][x]
              # fill left
              l = x
              until drawing[y][l-1]=='#' or drawing[y+1][l] == '.' or drawing[y+1][l] == '|'
                unless drawing[y][l-1] == '|'
                  drawing[y][l-1] = '|'
                  count_water+=1
                  if was_updated
                    heighest_update = [heighest_update, y].min
                  else
                    heighest_update = y-1
                    was_updated = true
                  end
                end
                l-=1
              end
              # fill right
              r = x
              until drawing[y][r+1]=='#' or drawing[y+1][r]=='.' or drawing[y+1][r]=='|'
                unless drawing[y][r+1] == '|'
                  drawing[y][r+1] = '|'
                  count_water+=1
                  if was_updated
                    heighest_update = [heighest_update, y].min
                  else
                    heighest_update = y-1
                    was_updated = true
                  end
                end
                r+=1
              end
              # check if filled with water
              if (drawing[y][l-1] == '#') and (drawing[y][r+1] == '#')
                (l..r).each do |i|
                  unless drawing[y][i] == '~'
                    drawing[y][i]='~'
                    if was_updated
                      heighest_update = [heighest_update, y].min
                    else
                      heighest_update = y-1
                      was_updated = true
                    end
                  end
                end
              end
            end
          end
        end
      end
    end
    unless was_updated
      puts "#{"\n"*100}#{drawing.join("\n")}"
      return count_water
    end
  end
end

puts "RESULT #{update_drawing(input_parse(input))}"
