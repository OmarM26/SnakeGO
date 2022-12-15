package main

/*
// REGLAS BÁSICAS:
Grilla de tamaño N×M.
Pueden haber 'p' serpientes.
Las serpientes no pueden transitar por una celda ocupada por otra serpiente.
Cada serpiente se desplaza de forma autónoma intentando llegar primero que las demás al alimento.
Cada vez que una serpiente alcance el alimento, ésta crecerá y, a la vez, deberá surgir una nueva porción de comida en otra celda (definida aleatoriamente).
Cada 'x' segundos cada serpiente se mueve una celda.
El juego debe continuar hasta que las serpientes no puedan seguir moviéndose.
// REGLAS CON CONCURRENCIA:
Las serpientes deben ser representadas por un Thread (liviano o pesado) independiente (con exclusion mutua).
Si una serpiente no puede moverse entra a deadlock permanente (o momentaneo).
El programa debe ser tan asíncrono y libre de barreras de sincronización como sea posible.
// GRILLA:
Defina un mecanismo que permita observar los diferentes estados de las serpientes en la grilla.
*/

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MAIN------------------------------------------------------------------------------------------------------------------------------------------------
var ancho, largo int //Tamaño de la grilla.
var Dir byte         //Direccion de la serpiente.
var oldDir byte      //Direccion de la serpiente.
var posComida Coords //Coordenadas de comida
var list []Snake     // Lista serpientes
var wg sync.WaitGroup

type MyGrilla [][]string //Grilla

func (grilla MyGrilla) String() string {
	out := "█"
	for i := 0; i < len(grilla[0]); i++ {
		out += "▀"
	}
	out += "█\n"
	for i := 0; i < len(grilla); i++ {
		out += "█"
		for j := 0; j < len(grilla[i]); j++ {
			out += grilla[i][j]
		}
		out += "█\n"
	}
	out += "▀"
	for i := 0; i < len(grilla[0]); i++ {
		out += "▀"
	}
	out += "▀\n"
	return out
}

type Coords struct {
	X int
	Y int
}

type Snake struct {
	Cola []Coords
	lost bool
}

func (s *Snake) agregarCola(newLoc Coords) Coords {
	s.Cola = append(s.Cola, newLoc)
	return newLoc
}

func (s *Snake) quitarCola() Coords {
	temp := s.Cola[0]
	s.Cola = s.Cola[1:]
	return temp
}

//RETORNO DE CARRO GOLANG (hacer que la matriz se refresque bonito)

// FUNCIONES------------------------------------------------------------------------------------------------------------------------------------------------

func colocarComida(grilla MyGrilla) {
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	x, y := random.Intn(int(ancho-1)), random.Intn(int(largo-1))
	for grilla[x][y] != " " {
		x, y = random.Intn(int(ancho-1)), random.Intn(int(largo-1))
	}
	grilla[x][y] = "◈"
	posComida.X = x
	posComida.Y = y
}

// 0 arriba;		1 abajo;	2 izquierda;	3 derecha
func verificar() {
	n, o := Dir, oldDir
	if (n == 0 && o == 1) || (n == 1 && o == 0) || (n == 3 && o == 2) || (n == 2 && o == 3) {
		Dir = oldDir

	}
}

func celdaSig(coord Coords, dir byte, grilla MyGrilla) (Coords, bool) {
	tempCoord := coord
	switch dir {
	case 0:
		tempCoord.X -= 1
		break
	case 2:
		tempCoord.Y -= 1
		break
	case 1:
		tempCoord.X += 1
		break
	case 3:
		tempCoord.Y += 1
		break
	}

	// Golpear la pared
	if tempCoord.X >= ancho || tempCoord.Y >= largo || tempCoord.X < 0 || tempCoord.Y < 0 {
		return tempCoord, false
	}
	// Comerse a si mismo
	if grilla[tempCoord.X][tempCoord.Y] != " " && grilla[tempCoord.X][tempCoord.Y] != "◈" {
		return tempCoord, false
	}
	return tempCoord, true
}

func actualizarGrilla(grilla MyGrilla, snake *Snake) {
	last := int(len(snake.Cola) - 1)
	newCell, cont := celdaSig(snake.Cola[last], Dir, grilla)

	if !cont {
		snake.lost = true
		return
	}

	snake.agregarCola(newCell)

	last = len(snake.Cola) - 1
	grilla[snake.Cola[last].X][snake.Cola[last].Y] = "□"

	if newCell == posComida {
		colocarComida(grilla)
	} else {
		grilla[snake.Cola[0].X][snake.Cola[0].Y] = " "
		snake.quitarCola()
	}
	oldDir = Dir
}

