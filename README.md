<p align="center">
<img src="https://user-images.githubusercontent.com/7521600/33375172-14b21f68-d52f-11e7-9b30-477682bccf8f.png">
</p>

<p align="center">A lightweight, open-source, privacy-focused comment engine alternative to Disqus.</p>

<p align="center">
<a href="https://www.patreon.com/adtac"><img src="https://img.shields.io/badge/support-patreon-red.svg?style=for-the-badge&colorA=1e2127&colorB=e06c75&label=support"></a>
<a href="https://commento.adtac.pw"><img src="https://img.shields.io/badge/demo-live-red.svg?style=for-the-badge&colorA=1e2127&colorB=98c379&label=demo"></a>
<a href="https://gitter.im/commento-dev/commento"><img src="https://img.shields.io/badge/live-gitter-red.svg?style=for-the-badge&colorA=1e2127&colorB=c678dd&label=chat"></a>
<a href="https://hub.docker.com/r/adtac/commento/"><img src="https://img.shields.io/badge/live-commento-red.svg?style=for-the-badge&colorA=1e2127&colorB=56b6c2&label=docker"></a>
</p>

<h2 align="center"></h2>

### Introduction

Commento is a comment engine. You can embed it on your blog, news articles, and any place where you want your readers to add comments. It's free software, meaning you are allowed to modify and redistribute the source code. It's lightweight and simple, allowing for lightning fast page loads. It's privacy-focused because that's the way comment engines should be. Embedding comments on your blog on why Vim is the greatest editor shouldn't inject tons of ads and tracking scripts into your reader's browser.

##### Principles and What Commento Isn't

* Commento will be free software forever. You are free to fork your own copy, run it as a service and charge your users, run it behind a closed platform or anything else. I only ask that you include the copyright notice in all copies.
* Commento will be simple, fast, and lightweight. Adding Commento to a page will cost you nothing more than a few kB.
* There won't be an analytics dashboard with a hundred different pie charts telling you every which way the user interacted with the comments box. (There will be a moderation dashboard in the future, however.)
* Commento is not a centralized multi-site commenting system. It's simply impossible to be that and not have some sort of tracking.
* The project isn't a service either. Right now, you need to self-host an instance of Commento. However, I'm working on a service-based option where you can just plug a script into your site and be done with it without having to bother with setting and maintaining servers.

### Getting Started

#### Hosting the Backend

If you want to self-host Commento, you have three options:

##### Using Docker

If you're going down the self-hosting route, using Docker to run Commento is recommended. A minimal Docker image is provided for this: [`adtac/commento`](https://hub.docker.com/r/adtac/commento/). You can get a container running by pulling the image and starting it:

```bash
$ docker pull adtac/commento
$ docker run -it -d -p 80:8080 adtac/commento
```

That's it. This will expose the server on your machine on port `80`; point your Commento frontend configuration to this.

##### Using binary releases

If you don't want to install the whole of Docker, you can simply use the latest release. Right now, binary releases are available for linux-amd64, but upon request, I can easily create and publish stable releases for other operating systems and architectures.

##### Building from source

Commento is written in Go. Building Commento from source is simply a matter of running `go build`. This will produce a binary `commento` for your machine that you can simply execute.

#### Embedding Commento in HTML

Embedding Commento in your static website is easy. Simply load Commento's JS library, add a `div` tag with `id="commento"`, then call the `init` function when you please:

```html
<！-- Load Commento's js library -->
<script src="http://127.0.0.1/assets/commento.min.js"></script>
```

```html
<！-- Where textareas and buttons would go -->
<div id="commento"></div>
```

```html
<！-- Call init() function -->
<script>
window.onload = function() {
    Commento.init({
        serverUrl: "http://127.0.0.1",
    });
}
</script>
```


### Configuration

#### Configuring the Backend

| Parameter | Default Value | Meaning |
| --------- | ------------- | ------- |
| `COMMENTO_PORT` | 8080 | The default port on which the server will listen. |
| `COMMENTO_DATABASE_FILE` | `sqlite3.db` | The database file that Commento will use to store comments. |

Commento uses environment variables as a way of configuring parameters. You can either use `.env` files or give the parameters through the command line. For example, an example `.env` file would be:

```bash
COMMENTO_PORT=8001
COMMENTO_DATABASE_FILE=/app/commento.db
```

You can give the exact same parameters through command line when you're start the docker container:

```bash
$ docker run -it -d -p 80:8001               \
  -e COMMENTO_PORT=8001                      \
  -e COMMENTO_DATABASE_FILE=/app/commento.db \
  adtac/commento
```

#### Configuring the frontend

| Parameter | Default Value | Meaning |
| --------- | ------------- | ------- |
| `serverUrl` | the same domain as the webpage | The backend server URL of the form `https://example.com` without a trailing slash. |
| `honeypot` | `false` | Whether or not to use a honeypot to filter spammers. |

### Purpose

Let's take a look at the popular options if you want to embed comments on your blog:

 - **Disqus**: Disqus is probably the biggest commenting system with 500 million unique visitors every month<sup>[[1]](https://blog.disqus.com/the-numbers-of-disqus)</sup>. However, a recent blog post<sup>[[2]](http://donw.io/post/github-comments/)</sup> revealed that every embedded Disqus frame made about 90 network requests increasing the load time by a full 4 seconds. It was discovered that they used tens of third-party tracking services on all pages. A simple commenting system should neither need that many requests nor hand over so much information to so many third-party tracking agents.

 - **Facebook comments**: Facebook comments is equally worse. I just opened a random website with embedded Facebook comments and requests to `facebook.com` accounted for 1.5 MB of the 2.4 MB transferred to load the whole page. This included 87 network requests and 35 Javascript files. And this didn't even load all the comments -- I had to click a "Load more comments" button!

While some open source solutions exist, I didn't find any attractive enough -- either they were discontinued or development was virtually non-existant. Open source is about choice so I figured I'd write my own software.

### Contributing

Everybody is welcome to contribute to the project. The project is still in beta stage and lacks some nice features such as spam protection, moderation, and live comments and I'd be thankful for any contribution – small or big. Please go through the [development guidelines](docs/development.md) before you start. You can join our [Gitter channel](https://gitter.im/commento-dev/commento) too.

While I eventually plan to make some revenue through the project in the form of a service, I also want everyone to have the option to self-host. That's the single reason why I will keep Commento open-source forever. Besides, I have seen and used some incredible programs and tools that were only possible because several people came together to make it.


### Sponsors

Commento development is sponsored by [DigitalOcean](https://www.digitalocean.com/).

<p align="center">
<a href="https://www.digitalocean.com"><img src="https://user-images.githubusercontent.com/7521600/32265839-d093c7da-bf0a-11e7-8d99-96a940041d06.png" title="DigitalOcean" height="40"></a>
</p>

### License

```
Copyright 2017 Adhityaa Chandrasekar

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
