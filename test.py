import os

print('no param')
os.system('curl "localhost:8080/books" --request "GET"')
print('\n')

print('param id')
os.system('curl "localhost:8080/books?id=2" --request "GET"')
print('\n')

print('param slug')
os.system('curl "localhost:8080/books?slug=mein-kampft" --request "GET"')
print('\n')

print('double param')
os.system('curl "localhost:8080/books?id=2&slug=main-kampft" --request "GET"')
print('\n')
