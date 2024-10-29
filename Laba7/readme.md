# Лабораторая работа №7

## Для проверки tcp 
запустить serever.go следом client.go и отправит сообщение

## Для проверки http
запустить http.go зайти в браузер вставить _http://localhost:8080/hello_
для проверки Post/data
_-X POST http://localhost:8080/data -d '{"key": "value"}' -H "Content-Type: application/json"_


## Для проверки websocket
запустить websocket.go открыть консоль в браузере и вставить следующее 


`const socket = new WebSocket("ws://localhost:8080/ws");`

`socket.onopen = () => {
    console.log("Соединение установлено!");
    socket.send("Привет от клиента!");
};`

`socket.onmessage = (event) => {
    console.log("Новое сообщение:", event.data);
};`

`socket.onclose = () => {
    console.log("Соединение закрыто.");
};`

// Отправить сообщение
`function sendMessage(message) {
    socket.send(message);
}`

если не ворк, скрипнуть предупреждение "allow pasting" и попробовать еще раз 
