# Contributing to Lambda Machine Local

Thank you for your interest in contributing to Lambda Machine Local.

Lambda Machine Local is a downstream fork
of [Docker Machine](https://github.com/docker/machine) with a focus on stability
and security. We
use [libmachine](https://github.com/docker/machine/tree/master/libmachine) in
Lambda Machine Local and see Docker Machine as our upstream project. We have
adopted ideas and code from Docker Machine and we would like to thank upstream
project for that.

Lambda Machine Local is developed and maintained as a **out-of-tree patch set**
on top of libmachine. We regularly rebase our patch set to the latest stable
version of libmachine. You can find the patch set by checking out the git tag
which has the form `vYYMM.X.Y`, where `YY` is the year, `MM` is the month, `X`
and `Y` are point release version numbers.

We welcome bug reports, feedback and code contributions to Lambda Machine Local.

Prior to making a code contribution, we request that you please consider if it
makes sense to make your contribution to upstream Docker Machine instead? Since
we regularly rebase from Docker Machine, the upstream project might also benefit
from your contribution.

Lambda Machine Local also has
a [plugin architecture](https://github.com/docker/machine/pull/1902) similar to
Docker Machine. This is also another avenue available to you to extend and
customize Lambda Machine Local.

If your contribution is Lambda Machine Local specific, kindly please open an
issue before submitting a pull request so we can guide you through the
contribution process.
