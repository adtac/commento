<p align="center">
<a href="https://commento.io"><img src="https://user-images.githubusercontent.com/7521600/33375172-14b21f68-d52f-11e7-9b30-477682bccf8f.png" width=300></a>
</p>

<p align="center"><b>An open source, privacy focused discussion platform.</b></p>

<p align="center">
<a href="https://demo.commento.io"><img src="https://img.shields.io/badge/demo-live-red.svg?style=for-the-badge&colorA=1e2127&colorB=98c379&label=demo"></a>
<a href="https://irc.commento.io"><img src="https://img.shields.io/badge/irc-%23commento-red.svg?style=for-the-badge&colorA=1e2127&colorB=c678dd&label=freenode"></a>
<a href="https://gitlab.com/commento/commento-ce/container_registry"><img src="https://img.shields.io/badge/live-commento-red.svg?style=for-the-badge&colorA=1e2127&colorB=56b6c2&label=docker"></a>
</p>

<h2 align="center"></h2>

Commento is a discussion platform that you can self-host. You can embed it on your blog, news articles, and any place where you want your readers to add comments. Unlike most alternatives, Commento is lightweight and privacy focused; we'll never sell your data, show ads, embed third-party tracking scripts, or inject affiliate links.

### Features

 - Privacy focused
 - Modern interface with a clean design
 - Automatic spam filtering
 - Review and approve or delete comments through the moderation interface
 - Extremely lightweight, allowing for faster pageloads
 - Import from existing services (like Disqus)
 - Completely free and open source (MIT Expat license)

### Editions

There are three editions of Commento.

 - **Commento Community Edition (CE)** is open source software that's freely available under the MIT Expat license.
 - **Commento Enterprise Edition (EE)** includes extra features geared towards organizations. If you're interested in this, [contact me](mailto:c.adhityaa@gmail.com).
 - [**Commento Hosted**](https://commento.io) is a hosted version of Commento if you don't want to host and manage servers on your own. This is currently in private beta and you can [add yourself to the waiting list here](https://commento.io).

### Installation

The recommended way to install Commento is with [Docker Compose](https://docs.docker.com/compose). Docker Compose allows you to install and manage the service painlessly. [Read our documentation](http://docs.commento.io/installation-docker.html) to find our more on how to get Commento running with Docker Compose.

There are other options to install Commento, including running the binary directly and compiling from source. Please refer to the [installation page](https://docs.commento.io/installation.html) for more information.

Once you've installed the software, you need to configure it with various environment variables before starting the service. To learn more about this, refer to our documentation on [configuring Commento](https://docs.commento.io/configuration.html).

### Contributing

Commento is possible only because of its community. If this is your first contribution to Commento, please go through the [development documentation](https://docs.commento.io/contributing.html) before you begin.

Help will always be given at Commento to those who ask for it. We use IRC for chat to collaborate with other developers. You're invited to [hang out with us](https://irc.commento.io) in the `#commento-dev` channel on freenode if you want to contribute to Commento!

### License

```
Copyright 2018 Commento, Inc.

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
