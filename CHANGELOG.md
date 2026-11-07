# Changelog

## 0.1.0 (2026-11-07)


* replace prose double-dashes with single dashes ([54e4c82](https://github.com/st0o0/vordr/commit/54e4c8230386448754843784e6a2a600ebbc2c1b))


### Features

* add Dockerfile with bundled Trivy binary and single-container compose setup ([e491658](https://github.com/st0o0/vordr/commit/e49165882e3a244cc9e3d1b8bd99333da9b8c5bf))
* add E2E dev testing infrastructure with docker-compose.dev.yml and test plan ([cd3fe58](https://github.com/st0o0/vordr/commit/cd3fe58e8231a30319c5aeed2459448e08fa5e2e))
* add project scaffolding, CI/CD, and community files ([246db59](https://github.com/st0o0/vordr/commit/246db591b877ec3d0ca32a121516c5b911b40853))
* add Trivy server mode for centralized DB and vordr-server image ([f18e7cb](https://github.com/st0o0/vordr/commit/f18e7cb021eddddb5380f1bf3d924c1d07950e61))
* add v0.2 observability - structured logging, /details endpoint, Trivy DB age metric, Grafana dashboard ([81a000f](https://github.com/st0o0/vordr/commit/81a000fe325e15adb54fd328ec567f7f10f14a23))
* add v0.3 smarts - image exclusion, delta detection, severity threshold, retry backoff ([425403f](https://github.com/st0o0/vordr/commit/425403f5838bab167c076ddca7e9379a66e9f7c5))
* implement v0.1 CVE scanner with Trivy CLI, Docker watcher, worker pool, and Prometheus metrics ([883abad](https://github.com/st0o0/vordr/commit/883abad4b2e208736cad737c977cd59285e969f5))


### Bug Fixes

* remove unused snapshotToDetails function ([218667e](https://github.com/st0o0/vordr/commit/218667ea5c745195db666ee03c5ec3e99fd07959))


### Documentation

* add README, remove ROADMAP (all milestones implemented) ([1bbd93f](https://github.com/st0o0/vordr/commit/1bbd93f9661d2d404bc7a1200d1894e716a30813))
* add ROADMAP, CLAUDE.md, and OpenSpec planning artifacts ([664cbb0](https://github.com/st0o0/vordr/commit/664cbb0044a96deab2acfa0d6b72f3402118f4d0))
* add v0.2 observability OpenSpec artifacts ([215d4f3](https://github.com/st0o0/vordr/commit/215d4f31fe73cd71eca86fe82ab30f1c4f11d840))
* add v0.3 smarts OpenSpec artifacts ([44f5b42](https://github.com/st0o0/vordr/commit/44f5b42a14eb8c629492cfbb1cb6ac8444020478))
* archive v0.2 change and sync specs ([1387db5](https://github.com/st0o0/vordr/commit/1387db5a5a3b4f0465ba659ae201c4deb293a5ca))
* archive v0.3 change and sync specs ([ad7e2eb](https://github.com/st0o0/vordr/commit/ad7e2ebd1cd6aa89c4cf232ed865646db48598bd))
