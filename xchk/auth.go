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

package xchk

import (
	"fmt"

	"github.com/etcd-io/zetcd"
)

// NewAuth takes a candidate AuthFunc and an oracle AuthFunc
func NewAuth(cAuth, oAuth zetcd.AuthFunc, errc chan<- error) zetcd.AuthFunc {
	sp := newSessionPool()
	return func(zka zetcd.AuthConn) (zetcd.Session, error) {
		s, err := Auth(sp, zka, cAuth, oAuth)
		if _, ok := err.(*XchkError); ok {
			select {
			case errc <- err:
			default:
			}
		}
		return s, err
	}
}

// NewZK takes a candidate ZKFunc and an oracle ZKFunc, returning a cross checker.
func NewZK(cZK, oZK zetcd.ZKFunc, errc chan<- error) zetcd.ZKFunc {
	return func(s zetcd.Session) (zetcd.ZK, error) {
		ss, ok := s.(*session)
		if !ok {
			return nil, fmt.Errorf("expected xchk.session")
		}
		return newZK(ss, cZK, oZK, errc)
	}
}
