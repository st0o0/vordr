# Changelog

## [0.1.1](https://github.com/st0o0/vordr/compare/v0.1.0...v0.1.1) (2026-11-20)


### Features

* add embedded web dashboard on GET / ([2c8662a](https://github.com/st0o0/vordr/commit/2c8662ab75737731e3f770c54bee3486179c2493))


### Bug Fixes

* **security:** add trivy.yaml with DS-0002 skip and wire trivy-config ([92197f0](https://github.com/st0o0/vordr/commit/92197f0fd386ea3f78511446dc08cbc1b327ece1))
* **security:** use .trivyignore for DS-0002 skip instead of trivy.yaml ([ce7dddc](https://github.com/st0o0/vordr/commit/ce7dddcef417f03b78dfbf8c614b81afd30581e8))


### Refactoring

* **security:** remove unused trivy.yaml placeholder ([531eb24](https://github.com/st0o0/vordr/commit/531eb240338296acbc9b2ac7acf9c39473a46ba0))

## 0.1.0 (2026-11-19)


* replace prose double-dashes with single dashes ([54e4c82](https://github.com/st0o0/vordr/commit/54e4c8230386448754843784e6a2a600ebbc2c1b))


### Features

* add Dockerfile with bundled Trivy binary and single-container compose setup ([e491658](https://github.com/st0o0/vordr/commit/e49165882e3a244cc9e3d1b8bd99333da9b8c5bf))
* add E2E dev testing infrastructure with docker-compose.dev.yml and test plan ([cd3fe58](https://github.com/st0o0/vordr/commit/cd3fe58e8231a30319c5aeed2459448e08fa5e2e))
* add project scaffolding, CI/CD, and community files ([246db59](https://github.com/st0o0/vordr/commit/246db591b877ec3d0ca32a121516c5b911b40853))
* add Trivy server mode for centralized DB and vordr-server image ([f18e7cb](https://github.com/st0o0/vordr/commit/f18e7cb021eddddb5380f1bf3d924c1d07950e61))
* add v0.2 observability - structured logging, /details endpoint, Trivy DB age metric, Grafana dashboard ([81a000f](https://github.com/st0o0/vordr/commit/81a000fe325e15adb54fd328ec567f7f10f14a23))
* add v0.3 smarts - image exclusion, delta detection, severity threshold, retry backoff ([425403f](https://github.com/st0o0/vordr/commit/425403f5838bab167c076ddca7e9379a66e9f7c5))
* implement v0.1 CVE scanner with Trivy CLI, Docker watcher, worker pool, and Prometheus metrics ([883abad](https://github.com/st0o0/vordr/commit/883abad4b2e208736cad737c977cd59285e969f5))
* **security:** add source scan on PR/push, registry scan after release ([137b8e0](https://github.com/st0o0/vordr/commit/137b8e07e049fa783a2604434d848b42eea0ec98))


### Bug Fixes

* remove unused snapshotToDetails function ([218667e](https://github.com/st0o0/vordr/commit/218667ea5c745195db666ee03c5ec3e99fd07959))
* resolve all Trivy security alerts without suppression ([8b8d44d](https://github.com/st0o0/vordr/commit/8b8d44dc06d23b403acc49fbaa0f4372e25c0f19))
* **security:** add packages:read permission for reusable workflow compatibility ([31c8ffb](https://github.com/st0o0/vordr/commit/31c8ffbc796bd6540b115015dc2e5da6cc3e55cc))


### Documentation

* add README, remove ROADMAP (all milestones implemented) ([1bbd93f](https://github.com/st0o0/vordr/commit/1bbd93f9661d2d404bc7a1200d1894e716a30813))
* add ROADMAP, CLAUDE.md, and OpenSpec planning artifacts ([664cbb0](https://github.com/st0o0/vordr/commit/664cbb0044a96deab2acfa0d6b72f3402118f4d0))
* add v0.2 observability OpenSpec artifacts ([215d4f3](https://github.com/st0o0/vordr/commit/215d4f31fe73cd71eca86fe82ab30f1c4f11d840))
* add v0.3 smarts OpenSpec artifacts ([44f5b42](https://github.com/st0o0/vordr/commit/44f5b42a14eb8c629492cfbb1cb6ac8444020478))
* archive v0.2 change and sync specs ([1387db5](https://github.com/st0o0/vordr/commit/1387db5a5a3b4f0465ba659ae201c4deb293a5ca))
* archive v0.3 change and sync specs ([ad7e2eb](https://github.com/st0o0/vordr/commit/ad7e2ebd1cd6aa89c4cf232ed865646db48598bd))


### Refactoring

* **security:** simplify to fs scan on PR/push and local image scan on schedule ([120004b](https://github.com/st0o0/vordr/commit/120004b2b7be2b0e9e014a4554130b666de7b10f))
