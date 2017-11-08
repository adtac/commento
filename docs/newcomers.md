## Newcomer Guide

In this document, I'll explain how to do some basic actions using git and
GitHub. If something isn't clear, Google is your friend: search Google and ask
on StackOverflow.

#### How to Submit a Pull Request

You should really be learning this from a tutorial dedicated to teaching you git
and GitHub, but I would like the recommended practices here as well. For the
purposes of this example, I'll assume you're fixing a typo in the README page.

First, you'll want to fork the repository. You can do this by clicking on the
'Fork' button at the top. This will create an identical copy of the project that
you can write to.

Secondly, you want to clone your fork. The fork exists only in GitHub's
servers. To make changes, you need to download the project to your local hard
disk:

```bash
$ git clone https://github.com/<your username>/commento.git
```

After cloning the repository, you create a branch. This a pretty important step
most people miss. Right after you clone and go into the directory, you're in the
`master` branch of the project. This is essentially a default branch. You do not
want to make changes over. To fix the README typo, create a new branch with:

```bash
$ git checkout -b readme_fix
```

This will create a new branch call `readme_fix`. Now you can make the changes
with your favorite editor. When you're done it's time to commit and push your
changes. To commit, do:

```bash
$ git commit -m "your commit message
```

See the section title *Commit Guidelines* to understand what makes a good
commit. Anyway, now you just need to push your changes. To do this, run:

```bash
$ git push --set-upstream origin readme_fix
```

With that done, you'll have your changes pushed to GitHub's servers. Now all
that's left is to create a pull request. To do this, go to the webpage of your
fork. Click on "New pull request" and fill in the appropriate fields.
