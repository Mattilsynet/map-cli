#INFO: use command underneath:
#nats kv get wadm_state host_default --raw | python3 plot_hosts.py | graph-easy --from=dot

#!/usr/bin/env python3
import sys, json
from graphviz import Digraph

data = json.load(sys.stdin)
dot = Digraph(format='svg')
dot.attr('node', shape='box', style='filled', fillcolor='lightgrey')

def sanitize(s):
    return s.replace('-', '_').replace('.', '_').replace(' ', '_')

for host_id, host_info in data.items():
    host_label = host_info.get("friendly_name", host_id)
    host_node_id = sanitize(host_label)
    dot.node(host_node_id, label=host_label)

    for comp_name in host_info.get("components", {}):
        comp_node_id = sanitize(comp_name)
        dot.node(comp_node_id, label=f"Component: {comp_name}", shape='ellipse', fillcolor='lightblue')
        dot.edge(host_node_id, comp_node_id)

    for provider in host_info.get("providers", []):
        provider_id = provider.get("provider_id", "unknown")
        app = provider.get("annotations", {}).get("wasmcloud.dev/appspec", "unknown")
        provider_node_id = sanitize(provider_id)
        dot.node(provider_node_id, label=f"Provider: {provider_id}\n({app})", shape='ellipse', fillcolor='lightgreen')
        dot.edge(host_node_id, provider_node_id)

print(dot.source)
