#!/usr/bin/env ruby
require 'set'

class Planet
  attr_accessor :name, :is_orbiting, :orbeted_by
  def initialize(name, is_orbiting, orbeted_by = [])
    @name = name
    @is_orbiting = is_orbiting
    @orbeted_by = orbeted_by
  end
end

def build_system()
  known_planets = {}
  input = File.read("input.txt")
  #input = "COM)B\nB)C\nC)D\nD)E\nE)F\nB)G\nG)H\nD)I\nE)J\nJ)K\nK)L\nK)YOU\nI)SAN"
  input.split("\n").each do |relation|
    center, orbiting = relation.split(')')
    known_planets[orbiting] = Planet.new(orbiting, nil) unless known_planets.has_key? orbiting
    if known_planets.has_key? center
      known_planets[center].orbeted_by.push known_planets[orbiting]
    else
      known_planets[center] = Planet.new(center, nil, [known_planets[orbiting]])
    end
    known_planets[orbiting].is_orbiting = known_planets[center]
  end
  return known_planets
end

def find_center(known_planets)
  name, center = known_planets.first
  until center.is_orbiting.nil?
    center = center.is_orbiting
  end
  return center
end

def find_leafs(known_planets, center, leafs=Set[])
  leafs.add(center) if center.orbeted_by.empty?
  center.orbeted_by.each do |orbiting|
    leafs = leafs | find_leafs(known_planets, orbiting, leafs)
  end
  return leafs
end

def count_orbits(known_planets, leaf, i=0)
  return i if leaf.is_orbiting.nil?
  return i + count_orbits(known_planets, leaf.is_orbiting, i) + 1
end

def reverse_path(known_planets, leaf)
  return {leaf.name => leaf.name} if leaf.is_orbiting.nil?
  res = reverse_path(known_planets, leaf.is_orbiting)
  res[leaf.name] = "#{res[leaf.is_orbiting.name]})#{leaf.name}"
  return res
end

def stringsCommonStart(a, b)
  res = ''
  a.split('').each_with_index do |char, i|
    return res unless a[i]==b[i]
    res+=char
  end
end

known_planets = build_system()
center = find_center(known_planets)
leafs = find_leafs(known_planets, center)
puts "Number of orbits: #{known_planets.map{|name, planet| count_orbits(known_planets, planet)}.sum}"
youPath = reverse_path(known_planets, known_planets['YOU'])['YOU']
santaPath = reverse_path(known_planets, known_planets['SAN'])['SAN']
comonPath = stringsCommonStart(youPath, santaPath)
youPath.slice! comonPath
santaPath.slice! comonPath
puts "Fastest routes: #{youPath.count(')')+santaPath.count(')')}"

