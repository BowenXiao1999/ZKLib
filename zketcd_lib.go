// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zetcd

import (
	// "bytes"
	// "encoding/binary"
	// "encoding/gob"
	"fmt"
	// "path"
	// "strings"
	// "time"

	etcd "github.com/coreos/etcd/clientv3"
	"errors"
	// v3sync "github.com/coreos/etcd/clientv3/concurrency"
	// "github.com/golang/glog"
	"context"
)


type ZKClient struct {
	z *zkEtcd
}



func NewZKClient(etcdEps []string) *ZKClient {
	// talk to the etcd3 server
	// create a etcd
	cfg := etcd.Config{Endpoints: etcdEps}
	c, err := etcd.New(cfg)
	if err != nil {
		panic(err)
	}

	// p.authf = zetcd.NewAuth(c)
	// p.ctx = c.Ctx()

	// mock a session here
	s, _ := newSessionForLib(c, 0)

	// zkf := NewZK(c)
	// zk, _ := zkf(s)
	zk := NewZKEtcd(c, s)
	ret := &ZKClient{zk.(*zkEtcd)}
	
	return ret
}


// not support flags, acl yet
func (z *ZKClient) Create(path string, data []byte) error {
	req := &CreateRequest{Path:path, Data:data, Acl:[]ACL{ACL{}}}
	resp := z.z.Create(0, req) // mock a id 0 here
	if resp.Err != nil {
		return errors.New("Create Error")
	}
	fmt.Printf("Error Code %d\n", resp.Hdr.Err)
	return nil
}

func (z *ZKClient) Delete(path string, data []byte) error {
	req := &DeleteRequest{Path:path}
	resp := z.z.Delete(0, req)
	if resp.Err != nil {
		return errors.New("Delete Error")
	}

	return nil
}

func (z *ZKClient) Get(path string) ([]byte, error) {
	req := &GetDataRequest{Path:path}
	resp := z.z.GetData(0, req)
	if resp.Err != nil {
		return resp.Resp.(*GetDataResponse).Data, errors.New("Delete Error")
	}
	if resp.Resp == nil {
		fmt.Printf("Error Code %d\n", resp.Hdr.Err)
		return []byte{}, errors.New("Get Error")
	}

	return resp.Resp.(*GetDataResponse).Data, nil
}

func (z *ZKClient) Set(path string, data []byte) error {
	req := &SetDataRequest{Path:path, Data:data}
	resp := z.z.SetData(0, req)
	if resp.Err != nil {
		return errors.New("Delete Error")
	}

	return nil
}

/*
Hack for Lib Mock a Session
*/
func newSessionForLib(c *etcd.Client, id etcd.LeaseID) (*session, error) {
	ctx, cancel := context.WithCancel(c.Ctx())
	s := &session{id: id, c: c, watches: newWatches(c)}

	_, kaerr := c.KeepAlive(ctx, id)
	if kaerr != nil {
		cancel()
		return nil, kaerr
	}

	// go func() {
	// 	glog.V(9).Infof("starting the session... id=%v", id)
	// 	defer func() {
	// 		glog.V(9).Infof("finishing the session... id=%v; expect revoke...", id)
	// 		cancel()
	// 		s.Close()
	// 	}()
	// 	for {
	// 		select {
	// 		case ka, ok := <-kach:
	// 			if !ok {
	// 				return
	// 			}
	// 			if ka.ResponseHeader == nil {
	// 				continue
	// 			}
	// 			s.mu.Lock()
	// 			s.leaseZXid = ZXid(ka.ResponseHeader.Revision)
	// 			s.mu.Unlock()
	// 		case <-s.StopNotify():
	// 			return
	// 		}
	// 	}
	// }()

	return s, nil
}


