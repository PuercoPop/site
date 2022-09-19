# Overview

The wiki lives under the /p/ path.
Image assets are under the /i/ path.
video assets under the /v/ path.
Each tag has its own index under the /t/:name path.
There are other hard-coded routes like /about/, /atom.xml and /archive/.

# Ideas/leads

Checkout zerok/webmentiond.
github.com/ichiban/prolog
https://github.com/shurcooL/home
I like write.as approach of a subdomain for each feature
- write.puercopop.com: blogging
- tweet.puercopop.com: microblogging
- snap.puercopop.com: photos
- read.puercopop.com: RSS Reader

## Process Supervision

We want to have the application supervised for two reasons.

We want to supervise the process to make sure it is running. If it exits with a
non zero code we restart it. If we receive sighup we reload.

See littleboss for an example of how to do this in the same executable.

We also want the supervisor to hold the sockets open while we reload the
application.

The related projects section of socket masters' documentation

https://zimbatm.github.io/socketmaster/#RELATED-PROJECTS
