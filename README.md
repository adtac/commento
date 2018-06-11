<p align="center">
<a href="https://commento.io"><img src="https://user-images.githubusercontent.com/7521600/33375172-14b21f68-d52f-11e7-9b30-477682bccf8f.png" width=300></a>
</p>

<p align="center"><b>An open source, privacy focused discussion platform.</b></p>

<p align="center">
<a href="https://commento.adtac.pw"><img src="https://img.shields.io/badge/demo-live-red.svg?style=for-the-badge&colorA=1e2127&colorB=98c379&label=demo"></a>
<a href="https://irc.commento.io"><img src="https://img.shields.io/badge/irc-%23commento-red.svg?style=for-the-badge&colorA=1e2127&colorB=c678dd&label=freenode"></a>
<a href="https://hub.docker.com/r/adtac/commento/"><img src="https://img.shields.io/badge/live-commento-red.svg?style=for-the-badge&colorB=e24329&colorA=554488&label=gitlab"></a>
</p>

<h2 align="center"></h2>

<p align="center">
:warning: :warning: :warning: &nbsp; <b>Note: All development has <a href="https://gitlab.com/commento/commento-ce">moved to GitLab</a>.</b> &nbsp; :warning: :warning: :warning:

You're currently viewing the deprecated `v0` branch. Version `0.x` will not receive new features any new features (except critical security patches). <a href="https://gitlab.com/commento/commento-ce">Version `1.x`</a> is the future with lots of new and shiny features (like moderation, voting). Please file new issues <a href="https://gitlab.com/commento/commento-ce">over there</a>.
</p>

<h2 align="center"></h2>

Commento is a discussion platform that you can self-host. You can embed it on your blog, news articles, and any place where you want your readers to add comments. Unlike most alternatives, Commento is lightweight and privacy focused; we'll never sell your data, show ads, embed third-party tracking scripts, or inject affiliate links.

### Features

 - Privacy focused
 - Modern interface with a clean design
 - Automatic spam filtering
 - Review and approve or delete comments through the moderation interface
 - Lightweight (under 10 KB total), allowing for faster pageloads
 - Import from existing services
 - Completely free and open source (MIT license)

### Editions

