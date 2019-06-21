package espat

import (
	"strconv"
	"strings"
	"time"
)

// DialUDP makes a UDP network connection. raadr is the port that the messages will
// be sent to, and laddr is the port that will be listened to in order to
// receive incoming messages.
func (d *Device) DialUDP(network string, laddr, raddr *UDPAddr) (*UDPSerialConn, error) {
	addr := raddr.IP.String()
	sendport := strconv.Itoa(raddr.Port)
	listenport := strconv.Itoa(laddr.Port)

	// disconnect any old socket
	d.DisconnectSocket()

	// connect new socket
	d.ConnectUDPSocket(addr, sendport, listenport)

	return &UDPSerialConn{SerialConn: SerialConn{Adaptor: d}, laddr: laddr, raddr: raddr}, nil
}

// ListenUDP listens for UDP connections on the port listed in laddr.
func (d *Device) ListenUDP(network string, laddr *UDPAddr) (*UDPSerialConn, error) {
	addr := "0"
	sendport := "0"
	listenport := strconv.Itoa(laddr.Port)

	// disconnect any old socket
	d.DisconnectSocket()

	// connect new socket
	d.ConnectUDPSocket(addr, sendport, listenport)

	return &UDPSerialConn{SerialConn: SerialConn{Adaptor: d}, laddr: laddr}, nil
}

// DialTCP makes a TCP network connection. raadr is the port that the messages will
// be sent to, and laddr is the port that will be listened to in order to
// receive incoming messages.
func (d *Device) DialTCP(network string, laddr, raddr *TCPAddr) (*TCPSerialConn, error) {
	addr := raddr.IP.String()
	sendport := strconv.Itoa(raddr.Port)

	// disconnect any old socket
	d.DisconnectSocket()

	// connect new socket
	d.ConnectTCPSocket(addr, sendport)

	return &TCPSerialConn{SerialConn: SerialConn{Adaptor: d}, laddr: laddr, raddr: raddr}, nil
}

// SerialConn is a loosely net.Conn compatible implementation
type SerialConn struct {
	Adaptor *Device
}

// UDPSerialConn is a loosely net.Conn compatible intended to support
// UDP over serial.
type UDPSerialConn struct {
	SerialConn
	laddr *UDPAddr
	raddr *UDPAddr
}

// TCPSerialConn is a loosely net.Conn compatible intended to support
// TCP over serial.
type TCPSerialConn struct {
	SerialConn
	laddr *TCPAddr
	raddr *TCPAddr
}

// Read reads data from the connection.
// TODO: implement the full method functionality:
// Read can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetReadDeadline.
func (c *SerialConn) Read(b []byte) (n int, err error) {
	// read only the data that has been received via "+IPD" socket
	return c.Adaptor.ReadSocket(b)
}

// Write writes data to the connection.
// TODO: implement the full method functionality for timeouts.
// Write can be made to time out and return an Error with Timeout() == true
// after a fixed time limit; see SetDeadline and SetWriteDeadline.
func (c *SerialConn) Write(b []byte) (n int, err error) {
	// specify that is a data transfer to the
	// currently open socket, not commands to the ESP8266/ESP32.
	c.Adaptor.StartSocketSend(len(b))
	return c.Adaptor.Write(b)
}

// Close closes the connection.
// Currently only supports a single Read or Write operations without blocking.
func (c *SerialConn) Close() error {
	c.Adaptor.DisconnectSocket()
	return nil
}

// LocalAddr returns the local network address.
func (c *UDPSerialConn) LocalAddr() UDPAddr {
	return *c.laddr
}

// RemoteAddr returns the remote network address.
func (c *UDPSerialConn) RemoteAddr() UDPAddr {
	return *c.laddr
}

// LocalAddr returns the local network address.
func (c *TCPSerialConn) LocalAddr() TCPAddr {
	return *c.laddr
}

// RemoteAddr returns the remote network address.
func (c *TCPSerialConn) RemoteAddr() TCPAddr {
	return *c.laddr
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
//
// A deadline is an absolute time after which I/O operations
// fail with a timeout (see type Error) instead of
// blocking. The deadline applies to all future and pending
// I/O, not just the immediately following call to Read or
// Write. After a deadline has been exceeded, the connection
// can be refreshed by setting a deadline in the future.
//
// An idle timeout can be implemented by repeatedly extending
// the deadline after successful Read or Write calls.
//
// A zero value for t means I/O operations will not time out.
func (c *SerialConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *SerialConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *SerialConn) SetWriteDeadline(t time.Time) error {
	return nil
}

// ResolveTCPAddr returns an address of TCP end point.
//
// The network must be a TCP network name.
//
func (d *Device) ResolveTCPAddr(network, address string) (*TCPAddr, error) {
	// TODO: make sure network is 'tcp'
	// separate domain from port, if any
	r := strings.Split(address, ":")
	ip, err := d.GetDNS(r[0])
	if err != nil {
		return nil, err
	}
	if len(r) > 1 {
		port, e := strconv.Atoi(r[1])
		if e != nil {
			return nil, e
		}
		return &TCPAddr{IP: ip, Port: port}, nil
	}
	return &TCPAddr{IP: ip}, nil
}

// ResolveUDPAddr returns an address of UDP end point.
//
// The network must be a UDP network name.
//
func (d *Device) ResolveUDPAddr(network, address string) (*UDPAddr, error) {
	// TODO: make sure network is 'udp'
	// separate domain from port, if any
	r := strings.Split(address, ":")
	ip, err := d.GetDNS(r[0])
	if err != nil {
		return nil, err
	}
	if len(r) > 1 {
		port, e := strconv.Atoi(r[1])
		if e != nil {
			return nil, e
		}
		return &UDPAddr{IP: ip, Port: port}, nil
	}

	return &UDPAddr{IP: ip}, nil
}

// The following definitions are here to support a Golang standard package
// net-compatible interface for IP until TinyGo can compile the net package.

// IP is an IP address. Unlike the standard implementation, it is only
// a buffer of bytes that contains the string form of the IP address, not the
// full byte format used by the Go standard .
type IP []byte

// UDPAddr here to serve as compatible type. until TinyGo can compile the net package.
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone; added in Go 1.1
}

// TCPAddr here to serve as compatible type. until TinyGo can compile the net package.
type TCPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone
}

// ParseIP parses s as an IP address, returning the result.
func ParseIP(s string) IP {
	return IP([]byte(s))
}

// String returns the string form of the IP address ip.
func (ip IP) String() string {
	return string(ip)
}
