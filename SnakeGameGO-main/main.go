package main

/*
// REGLAS BÁSICAS:
Grilla de tamaño N×M.
Pueden haber 'p' serpientes.
Cada vez que una serpiente alcance el alimento, ésta crecerá y, a la vez, deberá surgir una nueva porción de comida en otra celda (definida aleatoriamente).
Cada 'x' segundos cada serpiente se mueve una celda.

Cada serpiente se desplaza de forma autónoma intentando llegar primero que las demás al alimento.
Las serpientes no pueden transitar por una celda ocupada por otra serpiente.
El juego debe continuar hasta que las serpientes no puedan seguir moviéndose.

// REGLAS CON CONCURRENCIA:
Las serpientes deben ser representadas por un Thread (liviano o pesado) independiente (con exclusion mutua).
Si una serpiente no puede moverse entra a deadlock permanente (o momentaneo).
El programa debe ser tan asíncrono y libre de barreras de sincronización como sea posible.

// GRILLA:
Defina un mecanismo que permita observar los diferentes estados de las serpientes en la grilla.

// OPCIONAL:
Hacer que la matriz se refresque bien (Retorno de Carro)

Movimiento serpientes:
	0 = arriba
	1 = abajo
	2 = izquierda
	3 = derecha
*/

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DEFINICIONES ------------------------------------------------------------------------------------------------------------------------------------------------
var ancho, largo int // Tamaño de la grilla.
var vel int          // Tamaño de la grilla.
var cantSerp int     // Tamaño de la grilla.
var Dir byte         // Direccion de la serpiente.
var oldDir byte      // Direccion de la serpiente.
var posComida Coords // Coordenadas de comida.
var list []Snake     // Lista serpientes.
var w sync.WaitGroup // Esperar goroutines.

type MyGrilla [][]string // Grilla

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

// FUNCIONES ------------------------------------------------------------------------------------------------------------------------------------------------
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

// MAIN ------------------------------------------------------------------------------------------------------------------------------------------------
func main() {
	// Argumentos
	ancho1 := flag.Int("ancho", 10, "El ancho de la grilla.")
	largo1 := flag.Int("largo", 25, "El largo de la grilla.")
	vel1 := flag.Int("velocidad", 100, "La velocidad de la grilla.")
	cantSerp1 := flag.Int("serpientes", 1, "Cantidad de Serpientes.")
	flag.Parse()

	// Dimensiones de la grilla
	ancho, largo = *ancho1, *largo1
	// Velocidad
	vel = *vel1
	// Cantidad de Serpientes
	cantSerp = *cantSerp1

	// Grilla Vacia
	var grilla MyGrilla
	grilla = make([][]string, ancho)
	for i := 0; i < int(ancho); i++ {
		grilla[i] = make([]string, largo)
	}

	// Inicializar Grilla
	for i := 0; i < ancho; i++ {
		for j := 0; j < largo; j++ {
			grilla[i][j] = string(" ")
		}
	}

	// Inicializar Serpientes
	for i := 0; i < cantSerp; i++ {
		var s Snake
		s.lost = false
		time.Sleep(time.Millisecond * time.Duration(100))
		rand.Seed(time.Now().UnixNano())
		serpX, serpY := rand.Intn(int(ancho-1)), rand.Intn(int(largo-1))
		s.agregarCola(Coords{serpX, serpY + 1})
		grilla[serpX][serpY+1] = "□"
		list = append(list, s)
	}

	// Colocar comida en la grilla
	colocarComida(grilla)
	fmt.Println(grilla)

	// Goroutines serpientes
	ch := make(chan byte)
	for i := 0; i < cantSerp; i++ {
		w.Add(1)
		go func(i int, co chan byte) {
			// Movimiento serpiente
			for {
				// Deadlock
				if list[i].lost {
					for true {
					}
				}
				// Movimientos
				last := int(len(list[i].Cola) - 1)
				time.Sleep(time.Millisecond * time.Duration(100))
				if posComida.X > list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					if Dir == 0 && list[i].Cola[last].Y-1 > 0 && (grilla[list[i].Cola[last].X][list[i].Cola[last].Y-1] == " " || grilla[list[i].Cola[last].X][list[i].Cola[last].Y-1] == "◈") {
						oldDir = 0
						co <- 2
					} else if Dir == 0 && list[i].Cola[last].Y-1 == 0 && (grilla[list[i].Cola[last].X][list[i].Cola[last].Y+1] == " " || grilla[list[i].Cola[last].X][list[i].Cola[last].Y+1] == "◈") {
						oldDir = 0
						co <- 3
					} else if Dir == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
						oldDir = 0
						co <- 0
					} else {
						co <- 1
					}
				} else if posComida.X < list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					co <- 0
				} else if posComida.X == list[i].Cola[last].X {
					time.Sleep(time.Millisecond * time.Duration(100))
					if posComida.Y > list[i].Cola[last].Y {
						if Dir == 2 && list[i].Cola[last].X-1 > 0 && (grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == "◈") {
							oldDir = 2
							co <- 0
						} else if Dir == 2 && list[i].Cola[last].X-1 == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
							oldDir = 2
							co <- 1
						} else {
							co <- 3
						}
					} else if posComida.Y < list[i].Cola[last].Y {
						if Dir == 3 && list[i].Cola[last].X-1 > 0 && (grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X-1][list[i].Cola[last].Y] == "◈") {
							oldDir = 3
							co <- 0
						} else if Dir == 3 && list[i].Cola[last].X-1 == 0 && (grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == " " || grilla[list[i].Cola[last].X+1][list[i].Cola[last].Y] == "◈") {
							oldDir = 3
							co <- 1
						} else {
							co <- 2
						}
					}
				}
			}
		}(i, ch)

	}
	// Ejecucion del Juego
	for {
		// Direccion Serpiente
		select {
		case stdin, _ := <-ch:
			Dir = stdin
			verificar()
			for i := 0; i < len(list); i++ {
				actualizarGrilla(grilla, &list[i])
				time.Sleep(time.Millisecond * time.Duration(vel))
			}
		}
		// Imprimir pasos de la Grilla
		time.Sleep(time.Millisecond * time.Duration(200))
		fmt.Println(grilla)
		// Verificar fin del Juego
		contar_Muertas := 0
		for i := 0; i < len(list); i++ {
			if list[i].lost == true {
				contar_Muertas++
			}
		}
		if contar_Muertas >= cantSerp {
			fmt.Println("Todas las serpientes han muerto.")
			break
		}
	}
	close(ch)
}
