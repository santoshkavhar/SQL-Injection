// Attack vectors
// ' or ''='
// ' or 1=1 #'
// ' UNION SELECT * from users1 where 1=1 or password='
// ' UNION SELECT @@hostname, @@hostname; #
// ' UNION SELECT @@hostname, load_file('/var/lib/mysql-files/secret_key'); #
// ' or if(database()="a", sleep(10),sleep(20)) and ''='
package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	fmt.Println("SQL Injection Vulnerable Application")

	// Open up our database connection.
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/Security")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	var key string

	fmt.Println("users table will be created. Press any key to continue...")
	fmt.Scanln(&key)

	query := `DROP TABLE IF EXISTS users1;`

	_, err = db.ExecContext(context.Background(), query)
	if err != nil {
		log.Printf("Error %s when deleting users1 table", err)
	}

	query = `CREATE TABLE IF NOT EXISTS users1(username VARCHAR(20) primary key, password VARCHAR(20))`

	_, err = db.ExecContext(context.Background(), query)
	if err != nil {
		log.Printf("Error %s when creating users table", err)
		return
	}

	fmt.Println("3 users will be inserted into users1 table. Press any key to continue...")
	fmt.Scanln(&key)

	// perform a db.Query insert
	insert, err := db.Query(`INSERT INTO users1 VALUES ( 'user1', 'password1' ), ( 'user2', 'password2' ), ( 'user3', 'password3' );`)
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	// get username and password
	var username, password string

	bio := bufio.NewReader(os.Stdin)

	fmt.Println("Enter user credentials to get user details")
	// in case you want a string which doesn't contain the newline
	fmt.Println("Enter username : ")
	line, _, err := bio.ReadLine()
	username = string(line)
	fmt.Println(username)

	fmt.Println("Enter password : ")
	line, _, err = bio.ReadLine()
	password = string(line)
	fmt.Println(password)

	query = fmt.Sprintf("SELECT username, password FROM users1 where username='%s' and password='%s';", username, password)
	fmt.Println(query)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var username string
		var password string

		err = results.Scan(&username, &password)
		if err != nil {
			panic(err.Error())
		}

		log.Println("Users: ", username, password)
	}

}
