package main

import (
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
	types.ConfigMapArgs
}

var KustomizePlugin plugin

func (p *plugin) Config(
	h *resmap.PluginHelpers, c []byte) error {
	err := yaml.Unmarshal(c, p)
	if err != nil {
		return err
	}
	if p.ConfigMapArgs.Name == "" {
		p.ConfigMapArgs.Name = p.Name
	}
	if p.ConfigMapArgs.Namespace == "" {
		p.ConfigMapArgs.Namespace = p.Namespace
	}
	p.h = h
	return nil
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	sources := []string{}
	for _, s := range p.LiteralSources {
		sources = append(sources, os.ExpandEnv(s))
	}
	p.LiteralSources = sources
	return p.h.ResmapFactory().FromConfigMapArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), &p.GeneratorOptions, p.ConfigMapArgs)
}
