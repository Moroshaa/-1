## Лабораторая работа №7

# Для проверки первого задания запустить serever.go следом client.go

# Для проверки http
запустить http.go зайти в браузер вставить http://localhost:8080/hello
для проверки Post/data
-X POST http://localhost:8080/data -d '{"key": "value"}' -H "Content-Type: application/json"


# Для проверки websocket
запустить websocket.go открыть консоль в браузере и вставить следующее 


const socket = new WebSocket("ws://localhost:8080/ws");

socket.onopen = () => {
    console.log("Соединение установлено!");
    socket.send("Привет от клиента!");
};

socket.onmessage = (event) => {
    console.log("Новое сообщение:", event.data);
};

socket.onclose = () => {
    console.log("Соединение закрыто.");
};

// Отправить сообщение
function sendMessage(message) {
    socket.send(message);
}

если не ворк, скрипнуть предупреждение "allow pasting" и попробовать еще раз 
