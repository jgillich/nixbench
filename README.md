# nixbench

[![Build Status](https://travis-ci.org/jgillich/nixbench.svg?branch=master)](https://travis-ci.org/jgillich/nixbench)

A better benchmarking tool for servers.

To invoke, run:

```sh
curl -sS https://raw.githubusercontent.com/jgillich/nixbench/master/nixbench.sh | sh
```

Alternatively, you can download the latest release for your platform from the
[releases page](https://github.com/jgillich/nixbench/releases).

### Supported platforms

nixbench has wide platform support: Linux, macOS, *BSD and more. If your
platform is missing, please open an issue.

### Sample output

```
nixbench v0.6 - https://github.com/jgillich/nixbench

cpu
---
Sha256 (1x) :  101.54 MB/s
Gzip (1x)   :   92.16 MB/s
AES (1x)    :  343.17 MB/s

Sha256 (32x) :  143.19 MB/s
Gzip (32x)   :  177.10 MB/s
AES (32x)    :  792.61 MB/s


disk
----
1. run   : 271 MB/s
2. run   : 248 MB/s
3. run   : 239 MB/s
4. run   : 188 MB/s
5. run   : 227 MB/s
Average  : 235 MB/s


geekbench
---------
Single-Core Score  : 1974
Multi-Core Score   : 6892
Result URL         : https://browser.geekbench.com/v4/cpu/3723526


host
----
OS        : linux
Platform  : ubuntu 14.04
CPU       : Intel(R) Xeon(R) CPU E5-2680 v2 @ 2.80GHz
Cores     : 32
Clock     : 2800 Mhz
RAM       : 60377 MB


net
---
CDN                           :  77.71 MB/s
Amsterdam, The Netherlands    :  12.80 MB/s
Dallas, USA                   :  15.10 MB/s
Frankfurt, Germany            :   9.09 MB/s
Hong Kong, China              :   6.96 MB/s
London, United Kingdoms       :  20.76 MB/s
Melbourne, Australia          :   5.40 MB/s
Oslo, Norway                  :   6.31 MB/s
Paris, France                 :   9.38 MB/s
Queretaro Mexico              :   7.63 MB/s
San Jose, USA                 :  13.20 MB/s
Sao Paulo, Brazil             :   6.94 MB/s
Seoul, Korea                  :   3.22 MB/s
Singapore, Singapore          :   6.67 MB/s
Tokyo, Japan                  :   3.22 MB/s
Toronto, Canada               :  27.34 MB/s
Washington, USA               :   6.75 MB/s
```
