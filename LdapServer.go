package main

import (
	"fmt"
	ldap "github.com/vjeantet/ldapserver"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var resultChan = make(chan string,100)


func main()  {
	//关闭log
	ldap.Logger = log.New(ioutil.Discard,"",0)

	server := ldap.NewServer()

	routes := ldap.NewRouteMux()
	server.Handle(routes)
	routes.Bind(handleBind)
	routes.Search(handleSearch)


	go server.ListenAndServe("127.0.0.1:9998")
	go func() {
		for {
			select {
			case result := <- resultChan:
				fmt.Println(result)
			}
		}
	}()


	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	close(ch)

	server.Stop()

}

func handleBind(w ldap.ResponseWriter, m *ldap.Message)  {
	res := ldap.NewBindResponse(ldap.LDAPResultSuccess)

	w.Write(res)
}

func handleSearch(w ldap.ResponseWriter, m *ldap.Message) {
	r := m.GetSearchRequest()
	resultChan <- string(r.BaseObject())

}



