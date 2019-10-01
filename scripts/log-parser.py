import typing
from typing import Collection, Mapping, Pattern, Match

import re
import sys
import time
import argparse
import fileinput
import collections

try:
    import orjson as json
except ImportError:
    sys.stderr.write('WARN: orjson library not found, json parsing is slower\n')
    sys.stderr.write('WARN: install orjson using: `pip3 install orjson` \n')

    import json as json

try:
    from termcolor import colored
except ImportError:
    sys.stderr.write('WARN: termcolor library not found, colors aren\'t available\n')
    sys.stderr.write('WARN: install termcolor using: `pip3 install termcolor` \n')

    def colored(string, *args, **kwargs):
        return string

try:
    from dateutil.parser import isoparse
except ImportError:
    sys.stderr.write('WARN: dateutil library not found, using standard parser\n')
    sys.stderr.write('WARN: install dateutil library using: `pip3 install python-dateutil`\n')

    import datetime

    parse_date = re.compile('^(-?(?:[1-9][0-9]*)?[0-9]{4})-(1[0-2]|0[1-9])-(3[01]|0[1-9]|[12][0-9])T(2[0-3]|[01][0-9]):([0-5][0-9]):([0-5][0-9])(\.[0-9]+)?(Z)?$')
    def isoparse(string):
        match = parse_date.match(string)

        if match is None:
            sys.stderr.write('WARN: failed to parse date "%s"\n/' % string)
            return None

        return datetime.datetime(
            year        = int(match.group(1)),
            month       = int(match.group(2)),
            day         = int(match.group(3)),
            hour        = int(match.group(4)),
            minute      = int(match.group(5)),
            second      = int(match.group(6)),
            microsecond = int(match.group(7)[1:7]),
        )

class FilterOptions:
    filter_caller:   typing.Collection[str]
    filter_loglevel: typing.Collection[str]
    filter_message:  typing.Collection[str]

    def __init__(self, args):
        self.filter_caller = [
            'network/',
        ]
        self.filter_caller.extend(args.skip_caller)

        self.filter_message  = args.filter_message
        self.filter_loglevel = args.filter_loglevel

    def match_message(self, line: Mapping[str, str]):
        message = line.get('message', None)
        if message is None:
            return True

        for msg in self.filter_message:
            if message.find(msg) != -1:
                return True

        return False

    def match_caller(self, line: Mapping[str, str]):
        caller = line.get('caller', None)
        if caller is None:
            return True

        for msg in self.filter_caller:
            if caller.find(msg) != -1:
                return False

        return True

    def filter_logline(self, line: Mapping[str, str]):
        if self.filter_loglevel and not (line['level'] in self.filter_loglevel):
            return True
        if len(self.filter_caller) != 0 and not self.match_caller(line):
            return True
        if len(self.filter_message) != 0 and not self.match_message(line):
            return True
        return False


fatal_prefix = 'FATAL: '
log_line_regexp: Collection[Pattern] = [
    re.compile('(\d+)/output.log:\d+:({.*})'),                             # dir-%d/output.log:line:<log>
    re.compile('(\d+)/output.log:({.*})'),                                 # dir-%d/output.log:line:<log>
    re.compile('[0-9\-TZ]+ .+-(\d+) [^ ]+: ({.*})'),                       # time host-%d ???: <log>
    re.compile('((?:virtual|heavy|light)?-?\d+)[^{]+({.*})'), # (virtual|heavy|light)-<number>.<id>.insolar.io <other>: <log>
    re.compile('(\d+).insolar.io[^{]+({.*})'),                             # dev-insolar-*-().vm.insolar.io
    re.compile('\d+/(\d+)/[^{]+({.*})'),                                   # 3165/32001/insolard_1569777807.log.gz:85947:{
]

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

