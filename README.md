# nixbench

[![Build Status](https://travis-ci.org/jgillich/nixbench.svg?branch=master)](https://travis-ci.org/jgillich/nixbench)

A better benchmarking tool for servers.

To invoke, run:

```sh
curl -sS https://raw.githubusercontent.com/jgillich/nixbench/master/nixbench.sh | sh
```

Alternatively, you can download the latest release for your platform from the
[releases page](https://github.com/jgillich/nixbench/releases).

Sample output:

```
nixbench 0.1 - github.com/jgillich/nixbench

Host
----

OS        : linux
Platform  : ubuntu 16.04
Virt      :
CPU       : Intel(R) Xeon(R) CPU E5-26xx (Sandy Bridge)
Cores     : 2
Clock     : 2099 Mhz


Disk
----

1. run   : 291 MB/s
2. run   : 310 MB/s
3. run   : 459 MB/s
4. run   : 422 MB/s
5. run   : 351 MB/s
Average  : 367 MB/s


Geekbench
---------

Single-Core Score  : 2061
Multi-Core Score   : 3717
Result URL         : https://browser.geekbench.com/v4/cpu/3666474


Net
---

CDN                           : 477 Mbit/s
Amsterdam, The Netherlands    : 366 Mbit/s
Dallas, USA                   : 96 Mbit/s
Frankfurt, Germany            : 142 Mbit/s
Hong Kong, China              : 29 Mbit/s
London, United Kingdoms       : 397 Mbit/s
Melbourne, Australia          : 41 Mbit/s
Oslo, Norway                  : 237 Mbit/s
Paris, France                 : 427 Mbit/s
Querétaro Mexico              : 72 Mbit/s
San Jose, USA                 : 87 Mbit/s
Sao Paulo, Brazil             : 46 Mbit/s
Seoul, Korea                  : 30 Mbit/s
Singapore, Singapore          : 38 Mbit/s
Tokyo, Japan                  : 45 Mbit/s
Toronto, Canada               : 126 Mbit/s
Washington, USA               : 35 Mbit/s
```
