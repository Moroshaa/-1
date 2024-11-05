# Лаба 8
## Настройка
редактируем следующую строку под нашу бд 
dsn := "postgres://username:password@localhost:8080/userdb?sslmode=disable"
## Проверка get all 
http://localhost:9090/users
## Проверка get user by id
http://localhost:9090/users/id
## Проверка post 
http://localhost:9090/users
Вводим информацию в формате json 
{"name": "Ромбастик", "age": 20}
## Проверка put 
http://localhost:9090/users/2 
и редактируем 
## Проверка delete 
http://localhost:9090/users/2
