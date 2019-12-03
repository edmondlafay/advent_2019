# calculate its fuel and add it to the total. Then, treat the fuel amount you
# just calculated as the input mass and repeat the process, continuing until
# a fuel requirement is zero or negative

def fuel_for_mass(mass)
  res = (mass/3)-2
  if res > 8
    res = res + fuel_for_mass(res)
  else
    res
  end
end

res = 0
input.each{|i| res=res+fuel_for_mass(i)}
puts res
