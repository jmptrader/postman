// Package smtp implements the Simple Mail Transfer Protocol as defined in RFC 5321.
// It also implements the following extensions:
//  8BITMIME  RFC 1652
//  STARTTLS  RFC 3207
// Additional extensions may be handled by clients.
package util

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

// A StmpClient represents a client connection to an SMTP server.
type StmpClient struct {
	// Text is the textproto.Conn used by the StmpClient. It is exported to allow for
	// clients to add extensions.
	Text *textproto.Conn
	// keep a reference to the connection so it can be used to create a TLS
	// connection later
	conn net.Conn
	// whether the StmpClient is using TLS
	tls        bool
	serverName string
	// map of supported extensions
	ext        map[string]string
	localName  string // the name to use in HELO/EHLO
	didHello   bool   // whether we've said HELO/EHLO
	helloError error  // the error from the hello
}

// Dial returns a new StmpClient connected to an SMTP server at addr.
// The addr must include a port number.
func Dial(addr string, local string) (*StmpClient, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	host, _, _ := net.SplitHostPort(addr)
	return NewStmpClient(conn, host, local)
}

// NewStmpClient returns a new StmpClient using an existing connection and host as a
// server name to be used when authenticating.
func NewStmpClient(conn net.Conn, host string, local string) (*StmpClient, error) {
	text := textproto.NewConn(conn)
	_, _, err := text.ReadResponse(220)
	if err != nil {
		text.Close()
		return nil, err
	}
	c := &StmpClient{Text: text, conn: conn, serverName: host, localName: local}
	return c, nil
}

// Close closes the connection.
func (c *StmpClient) Close() error {
	return c.Text.Close()
}

// hello runs a hello exchange if needed.
func (c *StmpClient) hello() error {
	if !c.didHello {
		c.didHello = true
		err := c.ehlo()
		if err != nil {
			c.helloError = c.helo()
		}
	}
	return c.helloError
}

// Hello sends a HELO or EHLO to the server as the given host name.
// Calling this method is only necessary if the client needs control
// over the host name used.  The client will introduce itself as "localhost"
// automatically otherwise.  If Hello is called, it must be called before
// any of the other methods.
func (c *StmpClient) Hello(localName string) error {
	if c.didHello {
		return errors.New("smtp: Hello called after other methods")
	}
	c.localName = localName
	return c.hello()
}

// cmd is a convenience function that sends a command and returns the response
func (c *StmpClient) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.Text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	c.Text.StartResponse(id)
	defer c.Text.EndResponse(id)
	code, msg, err := c.Text.ReadResponse(expectCode)
	return code, msg, err
}

// helo sends the HELO greeting to the server. It should be used only when the
// server does not support ehlo.
func (c *StmpClient) helo() error {
	c.ext = nil
	_, _, err := c.cmd(250, "HELO %s", c.localName)
	return err
}

// ehlo sends the EHLO (extended hello) greeting to the server. It
// should be the preferred greeting for servers that support it.
func (c *StmpClient) ehlo() error {
	_, msg, err := c.cmd(250, "EHLO %s", c.localName)
	if err != nil {
		return err
	}
	ext := make(map[string]string)
	extList := strings.Split(msg, "\n")
	if len(extList) > 1 {
		extList = extList[1:]
		for _, line := range extList {
			args := strings.SplitN(line, " ", 2)
			if len(args) > 1 {
				ext[args[0]] = args[1]
			} else {
				ext[args[0]] = ""
			}
		}
	}
	c.ext = ext
	return err
}

// StartTLS sends the STARTTLS command and encrypts all further communication.
// Only servers that advertise the STARTTLS extension support this function.
func (c *StmpClient) StartTLS(config *tls.Config) error {
	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(220, "STARTTLS")
	if err != nil {
		return err
	}
	c.conn = tls.Client(c.conn, config)
	c.Text = textproto.NewConn(c.conn)
	c.tls = true
	return c.ehlo()
}