def read_input(filter_options: FilterOptions) -> Collection[Mapping[str, str]]:
    lines: Collection[Mapping[str, str]] = []

    for line in sys.stdin.readlines():
        mo = deduce(line.strip())

        if mo != None:
            parsed_log_line = None
            try:
                parsed_log_line = json.loads(mo.group(2))
            except:
                print("failed to parse json: '%s'" % mo.group(2), file=sys.stderr)
                continue

            if not filter_options.filter_logline(parsed_log_line):
                time = isoparse(parsed_log_line.pop('time'))

                lines.append({
                    'instance': mo.group(1),
                    'time':     time,
                    'group':    parsed_log_line,
                })
        else:
            print("failed to match '%s'" % line)
            break

    return lines

color_flag: str = 'on'
stdout_is_term: bool = None
def wcolored(input, *args, **kwargs):
    global color_flag
    global stdout_is_term

    if stdout_is_term is None:
        stdout_is_term = sys.stdout.isatty()

    if (not stdout_is_term and not color_flag == 'force') or color_flag == 'off':
        return input
    else:
        return colored(input, *args, **kwargs)

prev_node = -1
node_pulses = collections.defaultdict((lambda: None))
def print_header(node: int, role: str, pulse: int):
    global prev_pulse
    global prev_node

    print_next_node = False
    print_next_pulse = False
    if prev_node != node:
        print_next_node = True

    if node_pulses[node] is not None and node_pulses[node] != pulse:
        print_next_pulse = True

    if print_next_node or print_next_pulse:
        sys.stdout.write('\n')
        sys.stdout.write('Node %s - %s - %s' % (node, role, pulse))
        if print_next_pulse:
            sys.stdout.write(' PULSE CHANGED')
        sys.stdout.write(' ===\n')

    prev_node = node
    node_pulses[node] = pulse

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

def print_output(lines: Collection[Mapping[str, str]], skip_field: Collection[str]):
    level_names = {
        'error': wcolored('ERR', 'red'),
        'info':  wcolored('INF', 'yellow', attrs=['bold']),
        'debug': wcolored('DBG', 'yellow'),
        'warn':  wcolored('WRN', 'magenta'),
    }

    prevPulse = -1
    prevNode = -1
    for line in lines:
        # format is like <path>/<instance>/output.log:<line>:<datetime> <LVL> <message>
        parsed_log_line: typing.Dict[str, str] = line['group']

        node      = line['instance']
        timestamp = line['time']

        level = parsed_log_line.pop('level')
        msg   = parsed_log_line.pop('message')
        role  = parsed_log_line.pop('role', None)
        pulse = parsed_log_line.pop('pulse', None)
        if pulse is None:
            pulse = parsed_log_line.pop('new_pulse', None)
        if pulse is not None:
            pulse = int(pulse)

        del parsed_log_line['writeDuration']
        del parsed_log_line['loginstance']
        del parsed_log_line['caller']

        print_header(node, role, pulse)

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
        action  = 'store',
        default = 'on',
        help    = 'enable colored output (default on, on|off|force)'
    )
    parser.add_argument(
        '--filter-loglevel',
        action  = 'append',
        default = [],
        help    = 'show only messages with given loglevel'
    )
    parser.add_argument(
        '--filter-message',
        action  = 'append',
        default = [],
        help    = 'show only messages that contains string'
    )
    parser.add_argument(
        '--debug',
        action  = 'store_true',
        default = False,
        help    = 'enable debug output'
    )
    args = parser.parse_args()
    if args.color not in ['on', 'off', 'force']:
        sys.exit(1)
    color_flag = args.color

    filter_message = FilterOptions(args)

    stime = time.perf_counter()
    lines = read_input(filter_message)
    if args.debug:
        print("read everything in %lf seconds (%d lines)" % (time.perf_counter() - stime, len(lines)), file=sys.stderr)

    stime = time.perf_counter()
    lines.sort(key=lambda val: val['time'])
    if args.debug:
        print("sorted everything in %lf seconds" % (time.perf_counter() - stime), file=sys.stderr)

    stime = time.perf_counter()
    print_output(lines, args.skip_field)
    if args.debug:
        print("printed everything in %lf seconds" % (time.perf_counter() - stime), file=sys.stderr)

    sys.exit(0)
