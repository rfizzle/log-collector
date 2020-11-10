## [0.3.1](https://github.com/rfizzle/log-collector/compare/v0.3.0...v0.3.1) (2020-11-09)


### Bug Fixes

* fixed hang when more events have been written to output that the current count ([8df98c6](https://github.com/rfizzle/log-collector/commit/8df98c61aa78553550ce2407d95499d4262e9cd7))



# [0.3.0](https://github.com/rfizzle/log-collector/compare/v0.2.1...v0.3.0) (2020-11-09)


### Features

* added limit for the span of time to retrieve logs to 2 hours in order to prevent large polls ([02524af](https://github.com/rfizzle/log-collector/commit/02524afd24bb029a459059257e629512cfc6f229))



## [0.2.1](https://github.com/rfizzle/log-collector/compare/v0.2.0...v0.2.1) (2020-11-06)


### Bug Fixes

* added some error handling in microsoft retryable HTTP call ([4744cfa](https://github.com/rfizzle/log-collector/commit/4744cfa7d8aa0d1c65d72abcad00ca2b679c1467))
* changed check to be equal to or greater for request loop ([b7b0b26](https://github.com/rfizzle/log-collector/commit/b7b0b26c7e7529afd044306be93cbe75517d1907))



# [0.2.0](https://github.com/rfizzle/log-collector/compare/v0.1.2...v0.2.0) (2020-10-21)


### Features

* added poll offset option in order to retrieve logs that have a delay ([46c879e](https://github.com/rfizzle/log-collector/commit/46c879e024fc6f779ca3dd11bc222b17f40a76be))



## [0.1.2](https://github.com/rfizzle/log-collector/compare/v0.1.1...v0.1.2) (2020-10-21)


### Bug Fixes

* added collector /etc folder in Dockerfile ([e3bffe4](https://github.com/rfizzle/log-collector/commit/e3bffe4b3aaade0a39d048fa21f83b69b5dd5ef8))
* fixed docker build push context bug https://github.com/docker/build-push-action/issues/160 ([edac0a7](https://github.com/rfizzle/log-collector/commit/edac0a7c3ec29ca644154ffb567f925f255d9e2d))
* fixed Docker Hub cred reference in release.yml ([b822a08](https://github.com/rfizzle/log-collector/commit/b822a08accc250f3dadfc77826b4b1f97b9af5c8))
* fixed Dockerfile cache copy issue https://github.com/moby/moby/issues/37965 ([9069cd0](https://github.com/rfizzle/log-collector/commit/9069cd0fd4ac4d1e5e9e647be9690c8c9787b4e0))
* fixed tag reference in docker build CI step ([48a030e](https://github.com/rfizzle/log-collector/commit/48a030efe9989d1bf9fbc7da0af1c4b3756f1a66))
* fixed timer in poll and stream methods ([cd803a4](https://github.com/rfizzle/log-collector/commit/cd803a4ec830d0c74657834ac7c9de3d7856bbab))



