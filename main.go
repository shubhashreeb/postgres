package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

type MyTableRow struct {
	Id   string
	Msg  string
	Data interface{}
}

func main() {
	fmt.Println("Running db methods")
	listenForTableChanges()
}

func listenForTableChanges() {
	dbConfig := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"default",
		"postgres")

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	l := pq.NewListener(dbConfig, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	})
	defer l.Close()

	err = l.Listen("mytable_changes")
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		select {
		case n := <-l.Notify:
			parts := strings.Split(n.Extra, "|")
			log.Printf("*** Received data %+v\n", parts)
			op := parts[0]
			id, _ := strconv.Atoi(parts[1])
			data := parts[2]
			log.Printf("*** Received data %+v\n", data)

			switch op {
			case "INSERT":
				// handle insert event
				var row MyTableRow
				err := json.Unmarshal([]byte(data), &row)
				if err != nil {
					log.Println(err.Error())
				}
				log.Printf("Inserted row with id %d and data %+v\n", id, row)
			case "UPDATE":
				// handle update event
				var row MyTableRow
				err := json.Unmarshal([]byte(data), &row)
				if err != nil {
					log.Println(err.Error())
				}
				log.Printf("Updated row with id %d and new data %+v\n", id, row)
			case "DELETE":
				// handle delete event
				log.Printf("Deleted row with id %d\n", id)
			}
			// case <-time.After(90 * time.Second):
			// 	log.Println("Listener timeout. No events received for 90 seconds.")
			// 	return
		}
	}
}
