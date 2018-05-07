package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
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
	setupGrid()

	clientID := 1
	ln, err := net.Listen("tcp", ":9000")
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
		io.WriteString(client.conn, strconv.Itoa(gridSize)+"/"+strconv.Itoa(client.id)+"\n")
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
	scanner.Split(customSplitFunc)
	for scanner.Scan() {
		ln = scanner.Text()
		if ln == "UP" {
			client.posY--
		} else if ln == "DOWN" {
			client.posY++
		} else if ln == "LEFT" {
			client.posX--
		} else if ln == "RIGHT" {
			client.posX++
		}
		sendUpdate(client.id, ln)
	}
}

func sendUpdate(id int, ln string) {
	for i := 0; i < len(clients); i++ {
		if clients[i].id != id {
			io.WriteString(clients[i].conn, strconv.Itoa(clients[i].id)+"/"+ln+"/"+"\n")
		}
	}
}

func customSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a newline followed by a
	// pound sign.
	if i := strings.Index(string(data), "/"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}
