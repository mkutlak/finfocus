---
title: CI/CD Integration
description: Integrating FinFocus into CI/CD pipelines
---

You can run FinFocus in your CI/CD pipeline to enforce cost policies
and visibility.

## GitHub Actions

The [finfocus-action](https://github.com/rshade/finfocus-action) GitHub
Action installs FinFocus and its plugins, runs the analysis, and posts the
result as a pull request comment. It also supports analyzer mode, budgets,
recommendations, and sustainability metrics. To see the comment it produces,
open [the live demo pull request](https://github.com/rshade/finfocus-demo/pull/1)
in [finfocus-demo](https://github.com/rshade/finfocus-demo).

To run the CLI directly instead:

```yaml
name: Cost Estimate
on: [pull_request]

jobs:
  estimate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Install FinFocus
        run: |
          curl -fsSL https://raw.githubusercontent.com/rshade/finfocus/v0.3.1/scripts/install.sh \
            -o install.sh
          sha256sum --check <(echo "$(curl -fsSL https://raw.githubusercontent.com/rshade/finfocus/v0.3.1/scripts/install.sh.sha256)")
          sh install.sh

      - name: Pulumi Preview
        run: pulumi preview --json > plan.json
        env:
          PULUMI_ACCESS_TOKEN: ${{ secrets.PULUMI_ACCESS_TOKEN }}

      - name: Estimate Cost
        run: finfocus cost projected --pulumi-json plan.json
```

## GitLab CI

Run the preview and the estimate as two jobs. The preview job runs in the
image your Pulumi program already uses, because it needs the program's
language runtime and cloud credentials. The estimate job runs on the
published FinFocus image, which holds only the `finfocus` binary. Both
images set an entrypoint (`pulumi` and `finfocus`), so clear it to let the
GitLab Docker executor run the script.

```yaml
stages: [preview, cost]

pulumi-preview:
  stage: preview
  image:
    name: pulumi/pulumi:latest
    entrypoint: [""]
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
  script:
    - pulumi install
    - pulumi stack select "$PULUMI_STACK"
    - pulumi preview --json > plan.json
  artifacts:
    paths: [plan.json]
    expire_in: 1 day

finfocus-cost:
  stage: cost
  image:
    name: ghcr.io/rshade/finfocus:0.5.4
    entrypoint: [""]
  needs: [pulumi-preview]
  rules:
    - if: $CI_PIPELINE_SOURCE == "merge_request_event"
  script:
    - finfocus plugin install aws-public
    - finfocus cost projected --pulumi-json plan.json
```

- Image tags are release versions without the leading `v` (`0.5.4`, `0.5`,
  `latest`). Pin a version so a release does not change the job.
- The preview job needs `PULUMI_STACK` (the stack to preview),
  `PULUMI_ACCESS_TOKEN` (or your state backend's login), and the cloud
  credentials the program uses, set as CI/CD variables (mask the secrets).
  `pulumi install` installs the program's dependencies and provider plugins.
- `finfocus plugin install` downloads the plugin from its GitHub release on
  every run. Use the plugin your stack needs; `github.com/<owner>/<repo>`
  installs one that is not in the registry. The runner needs access to
  `ghcr.io` and `github.com`.
- FinFocus reads a resource's region from its properties, then from the
  cloud's environment variable (`AWS_REGION` or `AWS_DEFAULT_REGION` for
  AWS, `SCW_DEFAULT_ZONE` or `SCW_DEFAULT_REGION` for Scaleway). It does not
  read provider or stack configuration, so set the variable on the estimate
  job when your resources take their region from there.
- The estimate job has the checkout's `Pulumi.yaml` but no `pulumi`, so it
  prints `Warning: stack state unavailable` and prices the plan alone. That
  is expected.
- Add `--output json` to keep a machine-readable result as an artifact, or
  `--exit-on-threshold` to fail the job when a budget is exceeded. The
  threshold needs a budget, such as `cost.budgets.global.amount` in the
  project's `.finfocus/config.hujson`.
