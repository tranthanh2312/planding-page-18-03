package db

import (
	"log"

	"github.com/fatih/color"
)

type AllUser struct {
	Users []User `json:"users"`
}

func (q *Querier) GetAllUser() (res AllUser, err error) {
	bolt, err := Open()
	if err != nil {
		return
	}

	defer bolt.Close()

	err = bolt.Get("users", &res.Users)
	if err == nil && res.Users != nil {
		log.Println("Cache is hit!")
		return res, err
	}
	res.Users, err = q.GetAllUsers()
	if err != nil {
		return
	}

	err = bolt.Set("users", res)
	if err != nil {
		log.Fatal("Cannot cache users:", err)
		return
	}

	color.Blue("Cache is successfully updated!")

	return res, err
}

func (a *AllUser) SearchUser(userId int) (index int) {
	// binary search
	low := 0
	high := len(a.Users) - 1

	for low <= high {
		mid := (low + high) / 2
		if a.Users[mid].ID == userId {
			return mid
		} else if a.Users[mid].ID < userId {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return -1
}

func (a *AllUser) SearchUserFullName(userId int) string {
	// binary search
	low := 0
	high := len(a.Users) - 1

	for low <= high {
		mid := (low + high) / 2
		if a.Users[mid].ID == userId {
			return a.Users[mid].FullName
		} else if a.Users[mid].ID < userId {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	log.Println("out")

	return ""
}
