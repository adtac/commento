The Commento Manifesto
----------------------

| Original author         | Version |
|----------------------------------|---------|
| Adhityaa Chandrasekar (`@adtac`) | 1.0   |

This document describes the goals, status and the purpose behind Commento. This
will also serve as a guide for what Commento should and shouldn't be.

### What is Commento?

Commento is an Italian word that stands for *comment*. It is a comment engine that is embeddable in static and dynamic websites. It is free software and anybody can get an instance of the backend running on their server.

### Current status

Currently Commento is still in its early stages. We do not have any protection against spam -- it is in fact trivial to write a simple bash script to continuously spam the server with comments. We're already working on this and there are some really good proposals to tackle spam. When we implement Captcha support and a moderation dashboard, Commento will have very good spam protection.

Several new features are proposed as well. A front-end pipeline to make concurrent frontend development is in the works. We also have a backend router improvement that is going on. Reply notifications through emails is being discussed. Integration tests for the UI and unit tests for the backend have been proposed.

### Why I created Commento

Let's take a look at a couple of options (and why each is non-ideal) if you want to embed comments on your blog:

 - **Disqus**: Disqus is probably the biggest commenting system. In 2011 Disqus had about 500 million unique visitors every month<sup>[[1]](https://blog.disqus.com/the-numbers-of-disqus)</sup>. However, a recent blog post<sup>[[2]](http://donw.io/post/github-comments/)</sup> revealed that every embedded Disqus frame made about 90 network requests increasing the load time by a full 4 seconds. To make matters worse, it was discovered that they used tens of third-party tracking services on all pages. This is an obvious violation of user privacy: a simple commenting system should never need that many requests, handing over user data to that many third-party tracking services.

 - **Facebook comments**: Facebook is one of the biggest data collecting companies in the world - every Share button and every Facebook comments box tracks every move you do. I just opened a random website with embedded Facebook comments and requests to `facebook.com` accounted for 1.5 MB of the 2.4 MB transferred to load the whole page. This included 87 network requests and 35 Javascript files -- and this didn't even load all the comments!

While some open source solutions exist, I didn't find any attractive enough -- either they were discontinued or development was virtually non-existant. Open source is about choice so I figured I'd write my own software.

### Core principles

This section covers what Commento stands for and what it shouldn't be.

 - Commento is (and will always be) free software in the truest sense. You are free to fork your own copy, run it as a service and charge your users, run it behind a closed platform or anything else. I only ask that you include the copyright and permission notice in all copies. The project is licensed under the MIT License.

 - Commento will **never** track its users. Currently we don't store any cookies on your machine; but even if we do in the future, it will be temporary and short-lived. We do not load and run any external Javascript (except for a markdown parser). We may require an email to post comments in the future, but that will only be used to verify users to reduce spam. And remember: all the data is owned by the entity running the Commento instance. That means that if you provide your email address to comment on a friend's blog, then you're giving your email address only to your friend - not to a corporation like Disqus or Facebook.

 - Being lightweight is one of the most important ideal I had when I was creating Commento. Embedding Commento in your website will increase your page load size by just a couple of tens of kilobytes. That's it.

### What Commento is not

Commento is **not** a multi-website commenting platform. It is simply impossible to be a centralized commenting platform *and* not track users. This is, of course, a bit inconvenient - you would have to verify yourself each time you use a different website (only once for a particular website though). But I strongly believe this is a necessity to stick to the core principles.

The Commento project is not a service either. You must setup your an instance in your own instance to get it running. However, I'm planning on offering a service (and charge for a subscription so that I can pay the server bills) but that will be entirely separate and independent from the open source project's direction.

### How you can contribute

Commento is has two parts: the backend that's written purely in Golang and the frontend client that's written in Javascript and CSS. You can contribute to either. Just go to our Github issue tracker and pick up any unassigned issue you like and start hacking!
