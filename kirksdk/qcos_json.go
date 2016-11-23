package kirksdk

import (
	"encoding/json"
)

func (p ServiceSpec) MarshalJSON() ([]byte, error) {
	type Alias ServiceSpec
	return json.Marshal(&struct {
		Command       *[]string           `json:"command,omitempty"`
		EntryPoint    *[]string           `json:"entryPoint,omitempty"`
		Envs          *[]string           `json:"envs,omitempty"`
		Hosts         *[]string           `json:"hosts,omitempty"`
		LogCollectors *[]LogCollectorSpec `json:"logCollectors,omitempty"`
		Confs         *[]ConfSpec         `json:"confs,omitempty"`
		GpuUUIDs      *[]string           `json:"gpuUUIDs,omitempty"`
		*Alias
	}{
		Command:       stringS2P(p.Command),
		EntryPoint:    stringS2P(p.EntryPoint),
		Envs:          stringS2P(p.Envs),
		Hosts:         stringS2P(p.Hosts),
		LogCollectors: logCollectorSpecS2P(p.LogCollectors),
		Confs:         ConfSpecS2P(p.Confs),
		GpuUUIDs:      stringS2P(p.GpuUUIDs),
		Alias:         (*Alias)(&p),
	})
}

func (p JobTaskSpec) MarshalJSON() ([]byte, error) {
	type Alias JobTaskSpec
	return json.Marshal(&struct {
		Command       *[]string           `json:"command,omitempty"`
		EntryPoint    *[]string           `json:"entryPoint,omitempty"`
		Envs          *[]string           `json:"envs,omitempty"`
		Hosts         *[]string           `json:"hosts,omitempty"`
		LogCollectors *[]LogCollectorSpec `json:"logCollectors,omitempty"`
		Confs         *[]ConfSpec         `json:"confs,omitempty"`
		*Alias
	}{
		Command:       stringS2P(p.Command),
		EntryPoint:    stringS2P(p.EntryPoint),
		Envs:          stringS2P(p.Envs),
		Hosts:         stringS2P(p.Hosts),
		LogCollectors: logCollectorSpecS2P(p.LogCollectors),
		Confs:         ConfSpecS2P(p.Confs),
		Alias:         (*Alias)(&p),
	})
}

func (p UpdateJobArgs) MarshalJSON() ([]byte, error) {
	type Alias UpdateJobArgs
	return json.Marshal(&struct {
		Metadata *[]string `json:"metadata,omitempty"`
		*Alias
	}{
		Metadata: stringS2P(p.Metadata),
		Alias:    (*Alias)(&p),
	})
}

func (p JobTaskSpecEx) MarshalJSON() ([]byte, error) {
	type Alias JobTaskSpecEx
	return json.Marshal(&struct {
		LogCollectors *[]LogCollectorSpec `json:"logCollectors,omitempty"`
		Confs         *[]ConfSpec         `json:"confs,omitempty"`
		Command       *[]string           `json:"command,omitempty"`
		EntryPoint    *[]string           `json:"entryPoint,omitempty"`
		Envs          *[]string           `json:"envs,omitempty"`
		Hosts         *[]string           `json:"hosts,omitempty"`
		Deps          *[]string           `json:"deps,omitempty"`
		*Alias
	}{
		LogCollectors: logCollectorSpecS2P(p.LogCollectors),
		Confs:         ConfSpecS2P(p.Confs),
		Command:       stringS2P(p.Command),
		EntryPoint:    stringS2P(p.EntryPoint),
		Envs:          stringS2P(p.Envs),
		Hosts:         stringS2P(p.Hosts),
		Deps:          stringS2P(p.Deps),
		Alias:         (*Alias)(&p),
	})
}

//-----------------------------------------------------------------

func stringS2P(s []string) *[]string {
	if s == nil {
		return nil
	}

	return &s
}

func logCollectorSpecS2P(s []LogCollectorSpec) *[]LogCollectorSpec {
	if s == nil {
		return nil
	}

	return &s
}

func confSpecS2P(s []ConfSpec) *[]ConfSpec {
	if s == nil {
		return nil
	}

	return &s
}
