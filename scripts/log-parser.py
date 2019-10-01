from __future__ import print_function

from dateutil.parser import isoparse
from termcolor import colored

import typing
from typing import Collection, Mapping, Pattern, Match

import re
import sys
import json
import time
import argparse
import fileinput

fatal_prefix = 'FATAL: '
log_line_regexp: Collection[Pattern] = [
    re.compile('(\d+)/output.log:\d+:(.*)'),       # dir-%d/output.log:line:<log>
    re.compile('(\d+)/output.log:(.*)'),           # dir-%d/output.log:line:<log>
    re.compile('[0-9\-TZ]+ .+-(\d+) [^ ]+: (.*)'), # time host-%d ???: <log>
    re.compile('(\d+).vm.insolar.io[^{]+(.*)'), # dev-insolar-*-().vm.insolar.io
    re.compile('\d+/(\d+)/[^{]+(.*)'), # 3165/32001/insolard_1569777807.log.gz:85947:{
]

force_color: bool = None

typing.Dict[str, str]

stdout_is_term: bool = None
def wcolored(input, *args, **kwargs):
    global force_color
    global stdout_is_term

    if stdout_is_term is None:
        stdout_is_term = sys.stdout.isatty()

    if not stdout_is_term and not force_color:
        return input
    else:
        return colored(input, *args, **kwargs)

parse_line: Pattern = None
def deduce(line: str) -> Match:
    global parse_line

    if parse_line is not None:
        match = parse_line.search(line)
        if match != None:
            return match

    for regexp in log_line_regexp:
        match = regexp.search(line)
        if match != None:
            parse_line = regexp
            return match

    return None

def is_ignored(parsed_line: Mapping[str,str], skip_caller: Collection[str]) -> bool:
    try:
        caller = parsed_line['caller']
    except KeyError:
        return True
    if caller.startswith('insolar/bus/bus.go'):
        return True
    if caller.startswith('network/'):
        return True

    return False


def construct_fields(parsed_log_line, skip_fields):
    fields, first_field = [], True

    fields.append(wcolored('[', 'green'))
    for fkey, fvalue in parsed_log_line.items():
        if fkey in skip_fields:
            continue

        if not first_field:
            fields.append(', ')
        else:
            first_field = False

        fields.append(wcolored(str(fkey), 'green'))
        fields.append('=')
        fields.append(repr(fvalue))

    if first_field == True:
        return ''

    fields.append(wcolored(']', 'green'))

    return ''.join(fields)

def read_input(skip_caller: Collection[str]) -> Collection[Mapping[str, str]]:
    lines: Collection[Mapping[str, str]] = []

    for line in sys.stdin.readlines():
        mo = deduce(line)

        if mo != None:
            parsed_log_line = None
            try:
                parsed_log_line = json.loads(mo.group(2))
            except:
                print("failed to parse json: ''%s''" % mo.group(2), file=sys.stderr)
                continue

            if not is_ignored(parsed_log_line, skip_caller):
                time = isoparse(parsed_log_line.pop('time'))

                lines.append({
                    'instance': int(mo.group(1)),
                    'time':     time,
                    'group':    parsed_log_line,
                })
        else:
            print("failed to match ''%s''" % line)
            continue

    return lines

def print_output(lines: Collection[Mapping[str, str]], skip_field: Collection[str]):
    level_names = {
        'error': wcolored('ERR', 'red'),
        'info':  wcolored('INF', 'yellow', attrs=['bold']),
        'debug': wcolored('DBG', 'yellow'),
        'warn':  wcolored('WRN', 'magenta'),
    }

    prevNode = -1
    for line in lines:
        # format is like <path>/<instance>/output.log:<line>:<datetime> <LVL> <message>
        parsed_log_line: typing.Dict[str, str] = line['group']

        node      = line['instance']
        timestamp = line['time']

        level = parsed_log_line.pop('level')
        msg   = parsed_log_line.pop('message')
        role  = parsed_log_line.pop('role', None)

        del parsed_log_line['writeDuration']
        del parsed_log_line['loginstance']
        del parsed_log_line['caller']

        if node != prevNode:
            if prevNode != -1:
                sys.stdout.write('\n')
            sys.stdout.write('Node %02d - %s ===\n' % (node, role))
            prevNode = node

        if msg.startswith(fatal_prefix):
            msg = msg[len(fatal_prefix):]
            level = wcolored('FTL', 'black', 'on_red')
        else:
            level = level_names[level]

        # [time, level] - message [field1, field2]
        sys.stdout.write('[')
        sys.stdout.write(timestamp.strftime('%H:%M:%S.%f'))
        sys.stdout.write(', ')
        sys.stdout.write(level)
        sys.stdout.write('] - ')
        sys.stdout.write(wcolored(msg, 'cyan'))
        sys.stdout.write(' ')
        sys.stdout.write(construct_fields(parsed_log_line, skip_field))
        sys.stdout.write('\n')


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Parse JSON logs of insolar')
    parser.add_argument(
        '--skip-field',
        dest    = 'skip_field',
        action  = 'append',
        default = [],
        help    = 'fields to skip'
    )
    parser.add_argument(
        '--skip-caller',
        dest    = 'skip_caller',
        action  = 'append',
        default = [],
        help    = 'caller to skip'
    )
    parser.add_argument(
        '--color',
        action  = 'store_true',
        default = False,
        help    = 'force wcolored output'
    )
    args = parser.parse_args()

    force_color = args.color

    stime = time.time()
    lines = read_input(args.skip_caller)
    print("read everything in %d seconds (%d lines)" % (time.time() - stime, len(lines)), file=sys.stderr)

    stime = time.time()
    lines.sort(key=lambda val: val['time'])
    print("sorted everything in %d seconds" % (time.time() - stime), file=sys.stderr)

    print_output(lines, args.skip_field)

    sys.exit(0)

