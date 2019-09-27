#!/usr/bin/env python3
"""
Changelog generator
"""

import argparse
import re
import subprocess
import sys

import jinja2
import requests
from logzero import logger


def main(args):
    issue_keys = get_issues_from_git(args.current_release, args.new_release)
    issues = get_issue_info(issue_keys, args.email, args.api_token)
    if len(issues) == 0:
        logger.warning('No issues in release')
        return

    if args.html:
        save_html(
            {
                'issues': issues,
                'tags': {
                    'current': args.current_release,
                    'new': args.new_release,
                },
            },
            args.html,
        )
        if sys.platform == 'darwin':
            subprocess.run(['open', args.html])
    else:
        print_changelog(args.current_release, args.new_release, issues)


# git log v0.9.0-rc3..v0.9.0-rc4 --oneline | grep -E '[A-Z]+-\d+' | perl -ne '/([A-Z]+-\d+)/; print "$1\n"' | sort -u
def get_issues_from_git(current_release, new_release):
    git_log = subprocess.run(
        ['git', 'log', f'{current_release}..{new_release}', '--oneline'],
        capture_output=True,
    )
    git_log = git_log.stdout.splitlines()
    issue_re = re.compile(r'[A-Z]+-\d+')
    issue_keys = []
    for msg in git_log:
        matches = issue_re.findall(msg.decode('utf-8'))
        for m in matches:
            issue_keys.append(m)
    return sorted(set(issue_keys))


def get_issue_info(issue_keys, email, api_token):
    issues = []
    for key in issue_keys:
        r = requests.get(
            f'https://insolar.atlassian.net/rest/api/3/issue/{key}',
            {
                'fields': 'summary,components,status,resolution,issuetype,priority,fixVersions',
            },
            auth=(email, api_token)
        )
        try:
            i = r.json()
            if 'errors' in i:
                logger.warning(f'Failed to get info for issue {key}')
            else:
                issues.append(i)
        except ValueError:
            logger.error(f'Failed to get info for issue {key}')
    return issues


def print_changelog(current, new, issues):
    print(f'Changelog from {current} to {new}:')
    for i in issues:
        print(f"{i['key']}: {i['fields']['summary']}")


def prepare_context(context):
    status_colors = {
        'new': 'secondary',
        'indeterminate': 'info',
        'done': 'success',
    }
    parsed_issues = []
    for i in context['issues']:
        try:
            key = i['fields']['status']['statusCategory']['key']
            if key in status_colors:
                i['fields']['status']['statusCategory']['status'] = status_colors[key]
            parsed_issues.append(i)
        except KeyError:
            logger.error(f'No fields in issue: {i}')
    context['issues'] = parsed_issues
    return context


def save_html(context, filename):
    context = prepare_context(context)
    html = render_html(context)
    with open(filename, 'w+') as f:
        f.write(html)


def render_html(context):
    tpl = jinja2.Environment().from_string(PAGE_TPL)
    return tpl.render(context)


PAGE_TPL = '''
<!doctype html>
<html lang="ru">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <link
        rel="stylesheet"
        href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css"
        integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T"
        crossorigin="anonymous"
    >

    <title>Release info</title>
</head>
<body>

<div class="container">
    <div class="row mt-5 mb-3">
        <div class="col-sm-12">
            <h3>Changelog from <b>{{ tags.current }}</b> to <b>{{ tags.new }}</b></h3>
        </div>
    </div>
    <div class="row">
        <div class="col-sm-12">
            <table class="table">
                <tbody>
                    {% for issue in issues %}
                        <tr>
                            <td>
                                <img src="{{ issue.fields.issuetype.iconUrl }}"
                                     title="{{ issue.fields.issuetype.name }}"
                                     width="16"
                                     height="16"
                                >
                            </td>
                            <td class="text-nowrap">
                                <a href="https://insolar.atlassian.net/browse/{{ issue.key }}"
                                   target="_blank"
                                >{{ issue.key }}</a>
                            </td>
                            <td>
                                <a href="https://insolar.atlassian.net/browse/{{ issue.key }}"
                                   target="_blank"
                                >{{ issue.fields.summary }}</a>
                            </td>
                            <td>
                                <img src="{{ issue.fields.priority.iconUrl }}"
                                     title="{{ issue.fields.priority.name }}"
                                     width="16"
                                     height="16"
                                >
                            </td>
                            <td class="h5">
                                <span
                                    class="badge badge-{{ issue.fields.status.statusCategory.status }}"
                                    title="{{ issue.fields.resolution.name }}"
                                    data-toggle="tooltip"
                                >{{ issue.fields.status.name }}</span>
                            </td>
                            <td class="text-nowrap">
                                {% for cmp in issue.fields.components %}
                                    {{ cmp.name }}{% if not loop.last %},{% endif %}
                                {% endfor %}
                            </td>
                        </tr>
                    {% endfor %}
                </tbody>
            </table>
        </div>
    </div>
</div>

<script
    src="https://code.jquery.com/jquery-3.3.1.slim.min.js"
    integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo"
    crossorigin="anonymous"
></script>
<script
    src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js"
    integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1"
    crossorigin="anonymous"
></script>
<script
    src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js"
    integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM"
    crossorigin="anonymous"
></script>
<script>
$(function () {
    $('[data-toggle="tooltip"]').tooltip()
})
</script>
</body>
</html>
'''

if __name__ == '__main__':
    ''' This is executed when run from the command line '''
    parser = argparse.ArgumentParser()

    parser.add_argument('current_release', help='Current release tag')
    parser.add_argument('new_release', help='New release tag')

    parser.add_argument('-e', '--email', help='JIRA account email', action='store', dest='email', required=True)
    parser.add_argument(
        '-t', '--api-token',
        help='JIRA API token. You can get one here: https://id.atlassian.com/manage/api-tokens',
        action='store', dest='api_token', required=True,
    )
    parser.add_argument('--html', help='Generate HTML file', action='store', dest='html', required=False)

    args = parser.parse_args()
    main(args)
