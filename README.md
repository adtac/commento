<h1 align="center">Commento</h1>

<p align="center">A lightweight, open source, tracking-free comment engine alternative to Disqus.</p>

<p align="center"><a href="https://adtac.github.io/commento/example/demo.html">
Check out a live demo of Commento here.
</a></p>

<p align="center"><img src="https://cloud.githubusercontent.com/assets/7521600/25660780/6a15a546-302b-11e7-8b55-f200ff856797.png" alt="Example"></p>

### Installation

You need to have an instance of Commento running in your server. Commento is written in Go so you'll need to [install Go first](https://golang.org/dl/). Once that's done, you can build the server. First get the source:

```bash
$ go get -v github.com/adtac/commento
```

Then go to the directory and run:

```bash
$ go build .
```

This should generate a `commento` binary.

To build the frontend you need to [install Node.js and `npm`](https://docs.npmjs.com/getting-started/installing-node). To install the build dependencies, run:

```bash
$ npm install
```

To build the frontend files, run:

```bash
$ npm run build
```

To start the server, run `./commento` from the build directory. By default the server will started on port 8080. If you want to change this, you can provide a environment variable. For example, if you want the server running on port `1234`:

```bash
$ PORT=1234 ./commento
```

Now you can embed Commento on your webpage. A trivial page would look like:

```html
<html>
    <head>
        <script src="http://127.0.0.1:8080/assets/commento.min.js"></script>
    </head>

    <script>
        Commento.init({
            serverUrl: "http://127.0.0.1:8080"
        })
    </script>

    <div id="commento">
        <!-- Commento will populate this div with comments -->
    </div>
</html>
```

You can run the entire server inside a Docker container too. To do this, run:

```bash
$ docker build . -t adtac/commento
$ docker run -d -p 8080:8080 adtac/commento
```

and you should have the server available on port `8080` on the IP address of the container.

### Options

The `Commento.init` function takes an object of parameters. This is documented below:

| Option    | Description                                                                                                                                                                                         |
|-----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| serverUrl | A server URL of the form `https://example.com` without a trailing slash. This will be used for all API requests and all the data will be stored on this server. Default: same server as the webpage |
| honeypot  | A boolean that determines whether you want a honeypot field. This is a spam-protection technique where all comment posts with a non-empty `honeypot` field are rejected silently. Default: `false`  |

### Why?

Let's take a look at a couple of options (and why each is non-ideal) if you want to embed comments on your blog:

 - **Disqus**: Disqus is probably the biggest commenting system. In 2011 Disqus had about 500 million unique visitors every month<sup>[[1]](https://blog.disqus.com/the-numbers-of-disqus)</sup>. However, a recent blog post<sup>[[2]](http://donw.io/post/github-comments/)</sup> revealed that every embedded Disqus frame made about 90 network requests increasing the load time by a full 4 seconds. To make matters worse, it was discovered that they used tens of third-party tracking services on all pages. This is an obvious violation of user privacy: a simple commenting system should never need that many requests, handing over user data to that many third-party tracking services.

 - **Facebook comments**: Facebook is one of the biggest data collecting companies in the world - every Share button and every Facebook comments box tracks every move you do. I just opened a random website with embedded Facebook comments and requests to `facebook.com` accounted for 1.5 MB of the 2.4 MB transferred to load the whole page. This included 87 network requests and 35 Javascript files -- and this didn't even load all the comments!

While some open source solutions exist, I didn't find any attractive enough -- either they were discontinued or development was virtually non-existant. Open source is about choice so I figured I'd write my own software.

### Contributing

Please read [The Commento Manifesto](https://github.com/adtac/commento/blob/master/manifesto.md) to understand what the project is and what it isn't. Commento is extremely simplistic in comparison to Disqus. It does not have voting, spam-protection, moderation, and some of the more advanced stuff. Patches are more than welcome!

### License

MIT License. See the [LICENSE](LICENSE) file for more information.

