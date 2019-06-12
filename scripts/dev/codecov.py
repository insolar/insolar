#!/usr/bin/env python3
#
# codecov.py is a try to implement https://codecov.io/ coverage algo for Go.
#
# Usage examples:
#
# ./scripts/dev/codecov.py [coverage.txt]
# grep 'FILTER' coverage.txt  | ./scripts/dev/codecov.py
#
# Official info about used algo:
# https://docs.codecov.io/docs/frequently-asked-questions#section-how-is-coverage-calculated-
# round((hits / (hits + partials + misses)) * 100, 5) = Coverage Ratio
#
# https://docs.codecov.io/docs/about-code-coverage

import argparse
import re
import sys

parser = argparse.ArgumentParser()
parser.add_argument('cover-file', nargs='?', default='coverage.txt', help="coverage file")
# https://docs.codecov.io/docs/fixing-reports
parser.add_argument("--fix", action="store_true", help="try to behave like fix")
args = parser.parse_args()

# stat per file with parsed data
stats = {}

def load_fixlines(fix_file_path):
    path_parts = fix_file_path.split('/')
    fix_file = '/'.join(path_parts[3:])
    # mark empty lines, comments, and lines with '}' only
    fixes = []
    with open(fix_file, 'r') as f:
        n = 0
        for line in f:
            n += 1
            line = line.rstrip("\n\r")
            analyze = line
            analyze = re.sub( r'^(.*)//.*', r'\1', analyze)
            analyze = re.sub( r'\s+', r'', analyze, 0, re.S)
            if (analyze == '}') or (len(analyze) == 0):
                fixes.append(n)
    return fixes


def code_line(offset_pair):
    return int(offset_pair.split('.', 2)[0])

cover_file = vars(args)['cover-file']

def append_pair_to_list(pairs_list, pair):
    for p in pairs_list:
        # pair starts on end of interval
        if pair[0] == p[1]:
            p[1] = pair[1]
            return
        # pair ends on start of interval
        if pair[1] == p[0]:
            p[0] = pair[0]
            return
    pairs_list.append(pair)

def parse_cover(f):
    for line in f:
        key, val = line.split(':', 2)
        if not key.endswith(".go"):
            continue

        cover_file = key
        st = stats.get(cover_file, None)
        if st is None:
            st = stats[cover_file] = {
                "hits": 0,
                "misses": 0,

                "fix_lines": load_fixlines(cover_file),
                "hits_borders": [],
                "miss_borders": [],
                "partials": [],
            }
# parse "27.47,28.10 1 1"
        cover_raw = val.split()
        line_offsets = cover_raw[0].split(',', 2)
        start, end = line_offsets[0], line_offsets[1]
# num_statements are used in std coverage tool
# https://github.com/golang/go/blob/50bd1c4d4eb4fac8ddeb5f063c099daccfb71b26/src/cmd/cover/func.go#L67
# https://github.com/golang/go/blob/50bd1c4d4eb4fac8ddeb5f063c099daccfb71b26/src/cmd/cover/func.go#L135:22
# but not in codecov
#        num_statements = int(cover_raw[1])
        not_covered = cover_raw[2] == '0'

        start_n, end_n = code_line(start), code_line(end)
        if not_covered:
            append_pair_to_list(st["miss_borders"], [start_n, end_n])
            continue
        append_pair_to_list(st["hits_borders"], [start_n, end_n])

if not sys.stdin.isatty():
    parse_cover(sys.stdin)
else:
    parse_cover(open(cover_file, 'r'))

def calc_partials(file_stat):
    for hit_pair in file_stat["hits_borders"]:
        for miss_pair in file_stat["miss_borders"]:
            partial = None
            # check only borders (can't overlap)
            if (hit_pair[0] in miss_pair):
                file_stat["partials"].append(hit_pair[0])
                file_stat["hits"] -= 1
                break
            if (hit_pair[1] in miss_pair):
                file_stat["partials"].append(hit_pair[1])
                file_stat["hits"] -= 1
                break

def sort_borders(borders):
    return sorted(borders, key=lambda x: x[0])

def apply_fixes(borders, fixes):
    result = borders
    i = 0
    applied_fixes = []
    for n in fixes:
        remove_pair = False
        splitted = []
        while i < len(result):
            pair = result[i]
            if n < pair[0]:
                break
            i += 1
            if n > pair[1]:
                # try next pair
                continue
            if n == pair[0] and n == pair[1]:
                remove_pair = True
                break
            if n-1 >= pair[0]:
                splitted.append([pair[0], n-1])
            if n+1<= pair[1]:
                splitted.append([n+1, pair[1]])
            break
        if len(splitted) > 0 or remove_pair:
            applied_fixes.append(n)
            result = result[:i-1] + splitted + result[i:]
            i -= 1
    return result

def recalc_hits_and_misses(file_stat):
    st = file_stat
    st["hits"], st["misses"] = 0, 0
    for pair in st["hits_borders"]:
        st["hits"] += pair[1] - pair[0] + 1
    for pair in st["miss_borders"]:
        st["misses"] += pair[1] - pair[0] + 1

# summary stats
t_hits, t_partials, t_misses = 0, 0, 0
components_stat = {}

# calc anf print per file stats
files = sorted(stats.keys())
max_file_length = len(max(files, key=lambda x: len(x)))
for file in files:
    st = stats[file]

    st["hits_borders"] = sort_borders(st["hits_borders"])
    st["miss_borders"] = sort_borders(st["miss_borders"])
    # DEBUG:
    # print("{}: before={}".format(file, st))
    if not args.fix:
        st["hits_borders"] = apply_fixes(st["hits_borders"], st["fix_lines"])
        st["miss_borders"] = apply_fixes(st["miss_borders"], st["fix_lines"])

    recalc_hits_and_misses(st)
    calc_partials(st)
    # DEBUG:
    # print("{}: {}".format(file, st))

    partials = len(st["partials"])
    hits = st["hits"]
    misses = st["misses"] - partials
    percents = hits/ (hits + partials + misses)

    fmt_string = "{:<" + str(max_file_length) + "}\t [{}, {}, {}]\t{:.2f}%"
    print(fmt_string.format(
        file,
        hits, partials, misses,
        round(percents, 5) * 100))

    # update global stat
    t_hits += hits
    t_partials += partials
    t_misses += misses

    # update per component stat
    path_parts = file.split('/')
    component = path_parts[3]
    comp_st = components_stat.get(component, None)
    if comp_st is None:
        comp_st = components_stat[component] = {"hits": 0, "partials": 0, "misses": 0}
    comp_st["hits"] += hits
    comp_st["partials"] += partials
    comp_st["misses"] += misses

# print per component stats
max_comp_length = len(max(components_stat.keys(), key=lambda x: len(x)))
print("-" * 75)
for comp in sorted(components_stat):
    st = components_stat[comp]
    perc = round(st["hits"]/(st["hits"]+st["partials"]+st["misses"]), 5) * 100
    line_stat = "[{}, {}, {}]".format(
        st["hits"], st["partials"], st["misses"],
    )
    fmt_string = "{:<" + str(max_comp_length) + "}\t{:<32} {:.02f}%"
    print(fmt_string.format(comp, line_stat, perc))

# print total stats
print("=" * 75)
print("Total:\t[{}, {}, {}] {:.2f}%".format(
    t_hits, t_partials, t_misses,
    round(t_hits/(t_hits+t_partials+t_misses), 5) * 100))
