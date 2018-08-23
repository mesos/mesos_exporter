# Releasing mesos_exporter

## Updating the CHANGELOG.

Before you begin a release, make sure that the [CHANGELOG](./CHANGELOG.md)
is up to date.

## Tagging the release.

First, increment the version number in the [VERSION](./VERSION)
file.  Next, tag the release. We use semantic version tags prefixed
with the letter "v", e.g. "v0.9.8". When tagging, you should
use an annotated tag. In the body of the tag commit message, put
the [shortlog](https://www.git-scm.com/docs/git-shortlog) to give
a summary of the changes included in the tags,
e.g.  `git shortlog --no-merge HEAD...$PREVIOUS`.

## Creating the Github release.

Once you have pushed a release tag, you need to make a draft release
in the Github interface. Use the [new release](https://github.com/mesos/mesos_exporter/releases/new)
link and save the release as a draft.

At this point, make sure you have a
[Github access token](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line)
that has edit permissions for the repository.

## Building and pushing the release.

Now you are ready to actually build and push the release. To do this,
you just need to run the [promu](https://github.com/prometheus/promu)
tool in the right sequence. Make sure that your
[GOPATH](https://github.com/golang/go/wiki/GOPATH) is set and that
you export the `GITHUB_TOKEN` environment variable with the token
you previously created.


```bash
$ export PATH="$GOPATH/bin:$PATH"
$ export GITHUB_TOKEN=xxxxxxxx
$ go get -u github.com/aktau/github-release
$ go get -u github.com/prometheus/promu
$ crossbuild
$ crossbuild tarballs
$ checksum .tarballs
$ release .tarballs
```

Once this completes successfully, you can edit any release notes
in the Github UI, and publish the release.
