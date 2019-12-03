# take its mass, divide by three, round down, and subtract 2

def fuel_for_mass(mass)
  return (mass/3)-2
end

res = 0
input.each{|i| res=res+fuel_for_mass}
puts res