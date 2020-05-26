package ddnet

import (
	"encoding/hex"
	"log"
	"net"
	"time"
)

/*
//
laddr := net.UDPAddr{
    IP:   net.IPv4(192, 168, 137, 224),
    Port: 3000,
}
// 这里设置接收者的IP地址为广播地址
raddr := net.UDPAddr{
    IP:   net.IPv4(255, 255, 255, 255),
    Port: 3000,
}
conn, err := net.DialUDP("udp", &laddr, &raddr)
if err != nil {
    println(err.Error())
    return
}
conn.Write([]byte(`hello peers`))
conn.Close()
*/

const (
	maxBufSize = 512
	mcAddress  = "XXX.XXX.XXX.XXX:XXXXX"
)

func newBroadcaster(address string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}

	return conn, nil

}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func listen(address string, handler func(*net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadBuffer(maxBufSize)

	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxBufSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		handler(src, numBytes, buffer)
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Println(n, "bytes read from", src)
	log.Println(hex.Dump(b[:n]))
}

func ping() {
	conn, _ := newBroadcaster(mcAddress)
	for {
		conn.Write([]byte("hello, world\n"))
		time.Sleep(1 * time.Second)
	}
}

//Init .
func Init() {
	go ping()
	go listen(mcAddress, msgHandler)
}
