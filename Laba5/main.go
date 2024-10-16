package main

import (
	"fmt"
	"math"
	"os"
)

// Структура Person
type Person struct {
	name string
	age  int
}

// Метод для вывода информации о человеке
func (p Person) Info() {
	fmt.Printf("Name: %s, Age: %d\n", p.name, p.age)
}

// Метод birthday увеличивает возраст на 1 год
func (p *Person) Birthday() {
	p.age++
}

// Структура Circle
type Circle struct {
	radius float64
}

// Метод для вычисления площади круга
func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

// Структура Rectangle
type Rectangle struct {
	width, height float64
}

// Метод для вычисления площади прямоугольника
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Интерфейс Shape
type Shape interface {
	Area() float64
}

// Функция, которая принимает срез интерфейсов Shape и выводит площадь каждого объекта
func printAreas(shapes []Shape) {
	for _, shape := range shapes {
		fmt.Printf("Area: %.2f\n", shape.Area())
	}
}

// Интерфейс Stringer
type Stringer interface {
	String() string
}

// Структура Book
type Book struct {
	title  string
	author string
}

// Реализация интерфейса Stringer для структуры Book
func (b Book) String() string {
	return fmt.Sprintf("Book: %s by %s", b.title, b.author)
}

func main() {
	for {
		var action int

		// Пример объектов для использования
		p := Person{name: "John", age: 30}
		c := Circle{radius: 5}
		r := Rectangle{width: 10, height: 5}
		book := Book{title: "Зеленая миля", author: "Стивенк Кинг"}

		fmt.Println("Выберите действие:")
		fmt.Println("1: Вывести информацию о человеке")
		fmt.Println("2: Отпраздновать день рождения (увеличить возраст)")
		fmt.Println("3: Вычислить площадь круга")
		fmt.Println("4: Вычислить площадь прямоугольника")
		fmt.Println("5: Вывести информацию о книге")
		fmt.Println("0: Выход")

		// Считываем выбор пользователя
		fmt.Scan(&action)

		switch action {
		case 1:
			// Вывод информации о человеке
			p.Info()
		case 2:
			// Увеличить возраст и вывести обновленную информацию
			p.Birthday()
			p.Info()
		case 3:
			// Вычислить площадь круга
			fmt.Printf("Circle Area: %.2f\n", c.Area())
		case 4:
			// Вычислить площадь прямоугольника
			fmt.Printf("Rectangle Area: %.2f\n", r.Area())
		case 5:
			// Вывод информации о книге
			fmt.Println(book)
		case 0:
			os.Exit(1)
		default:
			fmt.Println("Неверный выбор действия")
		}
	}
}