There are three editions of Commento:

 - **Commento Community Edition** (CE) is open source software that's freely available under the MIT Expat license.
 - **Commento Enterprise Edition** (EE) includes extra features geared towards organizations. If you're interested in this, [contact me](mailto:c.adhityaa@gmail.com).
 - [**Commento Hosted**](https://commento.io) is a service version of Commento if you want to embed Commento without running servers.  This is currently in private beta and you can [add yourself to the waiting list here](https://commento.io).

### Getting Started

If you're using version `1.x`, you can find [documentation here](https://docs.commento.io).

Below, you'll find documentation for `0.x`.

#### Hosting the Backend

If you want to self-host Commento, you have three options:

##### Run Binary

Go makes deployment easy by producing a single binary. This is the recommended way, as it has the least amount of abstraction. Download the [latest release](https://github.com/adtac/commento/releases/latest) for your operating system and architecture, extract, and run. For example, if I'm on linux-amd64, I would do:

```bash
$ tar xvf commento_0.2.1_linux_amd64.tar.gz
$ ./commento
```

This would expose the server on port `8080`. Point your Commento frontend configuration to this.

##### Use Docker

You can also use Docker to host Commento. A minimal Docker image is provided: [`adtac/commento`](https://hub.docker.com/r/adtac/commento/). You can get a container running by pulling the image and starting it to expose the server on port `8080`:

```bash
$ docker pull adtac/commento
$ docker run -it -d -p 8080:8080 adtac/commento
```

##### Build from Source

Commento is written in Go. Build the binary from source and run the produced executable:

```bash
$ go build -i -v -o commento
$ ./commento
```

If you're building the project for the first time, `go build` might take a few seconds since Commento's dependencies need to be pulled and compiled as well. However, subsequent builds will be faster.

#### Frontend Integration

To embed Commento in your website, paste the following HTML snippet wherever you want Commento to load:

```html
<div id="commento"></div>
<script defer
  src="http://server.com/assets/js/commento.min.js"
  data-div="#commento">
</script>
```

Commento will simply fill the container it is placed in. Remember to change `server.com` to your server.

### Configuration

#### Configuring the Backend

| Parameter | Default Value | Meaning |
| --------- | ------------- | ------- |
| `COMMENTO_PORT` | 8080 | Default port on which the server will listen. |
| `COMMENTO_DATABASE_FILE` | `commento.sqlite3` | Database file that Commento will use to store comments. |
| `AKISMET_KEY` | | Your Akismet key. Read [our docs](docs/akismet.md) for more details on how to obtain one. |

Commento uses environment variables to configure parameters. You can either use a `.env` file or give parameters through the command line. For example, a particular configuration can be achieved in three different ways:

```bash
$ cat .env
COMMENTO_PORT=9000
COMMENTO_DATABASE_FILE=/app/commento.db
AKISMET_KEY=abcdef012345
$ ./commento
```

```bash
$ COMMENTO_PORT=9000 COMMENTO_DATABASE_FILE=/app/commento.db AKISMET_KEY=abcdef012345 ./commento
```

```bash
$ export COMMENTO_PORT=9000
$ export COMMENTO_DATABASE_FILE=/app/commento.db
$ export AKISMET_KEY=abcdef012345
$ ./commento
```

Note that environment variables have precedence over `.env` values. If you're using Docker, you can pass environment variables too:

```bash
$ docker run adtac/commento -it -d -p 9000:9000 \
    -e COMMENTO_PORT=9000                       \
    -e COMMENTO_DATABASE_FILE=/app/commento.db  \
    -e AKISMET_KEY=abcdef012345
```

### Purpose

If you run a blog, and you want your readers to converse, you'll probably install Disqus. After all, they're the biggest system, with 500 million unique visitors every month<sup>[[1]](https://blog.disqus.com/the-numbers-of-disqus)</sup>. They have a free plan, too! They also have non-optional advertisement, third-party tracking, affiliate marketing, and link hijacking <sup>[[2]](http://donw.io/post/github-comments/)</sup> <sup>[[3]](https://www.cleversprocket.com/disqus-is-parsing-your-dom-and-adding-affiliate-links/)</sup> <sup>[[4]](http://chrislema.com/killed-disqus-commenting/)</sup> <sup>[[5]](https://medium.com/patrickleenyc/beware-of-disqus-17fb58cfab10)</sup>. Your readers will be mercilessly tracked, and you'll have no privacy whatsoever &mdash; even if you don't leave any comments. Thanks to the 90+ downloads Disqus will make, your website will be considerably slower.

Commento aims to solve this. Commento is [free software](https://www.fsf.org/about/what-is-free-software) that you can self-host for the cost of a cup of coffee. No third-party tracking scripts will be injected, there will be no affiliate marketing, no advertisements. Just comments.

While some open source solutions existed, I didn't find any attractive enough at the time Commento was created -- most were either they were discontinued or development was virtually non-existant. Free software is about, well, freedom, so I figured I'd write my own software.

### Contributing

Commento is possible only because of its community. The project is still in beta and I'd be thankful for any contribution. Please go through the [development guidelines](docs/development.md) before you start. If you're a newcomer, you want to go through our [newcomer docs](docs/newcomers.md) first. Pick up any issue and hack away!

If you have any questions, [please ask in our `#commento` IRC channel on freenode](https://irc.commento.io).

### Sponsors

Commento development is sponsored by [Mozilla](https://mozilla.org) and [DigitalOcean](https://www.digitalocean.com/) independently.

<p align="center">
<a href="https://www.mozilla.org/en-US/"><img src="https://user-images.githubusercontent.com/7521600/32265838-d05b2d08-bf0a-11e7-92e1-2cb183eae616.png" title="Mozilla" height="40"></a>
&nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp; &nbsp;
<a href="https://www.digitalocean.com"><img src="https://user-images.githubusercontent.com/7521600/32265839-d093c7da-bf0a-11e7-8d99-96a940041d06.png" title="DigitalOcean" height="40"></a>
</p>

### License

```
Copyright 2018 Adhityaa Chandrasekar

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
