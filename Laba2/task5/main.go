package main

import "fmt"

type Rectangle struct {
	width  int
	heigth int
}

func (r Rectangle) Area() int {
	return r.width * r.heigth

}
func main() {
	fmt.Println("Задание 5")
	rect := Rectangle{width: 10, heigth: 20}
	area := rect.Area()
	fmt.Println("Площадь прямугольника с высотой ", rect.heigth, " и шириной ", rect.width, " равна ", area)
}
