package persona

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/umk/phishell/bootstrap"
	"github.com/umk/phishell/tool/host/provider"
	"github.com/umk/phishell/tool/host/provider/schema"
	"github.com/umk/phishell/util/termx"
)

type Persona struct {
	profile  *bootstrap.Profile
	provider *provider.Provider
	server   *http.Server
}

func New(profile *bootstrap.Profile) *Persona {
	return &Persona{profile: profile}
}

func (p *Persona) Post(ctx context.Context, content string) error {
	switch p.provider.Info.Status {
	case provider.PsFailed:
		return fmt.Errorf("provider has failed: %s", p.provider.Info.StatusMessage)
	case provider.PsCompleted:
		return errors.New("provider has been completed")
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get working directory: %w", err)
	}

	arguments, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("cannot marshal content: %w", err)
	}

	r, err := p.provider.Post(&schema.Request{
		CallID: uuid.NewString(),
		Function: schema.Function{
			Name:      "Post",
			Arguments: string(arguments),
		},
		Context: schema.Context{
			Dir: wd,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to post request: %w", err)
	}

	termx.NewPrinter().Print(getContent(r))

	return nil
}

func getContent(res *schema.Response) (content string) {
	if c, ok := res.Content.(string); ok {
		content = c
	}

	if content == "" {
		content = "The request was completed, but provider did not return a meaningful response."
	}

	return
}
