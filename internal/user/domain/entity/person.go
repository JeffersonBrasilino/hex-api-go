package entity

import "github.com/hex-api-go/pkg/core/domain"

type Person struct {
	domain.Entity
	name       string
	age        int
	birthDate  string
	email      string
	dataSource string
}

func NewPerson(name string, age int, birthDate string, email string, dataSource string) *Person {
	return &Person{
		name:       name,
		age:        age,
		birthDate:  birthDate,
		email:      email,
		dataSource: dataSource,
		Entity:     domain.NewEntity(""),
	}
}

func (p *Person) GetName() string {
	return p.name
}

func (p *Person) GetAge() int {
	return p.age
}

func (p *Person) GetBirthDate() string {
	return p.birthDate
}

func (p *Person) GetEmail() string {
	return p.email
}

func (p *Person) GetDataSource() string {
	return p.dataSource
}
