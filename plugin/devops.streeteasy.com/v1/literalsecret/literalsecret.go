package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/api/kv"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	h *resmap.PluginHelpers

	types.ObjectMeta       `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	types.GeneratorOptions `json:"generatorOptions,omitempty" yaml:"generatorOptions,omitempty"`
	types.SecretArgs
}

var KustomizePlugin plugin

func (p *plugin) Config(
	h *resmap.PluginHelpers, c []byte) error {
	err := yaml.Unmarshal(c, p)
	if err != nil {
		return err
	}
	if p.SecretArgs.Name == "" {
		p.SecretArgs.Name = p.Name
	}
	if p.SecretArgs.Namespace == "" {
		p.SecretArgs.Namespace = p.Namespace
	}
	p.h = h
	return nil
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	ldr := kv.NewLoader(p.h.Loader(), p.h.Validator())

	all, err := ldr.Load(p.SecretArgs.KvPairSources)
	if err != nil {
		return nil, err
	}

	sources := []string{}
	for _, pair := range all {
		sources = append(sources, fmt.Sprintf("%s=%s", pair.Key, os.ExpandEnv(pair.Value)))
	}

	p.EnvSources = []string{}
	p.FileSources = []string{}
	p.LiteralSources = sources

	return p.h.ResmapFactory().FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), &p.GeneratorOptions, p.SecretArgs)
}
