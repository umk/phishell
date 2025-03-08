package cmd

import (
	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/util/execx"
)

type providers []*providerRef

type providerRef struct {
	args execx.Arguments

	process *provider.Provider
	info    *provider.Info
}

func (p *providers) refresh() {
	current := make([]*providerRef, 0, len(*p))
	for _, bj := range *p {
		if bj.info.Status == provider.PsRunning {
			current = append(current, bj)
		}
	}

	*p = current
}
