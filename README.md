<div align="center">

# WebBenchmark
*An http benchmark tool that can exhaust the bandwidth of your server.*

[![](https://img.shields.io/github/license/lvillis/webbenchmark?style=flat-square)](https://github.com/lvillis/webbenchmark)
[![](https://img.shields.io/github/repo-size/lvillis/webbenchmark?style=flat-square&color=328657)](https://github.com/lvillis/webbenchmark)
[![Github Actions](https://img.shields.io/github/workflow/status/lvillis/webbenchmark/Docker?style=flat-square)](https://github.com/lvillis/webbenchmark/actions)
[![](https://img.shields.io/github/last-commit/lvillis/webBenchmark?style=flat-square&label=commits)](https://github.com/lvillis/webbenchmark)
[![Docker Image Size (tag)](https://img.shields.io/docker/image-size/lvillis/webbenchmark/latest?style=flat-square)](https://hub.docker.com)
[![Docker Pulls](https://img.shields.io/docker/pulls/lvillis/webbenchmark?style=flat-square)](https://hub.docker.com)

</div>

---

## Features

* random User-Agent on every Request
* random x-forward-for and x-Real-ip on every Request
* customizable Referer Url
* cocurrent thread as you wish, depends on you server performance.
* post method.

## Usage
```
docker run -d \
    --name=webbenchmark \
    -e url="http://cachefly.cachefly.net/100mb.test" \
    -e method="GET" \
    -e thread=8 \
    --restart=always lvillis/webbenchmark:latest
```

## Related open source projects:

* webBenchmark：https://github.com/maintell/webBenchmark