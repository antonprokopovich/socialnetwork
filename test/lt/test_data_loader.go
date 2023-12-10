package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"time"
)

type User struct {
	FirstName      string
	LastName       string
	Age            int
	Gender         string
	Interests      string
	City           string
	Email          string
	HashedPassword string
}

var (
	maleFirstNames = []string{
		"Alexander", "Dmitry", "Maxim", "Sergey", "Andrey",
		"Alexey", "Artem", "Ivan", "Nikita", "Mikhail",
		"Vladimir", "Pavel", "Roman", "Viktor", "Evgeny",
		"Denis", "Yuri", "Anatoly", "Igor", "Oleg",
	}
	femaleFirstNames = []string{
		"Anastasia", "Maria", "Olga", "Tatiana", "Elena",
		"Irina", "Ekaterina", "Natalia", "Anna", "Svetlana",
		"Yulia", "Galina", "Ludmila", "Valentina", "Ksenia",
		"Sofia", "Veronika", "Daria", "Alina", "Viktoria",
	}
	surnames = []string{
		"Ivanov", "Smirnov", "Kuznetsov", "Popov", "Vasiliev",
		"Petrov", "Sokolov", "Mikhailov", "Fedorov", "Morozov",
		"Volkov", "Alexeev", "Lebedev", "Semyonov", "Egorov",
		"Pavlov", "Kozlov", "Stepanov", "Nikolaev", "Orlov",
	}

	interests = []string{
		"sports", "music", "traveling", "books", "movies",
		"technology", "cooking", "art", "gaming",
	}
	cities = []string{
		"New York", "Los Angeles", "Chicago", "Houston", "Phoenix",
		"Philadelphia", "San Antonio", "San Diego", "Dallas", "San Jose",
	}
)

func insertUser(db *sql.DB, user User) error {
	_, err := db.Exec("INSERT INTO users (first_name, last_name, age, gender, interests, city, email, hashed_password) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Age, user.Gender, user.Interests, user.City, user.Email, user.HashedPassword)
	return err
}

func generateName(firstNames, surnames []string, isFemale bool) User {
	rand.Seed(time.Now().UnixNano())

	firstName := firstNames[rand.Intn(len(firstNames))]
	surname := surnames[rand.Intn(len(surnames))]
	age := rand.Intn(100)
	gender := "male"

	if isFemale {
		surname += "a"
		gender = "female"
	}

	interest := interests[rand.Intn(len(interests))]
	city := cities[rand.Intn(len(cities))]
	email := fmt.Sprintf("%s.%s@example.com", firstName, surname)
	hashedPassword := fmt.Sprintf("hashed_%d", rand.Int63())

	return User{
		FirstName:      firstName,
		LastName:       surname,
		Age:            age,
		Gender:         gender,
		Interests:      interest,
		City:           city,
		Email:          email,
		HashedPassword: hashedPassword,
	}
}

func main() {
	const numberOfUsers = 500000

	db, err := sql.Open("mysql", "root:root@/socialnetwork?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for i := 0; i < numberOfUsers; i++ {
		user := generateName(maleFirstNames, surnames, false)

		if err = insertUser(db, user); err != nil {
			fmt.Printf("Failed to insert user: %v", err)
		}
	}

	for i := 0; i < numberOfUsers; i++ {
		user := generateName(femaleFirstNames, surnames, true)

		if err = insertUser(db, user); err != nil {
			fmt.Printf("Failed to insert user: %v", err)
		}

		fmt.Printf("inserted user #%d: %v\n", i, user)
	}
}
