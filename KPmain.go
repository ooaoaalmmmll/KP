package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// __________________________________________________________________________________________________________________ITEMS FUNCTIONALITY______________________________________
type Item struct {
	Id      int `json:"id" form:"id"`
	Address string `json:"address" form:"address"`
	Date    string `json:"date" form:"date"`
	Dims    string `json:"dims" form:"dims"` //dimentions + volume
	Desc    string `json:"desc" form:"desc"`
	OwnerId int `json:"ownerId" form:"ownerId"`
}

func (p *Item) AddItem() (id int64, err error) {
	rs, err := db.Exec("INSERT INTO item(address, date, dims, desc, ownerId) VALUES (?, ?, ?, ?, ?)", p.Address, p.Date, p.Desc, p.Dims, p.OwnerId)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

func (p *Item) GetItems() (items []Item, err error) {
	items = make([]Item, 0)
	rows, err := db.Query("SELECT id, address, date, dims, desc, ownerId FROM item")
	if err != nil {
		return
	}

	defer rows.Close()

	
	for rows.Next() {
		var item Item
		rows.Scan(&item.Id, &item.Address, &item.Date, &item.Desc, &item.Dims, &item.OwnerId)
		items = append(items, item)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (p *Item) GetItem() (item Item, err error) {
	err = db.QueryRow("SELECT id, address, date, dims, desc, ownerId FROM item WHERE id=?", p.Id).Scan(
		&item.Id, &item.Address, &item.Date, &item.Desc, &item.Dims, &item.OwnerId,
	)
	return
}

func (p *Item) ModItem() (ra int64, err error) {
	stmt, err := db.Prepare("UPDATE item SET address=?, date=?, dims=?, desc=?, ownerId=? WHERE id=?")
	if err != nil {
		return
	}
	defer stmt.Close()
	
	rs, err := stmt.Exec(p.Address, p.Date, p.Desc, p.Dims, p.OwnerId, p.Id)
	if err != nil {
		return
	}
	ra, err = rs.RowsAffected()
	return
}

func (p *Item) DelItem() (ra int64, err error) {
	rs, err := db.Exec("DELETE FROM item WHERE id=?", p.Id)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err = rs.RowsAffected()
	return
}


func GetRidOfReturnedErr(i int, e error) (int){				//	<----------------------<<	
	return i;
}
	//																					^		
func AddItemApi(c *gin.Context) {
	address := c.Request.FormValue("address")
	//																					^	
	date := c.Request.FormValue("date")
	//																					^	
	dims := c.Request.FormValue("dims")
	//																					^	
	desc := c.Request.FormValue("desc")
	//																					^																	
	ownerId := c.Request.FormValue("ownerId")
	//																					^																														
	p := Item{Address: address, Date: date, Dims: dims, Desc: desc, OwnerId:  GetRidOfReturnedErr(strconv.Atoi(ownerId))}

	ra, err := p.AddItem()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("insert successful %d", ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func GetItemsApi(c *gin.Context) {
	var p Item
	items, err := p.GetItems()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})

}

func GetItemApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Item{Id: id}
	item, err := p.GetItem()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"item": item,
	})

}

func ModItemApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Item{Id: id}
	err = c.Bind(&p)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err := p.ModItem()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Update item %d successful %d", p.Id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func DelItemApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Item{Id: id}
	ra, err := p.DelItem()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Delete item %d successful %d", id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

/*

//__________________________________________________________________________________________________________________PERSONS FUNCTIONALITY______________________________________

type Person struct {
	Id        int    `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
}

func (p *Person) AddPerson() (id int64, err error) {
	rs, err := db.Exec("INSERT INTO person(first_name, last_name) VALUES (?, ?)", p.FirstName, p.LastName)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

func (p *Person) GetPersons() (persons []Person, err error) {
	persons = make([]Person, 0)
	rows, err := db.Query("SELECT id, first_name, last_name FROM person")
	if err != nil {
		return
	}
	defer rows.Close()

	

	for rows.Next() {
		var person Person
		rows.Scan(&person.Id, &person.FirstName, &person.LastName)
		persons = append(persons, person)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return
}

func (p *Person) GetPerson() (person Person, err error) {
	err = db.QueryRow("SELECT id, first_name, last_name FROM person WHERE id=?", p.Id).Scan(
		&person.Id, &person.FirstName, &person.LastName,
	)
	return
}

func (p *Person) ModPerson() (ra int64, err error) {
	stmt, err := db.Prepare("UPDATE person SET first_name=?, last_name=? WHERE id=?")
	if err != nil {
		return
	}
	defer stmt.Close()
	
	rs, err := stmt.Exec(p.FirstName, p.LastName, p.Id)
	if err != nil {
		return
	}
	ra, err = rs.RowsAffected()
	return
}

func (p *Person) DelPerson() (ra int64, err error) {
	rs, err := db.Exec("DELETE FROM person WHERE id=?", p.Id)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err = rs.RowsAffected()
	return
}

func IndexApi(c *gin.Context) {
	c.String(http.StatusOK, "It works")
}

func AddPersonApi(c *gin.Context) {
	firstName := c.Request.FormValue("first_name")
	lastName := c.Request.FormValue("last_name")

	p := Person{FirstName: firstName, LastName: lastName}

	ra, err := p.AddPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("insert successful %d", ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func GetPersonsApi(c *gin.Context) {
	var p Person
	persons, err := p.GetPersons()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"persons": persons,
	})

}

func GetPersonApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	person, err := p.GetPerson()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"person": person,
	})

}

func ModPersonApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	err = c.Bind(&p)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err := p.ModPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Update person %d successful %d", p.Id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func DelPersonApi(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	ra, err := p.DelPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Delete person %d successful %d", id, ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}
*/

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()
	router.GET("/", IndexApi)

//	router.GET("/persons", GetPersonsApi)

//	router.GET("/person/:id", GetPersonApi)
//	router.POST("/person", AddPersonApi)
//	router.PUT("/person/:id", ModPersonApi)
//	router.DELETE("/person/:id", DelPersonApi)
 
	router.GET("/items", GetItemsApi)
	router.GET("/item/:ownerId", GetItemsByOwnerApi)
	router.GET("/item/:id", GetItemApi)
	router.POST("/item", AddItemApi)
	router.PUT("/item/:id", ModItemApi)
	router.DELETE("/item/:id", DelItemApi)

	

	router.Run(":8000")
}
