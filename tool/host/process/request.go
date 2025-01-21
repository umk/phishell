package process

import (
	"encoding/json"
	"fmt"

	"github.com/umk/phishell/provider"
)

func (p *Process) requestSend(req *provider.ToolRequest) (chan any, error) {
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
		close(ch)
		delete(p.requests, req.CallID)

		return nil, err
	}

	return ch, nil
}

func (p *Process) requestResolve(res *provider.ToolResponse) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.requests[res.CallID]; ok {
		ch <- res
	}
}

func (p *Process) requestClose(callID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, ok := p.requests[callID]; ok {
		close(ch)
		delete(p.requests, callID)
	}
}
