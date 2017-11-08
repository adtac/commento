## Contributing to Commento

Thank you for choosing to contribute to Commento. The project is still in beta
and any help that I can get is welcome. You could be fixing a typo somewhere or
churning out an entirely new feature - all of it is welcome. There are, however,
certain guidelines I'd like for the contributors to follow. It's okay to have a
slightly messy setup if the project is in pre-pre-alpha (whatever that means),
but as the project matures, its documentation and guidelines should too.

If you're a newcomer to open source and you haven't contributed to projects
before, you should consult the [newcomer docs](newcomer.md) first.

### Getting Started

#### Backend development

Commento uses the `dep` tool to manage dependencies. To retrieve all of
Commento's dependencies, simply do the following from the project root:

```bash
$ dep ensure
```

That's it. If you're adding a new dependency to the project, you can make
appropriate changes in the `Gopkg.toml` file.

After retrieving the dependencies, you can build the project:

```bash
$ go build -i -v -o commento
```

This will produce a binary. Run that and you'll have a backend server running.

#### Frontend developement

Commento has a simple and clean interface. To get started with frontend
development, first install `npm` on your machine. The way to do this varies on
each platform. Once that's done, install the frontend pipeline build
dependencies by running:

```bash
$ npm install
```

Now you can build the frontend minified files. To do that, simply do:

```bash
$ npm build
```

### Guidelines

#### Code Standards

Commento uses [coala](https://coala.io) for static code analysis and linting.
I'm a maintainer at coala too, so if you have any questions, you can direct them
towards me. Essentially, `go fmt` and `go vet` will be run on the project to
make sure your pull request is well-formatted.

This is by no means foolproof. That's why manual review is always done before a
pull request is merged.

#### Commit Guidelines

I believe that commit messages are code themselves. They document the changes
being and should be well explained. With that in mind, here's an example of a
good commit message:

    main.go: Move config parsing to config.go

This message straight away tells us several things:

* We're making this change in `main.go`.
* This change is moving code out of `main.go` into `config.go`.
* The code being moved out is the configuration parsing part.

Here's an example of a bad commit message:

    Fix tests

This doesn't tell us much apart from the fact that the commit is fixing tests.
Which file's tests are being fixed? Which individual test is being changed? What
was wrong before? Why is the new behavior the correct one?

#### Testing

Writing tests is almost always a good idea. I generally prefer one integration
test over 5 small unit tests; that's not to say that I don't want unit tests.
Unit tests are extremely useful, but they are what they are - modules testing
one small unit of the source code. It makes more sense to have integration tests
to make sure the entire system works well together instead of having 10 modules
working perfectly but breaking when they are put together.

Coverage is a trendy metric to quantitatively measure tests. However, I don't
think 100% coverage is a useful thing at all. While I'd still like to measure
coverage every now and then, I don't believe that rabidly adding unit tests to
achieve the magical 100% number is constructive.
