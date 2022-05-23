package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed template.go.tpl
var ginTemplate string

type serviceDesc struct {
	ServiceName     string
	ServiceFullName string
	Metadata        string
	Methods         []*methodDesc
	MethodSets      map[string]*methodDesc
}

type methodDesc struct {
	// method
	MethodName string
	Num        int
	Request    string
	Reply      string
	Comment    string

	// http_rule
	Path         string
	Method       string
	HasByte      bool
	HasFile		 bool
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.MethodName] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("gin").Parse(strings.TrimSpace(ginTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}

// initPathParams 转换参数路由 {xx} --> :xx
func (m *methodDesc) initPathParams() {
	paths := strings.Split(m.Path, "/")
	for i, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}') {
			paths[i] = ":" + p[1:len(p)-1]
		} else if len(p) > 0 && p[0] == ':' {
			paths[i] = p
		}
	}
	m.Path = strings.Join(paths, "/")
}

// HasPathParams 是否包含路由参数
func (m *methodDesc) HasPathParams() bool {
	paths := strings.Split(m.Path, "/")
	for _, p := range paths {
		if len(p) > 0 && (p[0] == '{' && p[len(p)-1] == '}' || p[0] == ':') {
			return true
		}
	}
	return false
}
