package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/envkey/envkey-fetch/fetch"
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

	SecretName string   `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	Prefix     string   `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Strip      bool     `json:"strip,omitempty" yaml:"strip,omitempty"`
	Templates  []string `json:"templates,omitempty" yaml:"templates,omitempty"`
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
	secret, err := getSecret(p.SecretName)
	if err != nil {
		return nil, err
	}

	newSecret := make(map[string]string)

	for k, v := range secret {
		if strings.HasPrefix(k, p.Prefix) {
			if p.Strip {
				k = strings.TrimPrefix(k, p.Prefix)
			}
			newSecret[k] = v
		}
	}

	for k, v := range newSecret {
		p.LiteralSources = append(p.LiteralSources, k+"="+v)
	}

	for _, t := range p.Templates {
		b, err := p.h.Loader().Load(t)
		if err != nil {
			return nil, err
		}

		tmpl, err := template.New(t).Parse(string(b))
		if err != nil {
			return nil, err
		}

		var buff bytes.Buffer
		err = tmpl.Execute(&buff, newSecret)
		if err != nil {
			return nil, err
		}
		p.LiteralSources = append(p.LiteralSources, filepath.Base(t)+"="+buff.String())
	}
	return p.h.ResmapFactory().FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), &p.GeneratorOptions, p.SecretArgs)

}

var secrets = make(map[string]map[string]string)

func getSecret(secretName string) (map[string]string, error) {
	if m, ok := secrets[secretName]; ok {
		return m, nil
	}

	envkey := os.Getenv(secretName)
	if envkey == "" {
		return nil, errors.New("missing ENVKEY")
	}

	res, err := fetch.Fetch(
		envkey,
		fetch.FetchOptions{false, "", "envkeygo", "", false, 15.0, 3, 1},
	)
	if err != nil {
		return nil, err
	}

	var m map[string]string
	err = json.Unmarshal([]byte(res), &m)

	secrets[secretName] = m

	return m, nil
}
