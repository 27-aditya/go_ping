package main

import ( 
	"fmt"
	"os"
	"net"
	"time"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	fmt.Println("Hello World")
	website := os.Args[1]
	dest, time, err := Ping(website)
	fmt.Println(dest, time, err)
}

func Ping(addr string) (*net.IPAddr, time.Duration, error) {
	dest, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		return dest, 0, err
	}
	
	con, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		return dest, 0, err
	}

	defer con.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Seq: 1,
			Data: []byte("hello-ping"),
		},
	}

	icmpMessage, err := message.Marshal(nil)
	if err != nil {
		return dest, 0, err
	}

	start := time.Now()
	_, err = con.WriteTo(icmpMessage, dest)
	if err != nil {
		return dest, 0, err
	}

	reply := make([]byte, 1500)
	con.SetReadDeadline(time.Now().Add(10*time.Second))
	n, _, err := con.ReadFrom(reply)
	if err != nil {
		return dest, 0, err
	}

	duration := time.Since(start)
	ReplyMessage, err := icmp.ParseMessage(1, reply[:n])

	if ReplyMessage.Type != ipv4.ICMPTypeEchoReply {
		return dest, 0, err
	}
	return dest, duration, nil
}