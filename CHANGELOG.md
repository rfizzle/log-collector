# [0.4.0](https://github.com/rfizzle/log-collector/compare/v0.3.1...v0.4.0) (2020-11-13)


### Features

* added cisco umbrella input ([7b18707](https://github.com/rfizzle/log-collector/commit/7b187074ef5a99e7e23578bb60bbef1e9c4afe9a))



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



