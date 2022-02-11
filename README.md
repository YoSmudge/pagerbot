# PagerBot
Heavily inspired by https://github.com/YoSmudge/pagerbot
Forked but the implementations are different.
- Uses go-slack github.com/slack-go/slack
- Uses go-pagerduty github.com/PagerDuty/go-pagerduty
- Posting message to channel is removed (not yet implemented)
- Specification are the same as the original

Update your Slack user groups based on your PagerDuty Schedules.

PagerBot is a simple program to rotate oncall user group. Provided with your PagerDuty and Slack API credentials, and some simple
configuration, it will update the usergroups automatically.

# Installation

Install the dependencies with go module
```shell go get .```

Then build
```shell go build```

You should have a nice `pagerbot` binary ready to go. You can also download prebuild binaries from
the [releases](https://github.com/qoharu/pagerbot/releases) page.

# Configuration

A basic configuration file will look like

```yaml
api_keys:
  slack: "abcd123"
  pagerduty:
    org: "songkick"
    key: "qwerty567"

groups:
  - name: firefighter
    schedules:
      - PAAAAAA
      - PBBBBBB
  - name: fielder
    schedules:
      - PCCCCCC
```

The configuration should be fairly straightforward, under API keys provide your Slack and Pagerduty keys. This can be also provided using environment variable by providing `SLACK_TOKEN` and `PAGERDUTY_TOKEN` respectively.

Under groups configure the Slack groups you'd like to update. Schedules is a list of PagerDuty schedule IDs.

Once done, you can run PagerBot with `./pagerbot -c /path/to/config.yaml`

It's recommended to run PagerBot under Upstart or some other process manager.

N.B. PagerBot matches PagerDuty users to Slack users by their email addresses, so your users must have the same email address in Slack as in PagerDuty. PagerBot will log warnings for any users it finds in PagerDuty but not in Slack.

# Take the Benefit of Github Actions
In my company, we use Github Actions to run PagerBot on Daily basis by using a scheduled workflow.

```yaml
name: Slack Pagerduty Oncall Rotation Scheduler

on:
  schedule:
  - cron: 0 4 * * *
  workflow_dispatch:

env:
  PAGERBOT_URL: https://github.com/qoharu/pagerbot/releases/download/v2.0.0/pagerbot-v2.0.0-linux-amd64.tar.gz
  CONFIG_FILE: ./pagerduty_oncall_rotation_sch.yaml
  SLACK_TOKEN: ${{ secrets.SLACK_TOKEN }}
  PAGERDUTY_TOKEN: ${{ secrets.PAGERDUTY_TOKEN }}
jobs:
  scheduler:
    name: Rotate Oncall
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - run: wget -c $PAGERBOT_URL -O - | tar -xz
    - run: ./pagerbot -c $CONFIG_FILE -v

```
