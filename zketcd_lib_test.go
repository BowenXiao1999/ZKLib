package zetcd

import (
	"fmt"
	"testing"

	"time"
)

var (
	Path         string = "/test1"
	SequencePath string = "/test2"
)

func TestCreateAndGetEmphereal(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	Delete(zk, Path, -1, t, ErrNoNode)

	// TODO: Flags Sequence Not Work
	Create(zk, Path, []byte("9999"), 1, []ACL{ACL{}}, t)

	Get(zk, Path, "9999", t)
}

func TestSetAndGet(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	// err = zk.Delete(Path, -1)

	Set(zk, Path, []byte("8888"), -1, t, nil)

	Get(zk, Path, "8888", t)

	//err = zk.Delete(Path, -1)
	Delete(zk, Path, -1, t, ErrNoNode)
}

func TestDeleteAndGet(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	Delete(zk, Path, -1, t, ErrNoNode)

	// expect to be empty and get a log Error Code -101 (Node Not Found)
	Get(zk, Path, "", t)
}

func TestVersionSetAndDelete(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	Create(zk, Path, []byte("9999"), 1, []ACL{ACL{}}, t) // mock ACL and flags

	Set(zk, Path, []byte("8888"), 2, t, ErrBadVersion)

	Set(zk, Path, []byte("8888"), 0, t, nil)

	Delete(zk, Path, 2, t, ErrBadVersion)

	Delete(zk, Path, 1, t, nil)
}

func TestCreateAndGetSequence(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	Delete(zk, SequencePath, -1, t, ErrNoNode)

	// TODO: Flags Sequence Not Work
	Create(zk, SequencePath, []byte("9999"), 1, []ACL{ACL{}}, t)

	//_, err = zk.Create(SequencePath+"/app", []byte("9999"), 2, []ACL{ACL{}}) // mock ACL and flags
	Create(zk, SequencePath+"/app", []byte("9999"), 2, []ACL{ACL{}}, t)

	Get(zk, SequencePath+"/app0000000001", "9999", t)

}

func TestWatches(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})
	zk.setCallBack(testCallBack)
	_ = zk.Delete(Path, -1)
	time.Sleep(100 * time.Millisecond)

	Create(zk, Path, []byte("9999"), 0, []ACL{ACL{}}, t) // mock ACL and flags

	time.Sleep(100 * time.Millisecond)
	// listen
	_, _, err := zk.ExistsW(Path)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(100 * time.Millisecond)

	Delete(zk, Path, -1, t, nil)

	time.Sleep(100 * time.Millisecond)
}

func TestGetChildren(t *testing.T) {
	zk := NewZKClient([]string{"127.0.0.1:2379"})

	Delete(zk, Path, -1, t, ErrNoNode)

	// TODO: Flags Sequence Not Work
	Create(zk, Path, []byte("9999"), 1, []ACL{ACL{}}, t)

	Create(zk, Path+"/child1", []byte("9999"), 1, []ACL{ACL{}}, t)

	resp, _, err := zk.Children(Path) // expect to child1
	if resp[0] != "child1" {
		t.Error(err)
	}

	Delete(zk, Path, -1, t, ErrNoNode)
	Delete(zk, Path+"/child1", -1, t, ErrNoNode)
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
