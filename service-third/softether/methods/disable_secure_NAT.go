package methods

import (
	"bytes"
	"encoding/json"

	"github.com/mizuki1412/go-core-kit/v2/service-third/softether/pkg"
)

type DisableSecureNAT struct {
	pkg.Base
	Params *DisableSecureNATParams `json:"params"`
}

func (g *DisableSecureNAT) Parameter() pkg.Params {
	return g.Params
}

func NewDisableSecureNAT(hubName string) *DisableSecureNAT {
	return &DisableSecureNAT{
		Base:   pkg.NewBase("DisableSecureNAT"),
		Params: newDisableSecureNATParams(hubName),
	}
}

func (g *DisableSecureNAT) Name() string {
	return g.Base.Name
}

func (g *DisableSecureNAT) GetId() int {
	return g.Id
}

func (g *DisableSecureNAT) SetId(id int) {
	g.Base.Id = id
}

func (g *DisableSecureNAT) Marshall() ([]byte, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}
	res := bytes.Replace(data, []byte("null"), []byte("{}"), -1)
	return res, nil
}

type DisableSecureNATParams struct {
	HubName string `json:"HubName_str"`
}

func newDisableSecureNATParams(hubName string) *DisableSecureNATParams {
	return &DisableSecureNATParams{
		HubName: hubName,
	}
}

func (p *DisableSecureNATParams) Tags() []string {
	return []string{"HubName_str"}
}
