package main

import (
	"text/template"
	"Data"
)

module: template.Execute(file, Data)

"""
		name = "{{ Data.repository }}"
		language = "tinygo"
		type = "component"
		version = "0.1.0"
		
		[component]
		wit_world = "{{ data.repository }}"
		wasm_target = "wasm32-wasi-preview2"
		destination = "build/{{ Data.repository }}_s.wasm"
	"""
