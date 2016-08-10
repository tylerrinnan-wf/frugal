// Autogenerated by Frugal Compiler (1.13.0)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package base

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Workiva/frugal/lib/go"
)

// (needed to ensure safety because of naive import list construction.)
var _ = thrift.ZERO
var _ = fmt.Printf
var _ = bytes.Equal

type FBaseFoo interface {
	BasePing(ctx *frugal.FContext) (err error)
}

type FBaseFooClient struct {
	transport       frugal.FTransport
	protocolFactory *frugal.FProtocolFactory
	oprot           *frugal.FProtocol
	mu              sync.Mutex
	methods         map[string]*frugal.Method
}

func NewFBaseFooClient(t frugal.FTransport, p *frugal.FProtocolFactory, middleware ...frugal.ServiceMiddleware) *FBaseFooClient {
	t.SetRegistry(frugal.NewFClientRegistry())
	methods := make(map[string]*frugal.Method)
	client := &FBaseFooClient{
		transport:       t,
		protocolFactory: p,
		oprot:           p.GetProtocol(t),
		methods:         methods,
	}
	methods["basePing"] = frugal.NewMethod(client, client.basePing, "basePing", middleware)
	return client
}

// Do Not Use. To be called only by generated code.
func (f *FBaseFooClient) GetWriteMutex() *sync.Mutex {
	return &f.mu
}

func (f *FBaseFooClient) BasePing(ctx *frugal.FContext) (err error) {
	ret := f.methods["basePing"].Invoke([]interface{}{ctx})
	if len(ret) != 1 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 1", len(ret)))
	}
	if ret[0] != nil {
		err = ret[0].(error)
	}
	return err
}

func (f *FBaseFooClient) basePing(ctx *frugal.FContext) (err error) {
	errorC := make(chan error, 1)
	resultC := make(chan struct{}, 1)
	if err = f.transport.Register(ctx, f.recvBasePingHandler(ctx, resultC, errorC)); err != nil {
		return
	}
	defer f.transport.Unregister(ctx)
	f.GetWriteMutex().Lock()
	if err = f.oprot.WriteRequestHeader(ctx); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageBegin("basePing", thrift.CALL, 0); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	args := BaseFooBasePingArgs{}
	if err = args.Write(f.oprot); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.WriteMessageEnd(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	if err = f.oprot.Flush(); err != nil {
		f.GetWriteMutex().Unlock()
		return
	}
	f.GetWriteMutex().Unlock()

	select {
	case err = <-errorC:
	case <-resultC:
	case <-time.After(ctx.Timeout()):
		err = frugal.ErrTimeout
	case <-f.transport.Closed():
		err = frugal.ErrTransportClosed
	}
	return
}

func (f *FBaseFooClient) recvBasePingHandler(ctx *frugal.FContext, resultC chan<- struct{}, errorC chan<- error) frugal.FAsyncCallback {
	return func(tr thrift.TTransport) error {
		iprot := f.protocolFactory.GetProtocol(tr)
		if err := iprot.ReadResponseHeader(ctx); err != nil {
			errorC <- err
			return err
		}
		method, mTypeId, _, err := iprot.ReadMessageBegin()
		if err != nil {
			errorC <- err
			return err
		}
		if method != "basePing" {
			err = thrift.NewTApplicationException(thrift.WRONG_METHOD_NAME, "basePing failed: wrong method name")
			errorC <- err
			return err
		}
		if mTypeId == thrift.EXCEPTION {
			error0 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
			var error1 thrift.TApplicationException
			error1, err = error0.Read(iprot)
			if err != nil {
				errorC <- err
				return err
			}
			if err = iprot.ReadMessageEnd(); err != nil {
				errorC <- err
				return err
			}
			if error1.TypeId() == frugal.RESPONSE_TOO_LARGE {
				err = thrift.NewTTransportException(frugal.RESPONSE_TOO_LARGE, "response too large for transport")
				errorC <- err
				return nil
			}
			err = error1
			errorC <- err
			return err
		}
		if mTypeId != thrift.REPLY {
			err = thrift.NewTApplicationException(thrift.INVALID_MESSAGE_TYPE_EXCEPTION, "basePing failed: invalid message type")
			errorC <- err
			return err
		}
		result := BaseFooBasePingResult{}
		if err = result.Read(iprot); err != nil {
			errorC <- err
			return err
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			errorC <- err
			return err
		}
		resultC <- struct{}{}
		return nil
	}
}

type FBaseFooProcessor struct {
	*frugal.FBaseProcessor
}

func NewFBaseFooProcessor(handler FBaseFoo, middleware ...frugal.ServiceMiddleware) *FBaseFooProcessor {
	p := &FBaseFooProcessor{frugal.NewFBaseProcessor()}
	p.AddToProcessorMap("basePing", &basefooFBasePing{handler: frugal.NewMethod(handler, handler.BasePing, "BasePing", middleware), writeMu: p.GetWriteMutex()})
	return p
}

type basefooFBasePing struct {
	handler *frugal.Method
	writeMu *sync.Mutex
}

func (p *basefooFBasePing) Process(ctx *frugal.FContext, iprot, oprot *frugal.FProtocol) error {
	args := BaseFooBasePingArgs{}
	var err error
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		p.writeMu.Lock()
		basefooWriteApplicationError(ctx, oprot, thrift.PROTOCOL_ERROR, "basePing", err.Error())
		p.writeMu.Unlock()
		return err
	}

	iprot.ReadMessageEnd()
	result := BaseFooBasePingResult{}
	var err2 error
	ret := p.handler.Invoke([]interface{}{ctx})
	if len(ret) != 1 {
		panic(fmt.Sprintf("Middleware returned %d arguments, expected 1", len(ret)))
	}
	if ret[0] != nil {
		err2 = ret[0].(error)
	}
	if err2 != nil {
		p.writeMu.Lock()
		basefooWriteApplicationError(ctx, oprot, thrift.INTERNAL_ERROR, "basePing", "Internal error processing basePing: "+err2.Error())
		p.writeMu.Unlock()
		return err2
	}
	p.writeMu.Lock()
	defer p.writeMu.Unlock()
	if err2 = oprot.WriteResponseHeader(ctx); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			basefooWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "basePing", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageBegin("basePing", thrift.REPLY, 0); err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			basefooWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "basePing", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = result.Write(oprot); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			basefooWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "basePing", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.WriteMessageEnd(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			basefooWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "basePing", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	if err2 = oprot.Flush(); err == nil && err2 != nil {
		if frugal.IsErrTooLarge(err2) {
			basefooWriteApplicationError(ctx, oprot, frugal.RESPONSE_TOO_LARGE, "basePing", "response too large: "+err2.Error())
			return nil
		}
		err = err2
	}
	return err
}

func basefooWriteApplicationError(ctx *frugal.FContext, oprot *frugal.FProtocol, type_ int32, method, message string) {
	x := thrift.NewTApplicationException(type_, message)
	oprot.WriteResponseHeader(ctx)
	oprot.WriteMessageBegin(method, thrift.EXCEPTION, 0)
	x.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
}
