package zetcd

import (
	"testing"
	"fmt"
)


func TestCreateAndGet(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	err := zk.Create("/test", []byte("9999"))
	if err != nil {
		t.Error()
	}

	resp, err := zk.Get("/test")	
	// expect to be 9999
	if string(resp) != "9999" {
		fmt.Printf("Not Equal %s\n", string(resp))
	}else{

	}
	
	// err := zk.Set("/hbase/meta-region-server", "8888")
	// if err != nil {
		
	// }

	// resp, err := zk.Get("/hbase/meta-region-server")
	// // expect to be 8888

	// err := zk.Delete("/hbase/meta-region-server")
	// if err != nil {

	// }

	// resp, err := zk.Get("/hbase/meta-region-server")
	// // expect to be 8888

}