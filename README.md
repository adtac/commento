# Commento

An open source, lightweight, and tracking-free comment engine.

![Example](https://cloud.githubusercontent.com/assets/7521600/25356132/d00013e0-2956-11e7-8dba-772a8040ae0c.png)

### Installation

It's really simple to embed a Commento section to your webpage. A trivial page would look like:

```html
<html>
    <body>
        <script src="https://cdn.rawgit.com/adtac/commento/0.0.4/vendor/commento.min.js"></script>
    </body>

    <script>
    init_commento("http://0.0.0.0:8080");
    </script>

    <div id="commento">
    </div>
</html>
```

And that's it! Source the script from CDN, add a `div` called `commento` (which will contain the comments) and initialize Commento with your server.

The client-side script accepts an optional second argument `options`, in the form of a plain object. Currently, the only option is `(boolean) honeypot`, which adds a hidden input field to fool spammers. If anything is input into this field, the submission is silently ignored. This option defaults to `false` if the options param is not set explicitly.

To get the server running, run:

```bash
$ go get -v github.com/adtac/commento
```

and build the project using a `go build .` to get a binary. Internally, I've used sqlite3 as the database. Take a look at the code for more details.

### Why?

[Disqus](https://disqus.com/) is one of the most popular commenting services. However, over the years it has become quite bloated - one [blog post](http://donw.io/post/github-comments/) has a detailed analysis. In short, a Disqus-free page makes about 16 HTTP requests while the same page makes 105 requests when Disqus is enabled! This is mostly due to various tracking services that record every action you take on any website that has Disqus embedded.

I ran a quick test: [go to this codepen](https://codepen.io/ryanbelisle/full/AwLgu/) and open your developer tools. You'll see that the sum total of all network requests related to Disqus comes to about ~250 kB! And there aren't even any comments!

So I thought I'd quickly write a simple comment engine in Go. I've been learning Go for the past month or so and it has been fantastic.

With Commento, you wouldn't need to worry about shady ad companies getting your data through hundreds of tracking services. You wouldn't need to worry about your page being slowed down - **Commento uses just 22 kB total**. And it's all open source.

### Contributing

Commento is extremely simplistic in comparison to Disqus. It does not have voting, moderation, and some of the more advanced stuff. Patches are more than welcome! But do keep in mind the whole purpose of this project - a lightweight comment engine with zero user tracking.

### License

MIT License. See the [LICENSE](LICENSE) file for more information.
