package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LassiHeikkila/flmnchll/content/contentdb"
)

func main() {
	var (
		dbPath string
	)

	flag.StringVar(&dbPath, "db", "content.db", "Path to database file")

	flag.Parse()

	if err := contentdb.Connect(dbPath); err != nil {
		fmt.Println("error connecting to database:", err)
		return
	}

	if err := contentdb.Init(); err != nil {
		fmt.Println("error initializing database:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin)

	prompt := func() {
		fmt.Println("What would you like to do?")
		fmt.Println("Options are:")
		fmt.Println("\tcreate (c)")
		fmt.Println("\tread   (r)")
		fmt.Println("\tupdate (u)")
		fmt.Println("\tdelete (d)")
		fmt.Println()
		fmt.Println("\tquit   (q)")
	}

	prompt()
	for scanner.Scan() {
		switch scanner.Text() {
		case "c", "create":
			handleCreate(scanner)
		case "r", "read":
			handleRead(scanner)
		case "u", "update":
			handleUpdate(scanner)
		case "d", "delete":
			handleDelete(scanner)
		case "q", "quit":
			fmt.Println("bye!")
			return
		default:
			fmt.Println("invalid option")
		}
		prompt()
	}
}

func handleCreate(scanner *bufio.Scanner) {
	fmt.Println("input video object:")
	scanner.Scan()
	var v contentdb.Video
	if err := json.Unmarshal(scanner.Bytes(), &v); err != nil {
		fmt.Println("failed to parse input:", err)
		return
	}

	if id, err := contentdb.AddVideo(v); err != nil {
		fmt.Println("failed to add video to database:", err)
		return
	} else {
		fmt.Println("inserted video with id:", id)
	}
}

func handleRead(scanner *bufio.Scanner) {
	fmt.Println("input file id:")
	scanner.Scan()
	id := scanner.Text()

	v, err := contentdb.GetVideo(id)
	if err != nil {
		fmt.Println("error getting video from database:", err)
		return
	}

	b, err := json.Marshal(&v)
	if err != nil {
		fmt.Println("failed to marshal video struct:", err)
		return
	}

	fmt.Println(string(b))
}

func handleUpdate(scanner *bufio.Scanner) {
	fmt.Println("update is not supported yet")
}

func handleDelete(scanner *bufio.Scanner) {
	fmt.Println("input file id:")
	scanner.Scan()
	id := scanner.Text()

	if err := contentdb.DeleteVideo(id); err != nil {
		fmt.Println("error deleting video from database:", err)
	}
}
