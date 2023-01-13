# site
## The code that powers my personal site.

It is inspired by the indieweb philosophy of self-hosting and ownership of
data. See [contributing.md](./docs/CONTRIBUTING.md) for information on how to
setup the project.

## Overview

The initial version of the site has the following components.

- chasqui: A proxy whose main purpose is to upgrade other HTTP services w/o any
  downtime.
- migrate: A database migration system.
- blog: Where I write my thoughts.
- finsta: An invite-only image board.
- webhookd: updates and redeploys other components upon learning of a new
  version.

For now all the components are written in Go, but that may change in the future.
