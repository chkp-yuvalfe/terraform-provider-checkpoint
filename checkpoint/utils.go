package checkpoint

import (
	"encoding/json"
	"fmt"
	checkpoint "github.com/Checkpoint/api_go_sdk/APIFiles"
	chkp "github.com/Checkpoint/api_go_sdk/APIFiles"
	"io/ioutil"
	"log"
	"os"
)

//var lock sync.Mutex

const (
	FILENAME = "sid.json"
)

type Session struct {
	Sid string `json:"sid"`
	Uid string `json:"uid"`
}

func (s *Session) Save() error {
	f, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(FILENAME, f, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetSession() (Session,error) {
	if _, err := os.Stat(FILENAME); os.IsNotExist(err) {
		_, err := os.Create(FILENAME)
		if err != nil {
			return Session{}, err
		}
	}
	b, err := ioutil.ReadFile(FILENAME)
	if err != nil || len(b) == 0 {
		return Session{}, err
	}
	var s Session
	if err = json.Unmarshal(b, &s); err != nil {
		return Session{}, err
	}
	return s, nil
}

func CheckSession(c *chkp.ApiClient, uid string) bool {
	if uid == "" || c.GetContext() != chkp.WebContext {
		return false
	}
	payload := map[string]interface{}{
		"uid": uid,
	}
	res, _ := c.ApiCall("show-session",payload,c.GetSessionID(),true,false)
	return res.Success
}

func Compare(a, b []string) []string {
	for i := len(a) - 1; i >= 0; i-- {
		for _, vD := range b {
			if a[i] == vD {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}
	return a
}

func PublishAction(client *checkpoint.ApiClient) (bool, error) {
	//lock.Lock()
	//defer lock.Unlock()

	if client.GetAutoPublish() {

		log.Println("publish current session")

		publishRes, _ := client.ApiCall("publish", map[string]interface{}{}, client.GetSessionID(),true,false)
		if !publishRes.Success {
			return false, fmt.Errorf(publishRes.ErrorMsg)
		}
		//time.Sleep(10 * time.Second)
	}
	return true, nil
}