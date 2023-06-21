// Copyright 2022 LiDonghai Email:5241871@qq.com Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package broker

import (
	"errors"
	"github.com/Coolpy7/Coolpy7_Community/mempool"
	"github.com/Coolpy7/Coolpy7_Community/packet"
)

func (c *Client) send(pkt packet.Generic) error {
	buf := mempool.Malloc(pkt.Len())
	defer mempool.Free(buf)
	n, _ := pkt.Encode(buf)
	if n != len(buf) {
		return errors.New("encode err")
	}
	if !c.isTls {
		_, err := c.conn.Write(buf)
		if err != nil {
			return err
		}
	} else {
		_, err := c.tlsConn.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}
