package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type user struct {
	id   int
	conn net.Conn
	posX int
	posY int
}

var clients []user
var gridSize int
var grid [][]int

func main() {
	fmt.Println("Server Initialised")
	setupGrid()

	clientID := 1
	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		client := user{clientID, conn, 4, 4}
		grid[4][4] = 2
		clientID++
		joinClient := ""
		joinClient = joinClient + strconv.Itoa(gridSize) + " " + strconv.Itoa(len(clients)) + " "
		fmt.Println(len(clients))
		if len(clients) > 0 {
			for i := 0; i < len(clients); i++ {
				joinClient = joinClient + strconv.Itoa(clients[i].posX) + " " + strconv.Itoa(clients[i].posY) + " "
			}
			fmt.Println(joinClient)
		}
		io.WriteString(client.conn, joinClient+"\n")
		clients = append(clients, client)
		join := "JOIN"
		sendUpdate(client.id, join)

		fmt.Println("New client joined game, number of clients = " + strconv.Itoa(len(clients)))
		for i := 0; i < len(clients); i++ {
			go handleInput(&clients[i])
		}
	}
}

func setupGrid() {
	//Initialise Grid
	gridSize = 10
	grid = make([][]int, gridSize) // rows
	for i := range grid {
		grid[i] = make([]int, gridSize) // cols
	}

	grid[6][7] = 3 // walls
	grid[6][8] = 3 // walls
	grid[5][8] = 3 // walls
	grid[1][2] = 2 // other player
}

func handleInput(client *user) {
	ln := ""
	scanner := bufio.NewScanner(client.conn)
	for scanner.Scan() {
		ln = scanner.Text()
		println(ln)
		//}
		if ln == "UP" {
			client.posY--
			sendUpdate(client.id, ln)
		} else if ln == "DOWN" {
			client.posY++
			sendUpdate(client.id, ln)
		} else if ln == "LEFT" {
			client.posX--
			sendUpdate(client.id, ln)
		} else if ln == "RIGHT" {
			client.posX++
			sendUpdate(client.id, ln)
		} else if ln == "/n" {
			break
		}
	}
}

func sendUpdate(id int, ln string) {
	for i := 0; i < len(clients); i++ {
		if clients[i].id != id {
			io.WriteString(clients[i].conn, ln+"\n")
			time.Sleep(1 * time.Second) // sleep 1 second
		}
	}
}
