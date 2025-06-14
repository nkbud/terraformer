// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rundeck

import (
	"errors"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/zclconf/go-cty/cty"
)

type RundeckProvider struct {
	terraformutils.Provider
	url         string
	token       string
	username    string
	password    string
	apiVersion  string
	insecureSSL bool
}

func (p *RundeckProvider) Init(args []string) error {
	p.url = args[0]
	p.token = args[1]
	p.username = args[2]
	p.password = args[3]
	p.apiVersion = args[4]

	if args[5] == "true" {
		p.insecureSSL = true
	} else {
		p.insecureSSL = false
	}

	return nil
}

func (p *RundeckProvider) GetName() string {
	return "rundeck"
}

func (p *RundeckProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{}
}

func (p *RundeckProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"url":         cty.StringVal(p.url),
		"token":       cty.StringVal(p.token),
		"username":    cty.StringVal(p.username),
		"password":    cty.StringVal(p.password),
		"api_version": cty.StringVal(p.apiVersion),
		"insecure":    cty.BoolVal(p.insecureSSL),
	})
}

func (p *RundeckProvider) GetBasicConfig() cty.Value {
	return p.GetConfig()
}

func (p *RundeckProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"url":         p.url,
		"token":       p.token,
		"username":    p.username,
		"password":    p.password,
		"api_version": p.apiVersion,
		"insecure":    p.insecureSSL,
	})
	return nil
}

func (p *RundeckProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		"jobs":     &JobGenerator{},
		"projects": &ProjectGenerator{},
	}
}

func (RundeckProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}
