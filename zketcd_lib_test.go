package zetcd

import (
	"testing"
	"fmt"
)


func TestCreateAndGet(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	_, err := zk.Create("/test", []byte("9999"), 0, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error()
	}

	resp, err := zk.Get("/test")	
	// expect to be 9999
	if string(resp) != "9999" {
		fmt.Printf("Not Equal %s\n", string(resp))
	}else{

	}

}

func TestSetAndGet(t *testing.T)  {

	zk := NewZKClient([]string{"127.0.0.1:2379"})
	_, err := zk.Set("/test", []byte("8888"), -1)
	if err != nil {
		t.Error()
	}

	resp, err := zk.Get("/test")	
	// expect to be 9999
	if string(resp) != "8888" {
		fmt.Printf("Not Equal %s\n", string(resp))
	}else{
		fmt.Printf("Equal %s\n", string(resp))
	}
}

func TestDeleteAndGet(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	err := zk.Delete("/test", -1)
	if err != nil {
		t.Error()
	}

	resp, err := zk.Get("/test")	
	// expect to be 9999
	if string(resp) != "" {
		fmt.Printf("Not Equal %s\n", string(resp))
	}else{
		fmt.Printf("Equal %s\n", string(resp))
	}
}