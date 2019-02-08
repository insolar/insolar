import sys
import re


def parse_log_line(line):
    regex = 'scripts/insolard/discoverynodes/(\d+)/output.log.*\[ NET Consensus (\d+) phase-3 \] BFT consensus passed\: (\d+)/(\d+)'
    m = re.search(regex, line)
    if m:
        return int(m.group(1)), m.group(2), m.group(3)
    return 0, '', ''


data = {}
nodes = set()

for line in sys.stdin:
    node, pulse, active_nodes = parse_log_line(line)
    if node == 0 and pulse == '' and active_nodes == '':
        continue
    if pulse not in data:
        data[pulse] = {}
    data[pulse][node] = active_nodes
    nodes.add(node)

s = 'pulse\t'
for key in sorted(nodes):
    s += str(key) + '\t'

print(s, end='\n')

for key in sorted(data):
    s = str(key) + '\t'
    for inner_key in sorted(nodes):
        if inner_key in data[key]:
            s += data[key][inner_key] + '\t'
        else:
            s += '?\t'
    print(s, end='\n')
