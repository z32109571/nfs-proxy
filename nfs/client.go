package nfs

import (
	"nfs/rpc"
)

type Client struct {
	*rpc.Client
}
