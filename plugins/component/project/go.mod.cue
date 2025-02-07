package main

import (
	"text/template"
	"cue.example/pkg"
)

module: template.Execute(file, Data)

file: """
	module {{.repository}}
	go 1.23.2
	
	require go.wasmcloud.dev/component v0.0.5
	
	require (
		github.com/Mattilsynet/mapis v0.0.2
		github.com/bytecodealliance/wasm-tools-go v0.3.1
		github.com/nats-io/nats.go v1.37.0
	)
	"""
