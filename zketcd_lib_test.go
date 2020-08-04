package zetcd

import (
	"testing"
	"fmt"
)

var (
	Path string = "/test1"
	SequencePath string = "/test2"

)

func TestCreateAndGetEmphereal(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	_ = zk.Delete(Path, -1)

	// TODO: Flags Sequence Not Work
	_, err := zk.Create(Path, []byte("9999"), 1, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get(Path)	
	// expect to be 9999
	if string(resp) != "9999" {
		t.Error(err)
	}

}

func TestSetAndGet(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	// err = zk.Delete(Path, -1)

	_, err := zk.Set(Path, []byte("8888"), -1)
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get(Path)	
	// expect to be 8888
	if string(resp) != "8888" {
		t.Error(err)
	}

	err = zk.Delete(Path, -1)
}

func TestDeleteAndGet(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	err := zk.Delete(Path, -1)
	if err != ErrNoNode {
		t.Error(err)
	}

	resp, err := zk.Get(Path)	
	// expect to be empty and get a log Error Code -101 (Node Not Found)
	if string(resp) != "" {
		t.Error(err)
	}
}

func TestVersionSetAndDelete(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	_, err := zk.Create(Path, []byte("9999"), 1, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error(err)
	}
	
	_, err = zk.Set(Path, []byte("8888"), 2)
	if err != ErrBadVersion {
		t.Error(err)
	}

	_, err = zk.Set(Path, []byte("8888"), 0)
	if err != nil {
		t.Error(err)
	}

	// expect to get a ErrBadVersion
	err = zk.Delete(Path, 2)	
	if err != ErrBadVersion {
		t.Error(err)
	}


	err = zk.Delete(Path, 1)	
	if err != nil {
		t.Error(err)
	}
}

func TestCreateAndGetSequence(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	err := zk.Delete(SequencePath, -1)	


	// TODO: Flags Sequence Not Work
	_, err = zk.Create(SequencePath, []byte("9999"), 1, []ACL{ACL{}}) // mock ACL and flags
	// if err != nil {
	// 	t.Error(err)
	// }

	_, err = zk.Create(SequencePath+"/app", []byte("9999"), 2, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error(err)
	}

	resp, err := zk.Get(SequencePath + "/app0000000001")	
	// expect to be 9999
	if string(resp) != "9999" {
		t.Error(err)
	}

}

func TestWatches(t *testing.T)  {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	zk.setCallBack(testCallBack)
	_ = zk.Delete(Path, -1)

	_, _, err := zk.ExistsW(Path)
	if err != nil {
		t.Error(err)
	}
	// TODO: Flags Sequence Not Work
	_, err = zk.Create(Path, []byte("9999"), 0, []ACL{ACL{}}) // mock ACL and flags
	if err != nil {
		t.Error(err)
	}


	// listen
	_, _, err = zk.ExistsW(Path)
	if err != nil {
		t.Error(err)
	}


	_ = zk.Delete(Path, -1)
}

// zk watch 回调函数
func testCallBack(event *WatcherEvent) {
	// zk.EventNodeCreated
	// zk.EventNodeDeleted
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type)
	fmt.Println("state: ", event.State)
	fmt.Println("---------------------------")
}