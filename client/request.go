package client

import (
	"fmt"
	"github.com/qmsk/snmpbot/snmp"
	"time"
)

/*
	PDU ::= SEQUENCE {
        request-id INTEGER (-214783648..214783647),
*/
type requestID int32

const requestIDWrapping = (1 << 31)

type requestMap map[ioKey]*Request

// send.Addr must be resolved using engine.Transport().Resolve(...)
func NewRequest(options Options, send IO) *Request {
	var request = makeRequest()

	request.send = send
	request.timeout = DefaultTimeout
	request.retry = DefaultRetry
	request.startTime = time.Now()

	if options.Timeout != 0 {
		request.timeout = options.Timeout
	}
	if options.Retry != 0 {
		request.retry = options.Retry
	}

	return &request
}

func makeRequest() Request {
	return Request{
		waitChan: make(chan error, 1),
	}
}

type Request struct {
	send      IO
	id        requestID
	timeout   time.Duration
	retry     uint
	startTime time.Time
	timer     *time.Timer
	waitChan  chan error
	recv      IO
	recvOK    bool
}

func (request Request) String() string {
	if request.recvOK {
		return fmt.Sprintf("%v<%v>@%v[%d] => %v", request.send.PDUType, request.send.PDU, request.send.Addr, request.id, request.recv.PDU)
	} else {
		return fmt.Sprintf("%v<%v>@%v[%d]", request.send.PDUType, request.send.PDU, request.send.Addr, request.id)
	}
}

// Return any SNMPError, or nil
func (request *Request) error() error {
	if pduError := request.recv.PDU.GetError(); pduError.ErrorStatus != 0 {
		return SNMPError{
			RequestType:   request.send.PDUType,
			ResponseType:  request.recv.PDUType,
			ResponseError: pduError,
		}
	} else {
		return nil
	}
}

func (request *Request) Result() (IO, error) {
	if !request.recvOK {
		return request.recv, fmt.Errorf("Request is not done")
	} else if err := request.error(); err != nil {
		return request.recv, err
	} else {
		return request.recv, nil
	}
}

func (request *Request) wait() error {
	if err, ok := <-request.waitChan; !ok {
		return fmt.Errorf("request canceled")
	} else {
		return err
	}
}

func (request *Request) init(id requestID) ioKey {
	request.id = id
	request.send.RequestID = int(id)

	return request.send.key()
}

func (request *Request) startTimeout(timeoutChan chan ioKey, key ioKey) {
	request.timer = time.AfterFunc(request.timeout, func() {
		timeoutChan <- key
	})
}

func (request *Request) close() {
	if request.timer != nil {
		request.timer.Stop()
	}
	close(request.waitChan)
}

func (request *Request) fail(err error) {
	request.waitChan <- err
	request.close()
}

func (request *Request) done(recv IO) {
	request.recv = recv
	request.recvOK = true
	request.waitChan <- nil
	request.close()
}

func (request *Request) failTimeout(transport Transport) {
	request.fail(TimeoutError{
		transport: transport,
		request:   request,
		Duration:  time.Now().Sub(request.startTime),
	})
}

type TimeoutError struct {
	transport Transport
	request   *Request
	Duration  time.Duration
}

func (err TimeoutError) Error() string {
	return fmt.Sprintf("SNMP<%v> timeout for %v after %v", err.transport, err.request, err.Duration)
}

type SNMPError struct {
	RequestType   snmp.PDUType
	ResponseType  snmp.PDUType
	ResponseError snmp.PDUError
}

func (err SNMPError) Error() string {
	return fmt.Sprintf("SNMP %v error: %v @ %v", err.RequestType, err.ResponseError.ErrorStatus, err.ResponseError.VarBind)
}
