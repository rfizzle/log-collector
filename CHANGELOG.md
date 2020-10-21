## [0.1.2](https://github.com/rfizzle/log-collector/compare/v0.1.1...v0.1.2) (2020-10-21)


### Bug Fixes

* added collector /etc folder in Dockerfile ([e3bffe4](https://github.com/rfizzle/log-collector/commit/e3bffe4b3aaade0a39d048fa21f83b69b5dd5ef8))
* fixed docker build push context bug https://github.com/docker/build-push-action/issues/160 ([edac0a7](https://github.com/rfizzle/log-collector/commit/edac0a7c3ec29ca644154ffb567f925f255d9e2d))
* fixed Docker Hub cred reference in release.yml ([b822a08](https://github.com/rfizzle/log-collector/commit/b822a08accc250f3dadfc77826b4b1f97b9af5c8))
* fixed Dockerfile cache copy issue https://github.com/moby/moby/issues/37965 ([9069cd0](https://github.com/rfizzle/log-collector/commit/9069cd0fd4ac4d1e5e9e647be9690c8c9787b4e0))
* fixed tag reference in docker build CI step ([48a030e](https://github.com/rfizzle/log-collector/commit/48a030efe9989d1bf9fbc7da0af1c4b3756f1a66))
* fixed timer in poll and stream methods ([cd803a4](https://github.com/rfizzle/log-collector/commit/cd803a4ec830d0c74657834ac7c9de3d7856bbab))



## [0.1.1](https://github.com/rfizzle/log-collector/compare/v0.1.0...v0.1.1) (2020-10-20)


### Bug Fixes

* fixed formatting in release pipeline ([82149cc](https://github.com/rfizzle/log-collector/commit/82149cc0557a5c15c0ee9b7aeda8cb4bab967962))
* fixed id syntax in release pipeline ([538014e](https://github.com/rfizzle/log-collector/commit/538014eb8bca696a360fb5bb2ebd4532f5ebbc4a))
* fixed name syntax in release pipeline ([2b1810f](https://github.com/rfizzle/log-collector/commit/2b1810fe477724129df580f165c74abdef6cca63))
* removed redundant release pipeline step after finding additional step output variable ([f6d872c](https://github.com/rfizzle/log-collector/commit/f6d872c859e3b1d822cb09a1efffbdb90b8b7431))
* updated Dockerfile to have a temp directory ([4327d07](https://github.com/rfizzle/log-collector/commit/4327d07632bd82de1d431c54ea27bb28b056b4b7))
* updated release pipeline to have a clean version number for Docker release ([32d27fb](https://github.com/rfizzle/log-collector/commit/32d27fbb31c91f5e6b245f285bf63e28b0046d07))
* updated state to start new collectors off at current time ([25d16e3](https://github.com/rfizzle/log-collector/commit/25d16e340160bc82daa082bbe1f6f9bbdb85759f))



# [0.1.0](https://github.com/rfizzle/log-collector/compare/c5ea3d31e50bd78bd7b01564377f7cb3d711dc93...v0.1.0) (2020-10-20)


### Features

* added Dockerfile ([546046b](https://github.com/rfizzle/log-collector/commit/546046b87222c0f06eb6f803827e867855a16e65))
* added syslog input and support for streaming clients ([fc68cb7](https://github.com/rfizzle/log-collector/commit/fc68cb735141d225505e9569ab04b02aa93d936e))
* initial commit ([c5ea3d3](https://github.com/rfizzle/log-collector/commit/c5ea3d31e50bd78bd7b01564377f7cb3d711dc93))



