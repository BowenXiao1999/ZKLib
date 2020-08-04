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

	// "net"
	"github.com/golang/glog"
	"context"
)

type wrapZKEtcd struct {
	*zkEtcd
	cb callback
}

type ZKClient struct {
	z *wrapZKEtcd
}



func NewZKClient(etcdEps []string) *ZKClient {
	// talk to the etcd3 server and new a etcd client
	cfg := etcd.Config{Endpoints: etcdEps}
	c, err := etcd.New(cfg)
	if err != nil {
		panic(err)
	}

	// 1. mock a session and new a ZKEtcd
	s, _ := newSessionForLib(c, 0)

	// // 2. mock a conn (seems not work)
	// cn, _ := net.Dial("tcp", "127.0.0.1:http")
	// con := NewConn(cn)
	// s, _ := newSession(c, con, etcd.LeaseID(0))

	zk := NewZKEtcd(c, s).(*zkEtcd)
	wrapZK := &wrapZKEtcd{zk, emptyCB}
	ret := &ZKClient{wrapZK}
	
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
		return []byte{}, errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}
	if resp.Resp == nil {
		return []byte{}, errors.New("Get Error") // TODO: convert errorCode to Error Type
	}

	return resp.Resp.(*GetDataResponse).Data, nil
}

func (z *ZKClient) Set(path string, data []byte, version int32) (*Stat, error) {
	req := &SetDataRequest{Path:path, Data:data, Version:Ver(version)}
	resp := z.z.SetData(0, req)
	if resp.Err != nil {
		return &Stat{}, resp.Err
	}

	if resp.Hdr.Err != 0 {
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

	// update the leaseID of session if any
	kach, kaerr := c.KeepAlive(ctx, id)
	if kaerr != nil {
		cancel()
		return nil, kaerr
	}

	// do not need session in lib
	go func() {
		glog.V(9).Infof("starting the session... id=%v", id)
		defer func() {
			glog.V(9).Infof("finishing the session... id=%v; expect revoke...", id)
			cancel()

			// s.Close() // hack for conn not exist
			s.watches.close()
		}()
		for {
			select {
			case ka, ok := <-kach:
				if !ok {
					return
				}
				if ka.ResponseHeader == nil {
					continue
				}
				s.mu.Lock()
				s.leaseZXid = ZXid(ka.ResponseHeader.Revision)
				s.mu.Unlock()
			// case <-s.StopNotify():
			// 	return
			}
		}
	}()

	return s, nil
}

// TODO
func (z *ZKClient) Exists(path string) (bool, *Stat, error) {
	req := &ExistsRequest{Path:path}
	resp := z.z.Exists(0, req) // mock a Xid 0 
	if resp.Err != nil {
		return false, &Stat{}, resp.Err 
	}

	if resp.Hdr.Err != 0 {
		return false, &Stat{}, errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}
	// return resp.Resp.(statResponse), nil
	return true, &resp.Resp.(*ExistsResponse).Stat, nil
}
func (z *ZKClient) ExistsW(path string) (bool, *Stat, error)  {

	// TODO: do SetWatchRequest
	req := &SetWatchesRequest{DataWatches:[]string{path}, ExistWatches:[]string{path}}
	resp := z.z.SetWatches(0, req) // mock a Xid 0 
	if resp.Err != nil {
		return false, &Stat{}, resp.Err 
	}

	if resp.Hdr.Err != 0 {
		return false, &Stat{}, errorCodeToErr[ErrCode(resp.Hdr.Err)]
	}
	// return resp.Resp.(statResponse), nil
	// return true, &resp.Resp.(*SetWatchesResponse).Stat, nil
	return true, &Stat{}, nil
}

func (z *wrapZKEtcd) SetWatches(xid Xid, op *SetWatchesRequest) ZKResponse {

	for _, dw := range op.DataWatches {
		dataPath := dw
		p := mkPath(dataPath)
		f := func(newzxid ZXid, evt EventType) {
			wresp := &WatcherEvent{
				Type:  evt,
				State: StateSyncConnected,
				Path:  dataPath,
			}
			glog.V(7).Infof("WatchData* (%v,%v,%v)", xid, newzxid, *wresp)
			// z.s.Send(-1, -1, wresp)
			// TODO: invoke the callback instead of sending resp in connection
			z.cb(wresp)
		}
		z.s.Watch(op.RelativeZxid, xid, p, EventNodeDataChanged, f)
	}

	ops := make([]etcd.Op, len(op.ExistWatches))
	for i, ew := range op.ExistWatches {
		ops[i] = etcd.OpGet(
			mkPathVer(mkPath(ew)),
			etcd.WithSerializable(),
			etcd.WithRev(int64(op.RelativeZxid)))
	}

	resp, err := z.c.Txn(z.c.Ctx()).Then(ops...).Commit()
	if err != nil {
		return mkErr(err)
	}
	curZXid := ZXid(resp.Header.Revision)

	for i, ew := range op.ExistWatches {
		existPath := ew
		p := mkPath(existPath)

		ev := EventNodeDeleted
		if len(resp.Responses[i].GetResponseRange().Kvs) == 0 {
			ev = EventNodeCreated
		}
		f := func(newzxid ZXid, evt EventType) {
			wresp := &WatcherEvent{
				Type:  evt,
				State: StateSyncConnected,
				Path:  existPath,
			}
			glog.V(7).Infof("WatchExist* (%v,%v,%v)", xid, newzxid, *wresp)
			// z.s.Send(-1, -1, wresp)
			z.cb(wresp)
		}
		z.s.Watch(op.RelativeZxid, xid, p, ev, f)
	}
	for _, cw := range op.ChildWatches {
		childPath := cw
		p := mkPath(childPath)
		f := func(newzxid ZXid, evt EventType) {
			wresp := &WatcherEvent{
				Type:  EventNodeChildrenChanged,
				State: StateSyncConnected,
				Path:  childPath,
			}
			glog.V(7).Infof("WatchChild* (%v,%v,%v)", xid, newzxid, *wresp)
			// z.s.Send(-1, -1, wresp)
			z.cb(wresp)
		}
		z.s.Watch(op.RelativeZxid, xid, p, EventNodeChildrenChanged, f)
	}

	swresp := &SetWatchesResponse{}

	glog.V(7).Infof("SetWatches(%v) = (zxid=%v, resp=%+v)", xid, curZXid, *swresp)
	return mkZKResp(xid, curZXid, swresp)
}

func (z *ZKClient) setCallBack(c callback) {
	z.z.cb = c;
}

type callback func(w *WatcherEvent) 

func emptyCB(w *WatcherEvent)  {
	
}

