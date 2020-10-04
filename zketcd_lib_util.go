package zetcd

import "testing"

func Create(zk *ZKClient, Path string, data []byte, flags int32, acl []ACL, t *testing.T) {
	_, err := zk.Create(Path, data, flags, acl) // mock ACL and flags
	if err != ErrNodeExists && err != nil {
		t.Error(err)
	}
}

func Set(zk *ZKClient, Path string, data []byte, version int32, t *testing.T, e error) {
	_, err := zk.Set(Path, data, version)
	if err != e {
		t.Error(err)
	}
}

func Delete(zk *ZKClient, Path string, version int32, t *testing.T, e error) {
	err := zk.Delete(Path, version)
	if e != ErrNoNode && err != e {
		t.Error(err)
	}
}

func Get(zk *ZKClient, Path string, expect string, t *testing.T) {
	resp, err := zk.Get(Path)
	if string(resp) != expect {
		t.Error(err)
	}
}