// MAIN------------------------------------------------------------------------------------------------------------------------------------------------
func main() {
	//Argumentos
	ancho1 := flag.Int("ancho", 0, "El ancho de la grilla.")
	largo1 := flag.Int("largo", 18, "El largo de la grilla.")
	vel1 := flag.Int("velocidad", 20, "La velocidad de la grilla.")
	flag.Parse()
	fmt.Println("Ancho = ", *ancho1)
	fmt.Println("Largo = ", *largo1)
	fmt.Println("Velocidad = ", *vel1)
	cant_snake := 3
	wg.Add(3)
	//Dimensiones de la grilla
	ancho, largo = 10, 25

	//Grilla Vacia
	var grilla MyGrilla
	grilla = make([][]string, ancho)
	for i := 0; i < int(ancho); i++ {
		grilla[i] = make([]string, largo)
	}

	//Inicializar Grilla
	for i := 0; i < int(ancho); i++ {
		for j := 0; j < int(largo); j++ {
			grilla[i][j] = string(" ")
		}
	}
	for i := 0; i <= int(cant_snake); i++ {
		// LA POSICIÓN DE LA SERPIENTE NO DEBE CREAR UNA NUEVA GRILLA
		var s Snake

		s.lost = false
		time.Sleep(time.Millisecond * time.Duration(100))
		rand.Seed(time.Now().UnixNano())
		serpX, serpY := rand.Intn(int(ancho-1)), rand.Intn(int(largo-1))
		s.agregarCola(Coords{serpX, serpY + 1})
		grilla[serpX][serpY+1] = "□"
		// SE AGREGA LA SERPIENTE A LA VARIABLE LIST
		list = append(list, s)
	}
	//Inicializar serpientes.
	for i := 0; i < int(cant_snake); i++ {
		ch1 := make(chan byte)

		go func(ch1 chan byte) {
			//Movimiento serpiente
			for {

				last := int(len((list[i]).Cola) - 1)
				time.Sleep(time.Millisecond * time.Duration(100))
				if posComida.X > list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					if Dir == 0 && list[i].Cola[last].Y-1 > 0 && (grilla[list[i].Cola[last].X][list[i].Cola[last].Y-1] == " " || grilla[list[i].Cola[last].X][list[i].Cola[last].Y-1] == "◈") {
						oldDir = 0
						ch1 <- 2
					} else if Dir == 0 && list[i].Cola[last].Y-1 == 0 && (grilla[list[i].Cola[last].X][list[i].Cola[last].Y+1] == " " || grilla[list[i].Cola[last].X][list[i].Cola[last].Y+1] == "◈") {
						oldDir = 0
						ch1 <- 3
					} else if Dir == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
						oldDir = 0
						ch1 <- 0
					} else {
						ch1 <- 1
					}
				} else if posComida.X < list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					ch1 <- 0
				} else if posComida.X == list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					if posComida.Y > list[i].Cola[last].Y {
						//Estos aun estan con pruebas
						//Condiciones para "evitar" Deadlock
						//if Dir == 2 && snakes[i].Cola[last].X-1 > 0 && grilla[snakes[i].Cola[last].X-1][snakes[i].Cola[last].Y] == "" {
						if Dir == 2 && list[i].Cola[last].X-1 > 0 && (grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == "◈") {
							fmt.Println("<- to up")
							oldDir = 2
							ch1 <- 0
						} else if Dir == 2 && list[i].Cola[last].X-1 == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
							fmt.Println("<- to down")
							oldDir = 2
							ch1 <- 1
						} else {
							ch1 <- 3
						}

					} else if posComida.Y < list[i].Cola[last].Y {
						if Dir == 3 && list[i].Cola[last].X-1 > 0 && (grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == "◈") {
							fmt.Println("-> to up")
							oldDir = 3
							ch1 <- 0
						} else if Dir == 3 && list[i].Cola[last].X-1 == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
							fmt.Println("-> to down")
							oldDir = 3
							ch1 <- 1
						} else {
							ch1 <- 2
						}
					}
				}
				//Deadlock
				if list[i].lost {
					for true {
					}
				}
			}
		}(ch1)

		//Colocar comida en la grilla
		colocarComida(grilla)
		fmt.Println(grilla)

		//Ejecución del juego
		for {
			//Direccion serpiente
			select {
			case stdin, _ := <-ch1:
				Dir = stdin
				verificar()
				actualizarGrilla(grilla, &list[i])
				time.Sleep(time.Millisecond * time.Duration(100))

			}
			time.Sleep(time.Millisecond * time.Duration(200))
			fmt.Println(grilla) //Imprimir pasos de la grilla
			if list[i].lost == true {
				fmt.Println("Todas las serpientes han muerto.")
				break
			}
		}
	}
}
