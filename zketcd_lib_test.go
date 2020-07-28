package zetcd

import (
	"testing"
)


func TestCreateAndGet(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	_, err := zk.Create("/test", []byte("9999"), 0, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get("/test")	
	// expect to be 9999
	if string(resp) != "9999" {
		t.Error(err)
	}

}

func TestSetAndGet(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	_, err := zk.Set("/test", []byte("8888"), -1)
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get("/test")	
	// expect to be 8888
	if string(resp) != "8888" {
		t.Error(err)
	}
}

func TestDeleteAndGet(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	err := zk.Delete("/test", -1)
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get("/test")	
	// expect to be empty and get a log Error Code -101 (Node Not Found)
	if string(resp) != "" {
		t.Error(err)
	}
}