# go-mustache

[![Build Status](https://img.shields.io/travis/cbroglie/go-mustache.svg)](https://travis-ci.org/cbroglie/go-mustache)

This repository is a (work in progress) reimplementation of [cbroglie/mustache](https://github.com/cbroglie/mustache), which is a fork of [hoisie/mustache](https://github.com/hoisie/mustache).

I created my fork to clean up the API and add a few small features which I needed for another project, and I was able to do this without making major changes to the code base. And while [cbroglie/mustache](https://github.com/cbroglie/mustache) works well enough for most use cases, it fails ~40% of the official spec tests (though most failures are related to whitespace handling). I set out to fix these test failures thinking only incremental changes would be needed, but ultimately I decided the parser needed to be reimplemented to ever become compliant with the spec, so here we are.

I don't have any ETA for completion, I just plan to work on it from time to time as my schedule allows.