// Verify checks the validity of an email address on the server.
// If Verify returns nil, the address is valid. A non-nil return
// does not necessarily indicate an invalid address. Many servers
// will not verify addresses for security reasons.
func (c *StmpClient) Verify(addr string) error {
	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(250, "VRFY %s", addr)
	return err
}

// Mail issues a MAIL command to the server using the provided email address.
// If the server supports the 8BITMIME extension, Mail adds the BODY=8BITMIME
// parameter.
// This initiates a mail transaction and is followed by one or more Rcpt calls.
func (c *StmpClient) Mail(from string) error {
	if err := c.hello(); err != nil {
		return err
	}
	cmdStr := "MAIL FROM:<%s>"
	if c.ext != nil {
		if _, ok := c.ext["8BITMIME"]; ok {
			cmdStr += " BODY=8BITMIME"
		}
	}
	_, _, err := c.cmd(250, cmdStr, from)
	return err
}

// Rcpt issues a RCPT command to the server using the provided email address.
// A call to Rcpt must be preceded by a call to Mail and may be followed by
// a Data call or another Rcpt call.
func (c *StmpClient) Rcpt(to string) error {
	_, _, err := c.cmd(25, "RCPT TO:<%s>", to)
	return err
}

type dataCloser struct {
	c *StmpClient
	io.WriteCloser
}

func (d *dataCloser) Close() error {
	d.WriteCloser.Close()
	_, _, err := d.c.Text.ReadResponse(250)
	return err
}

// Data issues a DATA command to the server and returns a writer that
// can be used to write the data. The caller should close the writer
// before calling any more methods on c.
// A call to Data must be preceded by one or more calls to Rcpt.
func (c *StmpClient) Data() (io.WriteCloser, error) {
	_, _, err := c.cmd(354, "DATA")
	if err != nil {
		return nil, err
	}
	return &dataCloser{c, c.Text.DotWriter()}, nil
}

func sendMail(addr string, local string, from string, to string, msg []byte, tlsForbiden bool) error {
	c, err := Dial(addr, local)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.hello(); err != nil {
		return err
	}
	if !tlsForbiden {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err = c.StartTLS(nil); err != nil {
				// resend mail with tls support forbidden
				return sendMail(addr, local, from, to, msg, true)
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	if err = c.Rcpt(to); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// random return send result
func debugSendMail(_ string, _ string, from string, to string) error {
	log.Printf("now is send mail from %s to %s", from, to)
	<-time.After(time.Second * 2)
	randNum := Rand(0, 100)
	if randNum < 90 {
		return nil
	}
	if randNum < 95 {
		return errors.New("no such user find")
	}
	return errors.New("ip block")
}

// SendMail connects to the server at addr, switches to TLS if
// and then sends an email from address from, to addresses to, with
// message msg.
func SendMail(from string, to string, msg []byte, hostname string) error {
	domain := strings.SplitN(to, "@", 2)[1]
	addr, err := MxRecord(domain)
	if err != nil {
		return err
	}
	if os.Getenv("POSTMAN_DEBUG_MODE") == "true" {
		return debugSendMail(addr+":25", hostname, from, to)
	}
	return sendMail(addr+":25", hostname, from, to, msg, false)
}

// Extension reports whether an extension is support by the server.
// The extension name is case-insensitive. If the extension is supported,
// Extension also returns a string that contains any parameters the
// server specifies for the extension.
func (c *StmpClient) Extension(ext string) (bool, string) {
	if err := c.hello(); err != nil {
		return false, ""
	}
	if c.ext == nil {
		return false, ""
	}
	ext = strings.ToUpper(ext)
	param, ok := c.ext[ext]
	return ok, param
}

// Reset sends the RSET command to the server, aborting the current mail
// transaction.
func (c *StmpClient) Reset() error {
	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(250, "RSET")
	return err
}

// Quit sends the QUIT command and closes the connection to the server.
func (c *StmpClient) Quit() error {
	if err := c.hello(); err != nil {
		return err
	}
	_, _, err := c.cmd(221, "QUIT")
	if err != nil {
		return err
	}
	return c.Text.Close()
}
