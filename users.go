package main

import (
	"net"
	"sync"

	"git.sr.ht/~runxiyu/meseircd/meselog"
)

type User struct {
	Clients []*Client
	UID     uint64
	Nick    string
	Ident   string
	Gecos   string
	Host    string
	Caps    map[string]struct{}
	Extra   map[string]any
	Server  *Server
	State   ClientState
}

func (user *User) SendToLocalClients(msg SMsg) (numSent uint) {
	for _, c := range user.Clients {
		if c.Server != self {
			continue
		}
		err := c.Send(msg)
		if err == nil {
			numSent++
		}
	}
	return
}

func (user *User) ClientSource() string {
	// TODO: Edge cases where these aren't available
	return user.Nick + "!" + user.Ident + "@" + user.Host
}

func (user *User) ServerSource() uint64 {
	return user.UID
}

// func (user *User) Delete() {
// 	if client.conn != nil {
// 		(*client.conn).Close()
// 	}
// 	if !cidToClient.CompareAndDelete(client.CID, client) {
// 		meselog.Error("cid inconsistent", "cid", client.CID, "client", client)
// 	}
// 	if client.State >= ClientStateRegistered || client.Nick != "*" {
// 		if !nickToClient.CompareAndDelete(client.Nick, client) {
// 			meselog.Error("nick inconsistent", "nick", client.Nick, "client", client)
// 		}
// 	}
// }

func NewLocalUser(conn *net.Conn) (*User, error) {
	var uidPart uint32
	{
		uidPartCountLock.Lock()
		defer uidPartCountLock.Unlock()
		if uidPartCount == ^uint32(0) { // UINT32_MAX
			return nil, ErrFullClients
		}
		uidPartCount++
		uidPart = uidPartCount
	}
	client := &Client{
		conn:   conn,
		Server: self,
		State:  ClientStatePreRegistration,
		Nick:   "*",
		Caps:   make(map[string]struct{}),
		Extra:  make(map[string]any),
		CID:    uint64(self.SID)<<32 | uint64(uidPart),
	}
	return client, nil
}

const (
	ClientStatePreRegistration ClientState = iota
	ClientStateCapabilities
	ClientStateCapabilitiesFinished
	ClientStateRegistered
	ClientStateRemote
)

var (
	cidToClient      = sync.Map{}
	nickToClient     = sync.Map{}
	uidPartCount     uint32
	uidPartCountLock sync.Mutex
)
