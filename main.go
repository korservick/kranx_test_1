package main

import (
	"encoding/json"
	"github.com/labstack/echo"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type Account struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Balance int64  `json:"balance"`
}

type Accounts struct {
	ID    map[string]*Account `json:"id"`
	Mutex sync.RWMutex        `json:"-"`
}

type Transaction struct {
	FromID string `json:"from_id"`
	ToID   string `json:"to_id"`
	Amount int64  `json:"amount"`
}

var accounts *Accounts

func accountList(c echo.Context) error {
	accounts.Mutex.RLock()
	a, err := json.Marshal(accounts)
	accounts.Mutex.RUnlock()
	if err != nil {
		log.Println(err)
	}
	return c.JSONBlob(http.StatusOK, a)
}

func accountGet(c echo.Context) error {
	id := c.Param("id")
	accounts.Mutex.RLock()
	a, err := json.Marshal(accounts.ID[id])
	accounts.Mutex.RUnlock()
	if err != nil {
		log.Println(err)
	}
	return c.JSONBlob(http.StatusOK, a)
}

func accountAdd(c echo.Context) error {
	account := new(Account)

	bodyBuffer, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err.Error())
		return c.String(http.StatusConflict, err.Error())
	} else {
		err = json.Unmarshal(bodyBuffer, &account)
		if err != nil {
			log.Println(err.Error())
			return c.String(http.StatusConflict, err.Error())
		} else {
			if account.ID == "" {
				return c.String(http.StatusConflict, "ID = ''")
			}
			if account.User == "" {
				return c.String(http.StatusConflict, "User = ''")
			}
			accounts.Mutex.Lock()
			accounts.ID[account.ID] = account
			accounts.Mutex.Unlock()
		}
	}
	return c.HTMLBlob(http.StatusOK, bodyBuffer)
}

func transaction(c echo.Context) error {
	transaction := new(Transaction)

	bodyBuffer, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err.Error())
	} else {
		err = json.Unmarshal(bodyBuffer, &transaction)
		if err != nil {
			log.Println(err)
		} else {
			accounts.Mutex.RLock()
			if accounts.ID[transaction.FromID] == nil {
				accounts.Mutex.RUnlock()
				return c.String(http.StatusConflict, "from_id not found")
			}
			if accounts.ID[transaction.ToID] == nil {
				accounts.Mutex.RUnlock()
				return c.String(http.StatusConflict, "to_id not found")
			}
			if accounts.ID[transaction.FromID].Balance < transaction.Amount {
				accounts.Mutex.RUnlock()
				return c.String(http.StatusConflict, "Balance less then amount")
			}
			accounts.Mutex.RUnlock()
			accounts.Mutex.Lock()
			accounts.ID[transaction.FromID].Balance -= transaction.Amount
			accounts.ID[transaction.ToID].Balance += transaction.Amount
			accounts.Mutex.Unlock()

		}
	}
	return c.HTMLBlob(http.StatusOK, bodyBuffer)
}

func main() {
	accounts = new(Accounts)
	accounts.ID = make(map[string]*Account)

	e := echo.New()
	e.HideBanner = true
	e.GET("/account/", accountList)
	e.GET("/account/:id", accountGet)
	e.POST("/account/", accountAdd)
	e.POST("/transaction/", transaction)
	e.Logger.Fatal(e.Start(":8080"))
}
