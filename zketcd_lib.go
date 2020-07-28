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
	// talk to the etcd3 server and new a etcd client
	cfg := etcd.Config{Endpoints: etcdEps}
	c, err := etcd.New(cfg)
	if err != nil {
		panic(err)
	}

	// mock a session and new a ZKEtcd
	s, _ := newSessionForLib(c, 0)
	zk := NewZKEtcd(c, s)
	ret := &ZKClient{zk.(*zkEtcd)}
	
	return ret
}


// not support flags, acl yet
func (z *ZKClient) Create(path string, data []byte, flags int32, acl []ACL) (string, error) {
	req := &CreateRequest{Path:path, Data:data, Acl:acl, Flags:flags}
	resp := z.z.Create(0, req) // mock a Xid 0 
	if resp.Err != nil {
		return "", resp.Err 
	}

	if resp.Hdr.Err != 0 {
		fmt.Printf("Error Code %d\n", resp.Hdr.Err)
		return "", errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}
	return resp.Resp.(*CreateResponse).Path, nil
}

func (z *ZKClient) Delete(path string, version int32) error {
	req := &DeleteRequest{Path:path,  Version:Ver(version)}
	resp := z.z.Delete(0, req)
	if resp.Err != nil {
		return resp.Err
	}

	if resp.Hdr.Err != 0 {
		fmt.Printf("Error Code %d\n", resp.Hdr.Err)
		return errorCodeToErr[ErrCode(resp.Hdr.Err)] 
	}

	return nil
}

func (z *ZKClient) Get(path string) ([]byte, error) {
	req := &GetDataRequest{Path:path}
	resp := z.z.GetData(0, req)
	if resp.Err != nil {
		return []byte{}, resp.Err
	}
	if resp.Hdr.Err != 0 {
		fmt.Printf("Error Code %d\n", resp.Hdr.Err)
		return []byte{}, errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}
	if resp.Resp == nil {
		return []byte{}, errors.New("Get Error") // TODO: convert errorCode to Error Type
	}

	return resp.Resp.(*GetDataResponse).Data, nil
}

func (z *ZKClient) Set(path string, data []byte, version int32) (*Stat, error) {
	req := &SetDataRequest{Path:path, Data:data}
	resp := z.z.SetData(0, req)
	if resp.Err != nil {
		return &Stat{}, resp.Err
	}

	if resp.Hdr.Err != 0 {
		fmt.Printf("Error Code %d\n", resp.Hdr.Err)
		return &Stat{}, errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}

	return &resp.Resp.(*SetDataResponse).Stat, nil
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

	// // do not need session in lib
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


