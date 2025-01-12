package host

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/umk/phishell/provider"
)

func (p *ToolProcess) Post(req *provider.ToolRequest) (*provider.ToolResponse, error) {
	ch, err := p.sendRequest(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		p.mu.Lock()
		defer p.mu.Unlock()

		delete(p.requests, req.CallID)

		close(ch)
	}()

	timeout := 15 * time.Second

	select {
	case res := <-ch:
		switch v := res.(type) {
		case *provider.ToolResponse:
			return v, nil
		case error:
			return nil, v
		default:
			panic("bad message type")
		}
	case <-time.After(timeout):
		return nil, fmt.Errorf("handler timed out after %s", timeout)
	}
}

func (p *ToolProcess) sendRequest(req *provider.ToolRequest) (chan any, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.requests[req.CallID]; ok {
		return nil, fmt.Errorf("request %s is already sent", req.CallID)
	}

	ch := make(chan any, 1)
	p.requests[req.CallID] = ch

	if _, err := fmt.Fprintln(p.stdin, string(content)); err != nil {
		return nil, err
	}

	return ch, nil
}

func (p *ToolProcess) resolveRequest(res *provider.ToolResponse) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.requests[res.CallID]; ok {
		ch <- res
		delete(p.requests, res.CallID)
	}
}
