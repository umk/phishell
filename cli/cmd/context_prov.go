package cmd

import (
	"github.com/umk/phishell/tool/host"
	"github.com/umk/phishell/util/execx"
)

type providers []*providerRef

type providerRef struct {
	args execx.Arguments

	process *host.Provider
	info    *host.ProviderInfo
}

func (p *providers) refresh() {
	current := make([]*providerRef, 0, len(*p))
	for _, bj := range *p {
		if bj.info.Status == host.PsRunning {
			current = append(current, bj)
		}
	}

	*p = current
}
