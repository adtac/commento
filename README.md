# Commento

An open source, lightweight, and tracking-free comment engine.

![Example](https://cloud.githubusercontent.com/assets/7521600/25356132/d00013e0-2956-11e7-8dba-772a8040ae0c.png)

### Installation

**[I've hosted a live demo here. Check it out!](https://adtac.github.io/commento/example/demo.html)**

It's really simple to embed a Commento section to your webpage. A trivial page would look like:

```html
<html>
    <body>
        <script src="http://0.0.0.0:8080/assets/commento.min.js"></script>
    </body>

    <script>
        Commento.init({
            serverUrl: "http://0.0.0.0:8080"
        })
    </script>

    <div id="commento">
    </div>
</html>
```

And that's it! Source the client-side script, add a `div` called `commento` (which will contain the comments) and initialize Commento with your server. The client-side script does all the hard work of building the markup and loading the CSS. The assets themselves (JavaScript and CSS) as served by the go application.

The client-side script accepts an optional second argument `options`, in the form of a plain object. Currently, the only option is `(boolean) honeypot`, which adds a hidden input field to fool spammers. If anything is input into this field, the submission is silently ignored. This option defaults to `false` if the options param is not set explicitly.

To get the server running first you should install [dep](https://github.com/golang/dep), Golangs official dependency management tool.

Once `dep` is installed you can get the dependencies needed by commento by entering `dep ensure`

To build the Commento binary you enter the ./cmd/commento directory and build the project using a `go build` to get a binary. 

The commands to get to a working binary are summarized here:

```bash
$ dep ensure
$ cd cmd/commento
$ go build
```

Internally, I've used sqlite3 as the database. Take a look at the code for more details.

### Why?

[Disqus](https://disqus.com/) is one of the most popular commenting services. However, over the years it has become quite bloated - one [blog post](http://donw.io/post/github-comments/) has a detailed analysis. In short, a Disqus-free page makes about 16 HTTP requests while the same page makes 105 requests when Disqus is enabled! This is mostly due to various tracking services that record every action you take on any website that has Disqus embedded.

I ran a quick test: [go to this codepen](https://codepen.io/ryanbelisle/full/AwLgu/) and open your developer tools. You'll see that the sum total of all network requests related to Disqus comes to about ~250 kB! And there aren't even any comments!

So I thought I'd quickly write a simple comment engine in Go. I've been learning Go for the past month or so and it has been fantastic.

With Commento, you wouldn't need to worry about shady ad companies getting your data through hundreds of tracking services. You wouldn't need to worry about your page being slowed down - **Commento uses just 22 kB total**. And it's all open source.

### Contributing

Commento is extremely simplistic in comparison to Disqus. It does not have voting, moderation, and some of the more advanced stuff. Patches are more than welcome! But do keep in mind the whole purpose of this project - a lightweight comment engine with zero user tracking.

#### Development

To run the server

```bash
$ docker build . -t adtac/commento:VERSION
$ docker run -d -p 8080:8080 adtac/commento:VERSION
```

For the front end any static server will do, you can grab any from [this list](https://gist.github.com/willurd/5720255).

### License

MIT License. See the [LICENSE](LICENSE) file for more information.
