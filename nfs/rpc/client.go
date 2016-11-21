package rpc

import (
	"bufio"
	"bytes"
	"fmt"
	"net"

	"nfs/xdr"
)

type Client struct {
	transport
}

func DialTCP(network, addr string) (*Client, error) {
	a, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP(a.Network(), nil, a)
	if err != nil {
		return nil, err
	}
	t := &tcpTransport{
		Reader:      bufio.NewReader(conn),
		WriteCloser: conn,
	}
	return &Client{t}, nil
}

func (c *Client) Call(call interface{}) ([]byte, error) {
	msg := &message{
		Xid:     0xcafebabe,
		Msgtype: 0,
		Body:    call,
	}
	w := new(bytes.Buffer)
	if err := xdr.Write(w, msg); err != nil {
		return nil, err
	}
	if err := c.send(w.Bytes()); err != nil {
		return nil, err
	}
	buf, err := c.recv()
	if err != nil {
		return nil, err
	}
	xid, buf := xdr.Uint32(buf)
	if xid != msg.Xid {
		return nil, fmt.Errorf("xid did not match, expected: %x, received: %x", msg.Xid, xid)
	}
	mtype, buf := xdr.Uint32(buf)
	if mtype != 1 {
		return nil, fmt.Errorf("message as not a reply: %d", mtype)
	}
	reply_stat, buf := xdr.Uint32(buf)
	switch reply_stat {
	case MSG_ACCEPTED:
		_, buf = xdr.Uint32(buf)
		opaque_len, buf := xdr.Uint32(buf)
		_ = buf[0:int(opaque_len)]
		buf = buf[opaque_len:]
		accept_stat, buf := xdr.Uint32(buf)
		switch accept_stat {
		case SUCCESS:
			return buf, nil
		case PROG_UNAVAIL:
			return nil, fmt.Errorf("PROG_UNAVAIL")
		case PROG_MISMATCH:
			// TODO(dfc) decode mismatch_info
			return nil, fmt.Errorf("rpc: PROG_MISMATCH")
		default:
			return nil, fmt.Errorf("rpc: %d", accept_stat)
		}
	case MSG_DENIED:
		rejected_stat, _ := xdr.Uint32(buf)
		switch rejected_stat {
		case RPC_MISMATCH:

		default:
			return nil, fmt.Errorf("rejected_stat was not valid: %d", rejected_stat)
		}
	default:
		return nil, fmt.Errorf("reply_stat was not valid: %d", reply_stat)
	}
	panic("unreachable")
}
