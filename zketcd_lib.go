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
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"path"
	"strings"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	v3sync "github.com/coreos/etcd/clientv3/concurrency"
	"github.com/golang/glog"
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

	// create a session here
	s Session;

	zkf := zetcd.NewZK(c)
	ret := &ZKClient{zkf(s)}
	
	return ret
}


func (z *ZKClient) Create(path string, data []byte) error {
	req := &CreateRequest{Path:path, Data:data}
	resp := z.z.Create(0, creq)
	if resp.Err != nil {
		return error
	}

	return nil
}

func (z *ZKClient) Delete(path string, data []byte) error {
	req := &DeleteRequest{Path:path, Data:data}
	resp := z.z.Delete(0, creq)
	if resp.Err != nil {
		return error
	}

	return nil
}

func (z *ZKClient) GetData(path string, data []byte) error {
	req := &GetDataRequest{Path:path, Data:data}
	resp := z.z.GetData(0, creq)
	if resp.Err != nil {
		return error
	}

	return nil
}

func (z *ZKClient) SetData(path string, data []byte) error {
	req := &DeleteRequest{Path:path, Data:data}
	resp := z.z.SetData(0, creq)
	if resp.Err != nil {
		return error
	}

	return nil
}


