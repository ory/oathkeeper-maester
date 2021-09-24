<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Change Log](#change-log)
  - [v0.1.5 (2021-09-24)](#v015-2021-09-24)
  - [v0.1.4 (2021-05-08)](#v014-2021-05-08)
  - [v0.1.3 (2021-05-07)](#v013-2021-05-07)
  - [v0.1.2 (2021-04-25)](#v012-2021-04-25)
  - [v0.1.1 (2021-03-31)](#v011-2021-03-31)
  - [v0.1.0 (2020-07-14)](#v010-2020-07-14)
  - [v0.0.10 (2020-02-18)](#v0010-2020-02-18)
  - [v0.0.9 (2020-02-18)](#v009-2020-02-18)
  - [v0.0.8+oryOS.15 (2019-12-26)](#v008oryos15-2019-12-26)
  - [v0.0.7 (2019-12-20)](#v007-2019-12-20)
  - [v0.0.6 (2019-12-20)](#v006-2019-12-20)
  - [v0.0.5 (2019-12-11)](#v005-2019-12-11)
  - [v0.0.4 (2019-12-11)](#v004-2019-12-11)
  - [v0.0.3 (2019-11-18)](#v003-2019-11-18)
  - [v0.0.2-beta-2 (2019-11-12)](#v002-beta-2-2019-11-12)
  - [v0.0.2-beta.1 (2019-08-12)](#v002-beta1-2019-08-12)
  - [v0.0.1-beta.3 (2019-07-30)](#v001-beta3-2019-07-30)
  - [v0.0.1-beta.2 (2019-07-29)](#v001-beta2-2019-07-29)
  - [v0.0.1-beta.1 (2019-07-29)](#v001-beta1-2019-07-29)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Change Log

## [v0.1.5](https://github.com/ory/oathkeeper-maester/tree/v0.1.5) (2021-09-24)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.1.4...v0.1.5)

**Closed issues:**

- Release v0.1.3 did not pass status check and was not pushed to docker registry [\#50](https://github.com/ory/oathkeeper-maester/issues/50)
- Handlers "config" attribute is replaced by an empty object [\#49](https://github.com/ory/oathkeeper-maester/issues/49)
- Controller crashes because logger can't create file [\#34](https://github.com/ory/oathkeeper-maester/issues/34)
- Sidecar mode - write access rules to file [\#21](https://github.com/ory/oathkeeper-maester/issues/21)
- Validate access rules against JSON Schema from Oathkeeper Upstream [\#9](https://github.com/ory/oathkeeper-maester/issues/9)

**Merged pull requests:**

- Chore/kind 0.11 [\#52](https://github.com/ory/oathkeeper-maester/pull/52) ([Demonsthere](https://github.com/Demonsthere))
- feat: including option to scope manager to namespace [\#51](https://github.com/ory/oathkeeper-maester/pull/51) ([janiskemper](https://github.com/janiskemper))

## [v0.1.4](https://github.com/ory/oathkeeper-maester/tree/v0.1.4) (2021-05-08)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.1.3...v0.1.4)

## [v0.1.3](https://github.com/ory/oathkeeper-maester/tree/v0.1.3) (2021-05-07)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.1.2...v0.1.3)

**Fixed bugs:**

- The field 'upstream' should not be required [\#43](https://github.com/ory/oathkeeper-maester/issues/43)

**Closed issues:**

- Rewrite integration tests to stretchr/testify [\#10](https://github.com/ory/oathkeeper-maester/issues/10)

**Merged pull requests:**

- build: Update CRDs and k8s dependencies [\#48](https://github.com/ory/oathkeeper-maester/pull/48) ([colunira](https://github.com/colunira))

## [v0.1.2](https://github.com/ory/oathkeeper-maester/tree/v0.1.2) (2021-04-25)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.1.1...v0.1.2)

**Merged pull requests:**

- feat: Make upstream field optional part II [\#47](https://github.com/ory/oathkeeper-maester/pull/47) ([Demonsthere](https://github.com/Demonsthere))

## [v0.1.1](https://github.com/ory/oathkeeper-maester/tree/v0.1.1) (2021-03-31)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.1.0...v0.1.1)

**Merged pull requests:**

- fix: upstream should not be required in Rule CRD [\#45](https://github.com/ory/oathkeeper-maester/pull/45) ([pommelinho](https://github.com/pommelinho))

## [v0.1.0](https://github.com/ory/oathkeeper-maester/tree/v0.1.0) (2020-07-14)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.10...v0.1.0)

**Merged pull requests:**

- build: Update k8s packages to 1.17.8 [\#42](https://github.com/ory/oathkeeper-maester/pull/42) ([colunira](https://github.com/colunira))

## [v0.0.10](https://github.com/ory/oathkeeper-maester/tree/v0.0.10) (2020-02-18)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.9...v0.0.10)

## [v0.0.9](https://github.com/ory/oathkeeper-maester/tree/v0.0.9) (2020-02-18)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.8+oryOS.15...v0.0.9)

**Closed issues:**

- \<.\*\> becomes \u003c.\*\u003e [\#38](https://github.com/ory/oathkeeper-maester/issues/38)

**Merged pull requests:**

- fix\(test\): Add failing test case for escaped chars [\#39](https://github.com/ory/oathkeeper-maester/pull/39) ([aeneasr](https://github.com/aeneasr))

## [v0.0.8+oryOS.15](https://github.com/ory/oathkeeper-maester/tree/v0.0.8+oryOS.15) (2019-12-26)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.7...v0.0.8+oryOS.15)

## [v0.0.7](https://github.com/ory/oathkeeper-maester/tree/v0.0.7) (2019-12-20)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.6...v0.0.7)

## [v0.0.6](https://github.com/ory/oathkeeper-maester/tree/v0.0.6) (2019-12-20)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.5...v0.0.6)

**Closed issues:**

- Broken handling of rules across namespaces [\#36](https://github.com/ory/oathkeeper-maester/issues/36)
- Add leader election flag to the controller [\#35](https://github.com/ory/oathkeeper-maester/issues/35)
- High memory consumption in clusters with many ConfigMaps [\#32](https://github.com/ory/oathkeeper-maester/issues/32)

**Merged pull requests:**

- Filter only if CM field is set [\#37](https://github.com/ory/oathkeeper-maester/pull/37) ([Demonsthere](https://github.com/Demonsthere))
- Enable integration tests in the CI [\#15](https://github.com/ory/oathkeeper-maester/pull/15) ([piotrmsc](https://github.com/piotrmsc))

## [v0.0.5](https://github.com/ory/oathkeeper-maester/tree/v0.0.5) (2019-12-11)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.4...v0.0.5)

**Merged pull requests:**

- Don't watch ConfigMaps [\#33](https://github.com/ory/oathkeeper-maester/pull/33) ([Tomasz-Smelcerz-SAP](https://github.com/Tomasz-Smelcerz-SAP))

## [v0.0.4](https://github.com/ory/oathkeeper-maester/tree/v0.0.4) (2019-12-11)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.3...v0.0.4)

## [v0.0.3](https://github.com/ory/oathkeeper-maester/tree/v0.0.3) (2019-11-18)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.2-beta-2...v0.0.3)

**Closed issues:**

- oathkeeper-maester crashes on updating configmap w/ rules [\#25](https://github.com/ory/oathkeeper-maester/issues/25)

**Merged pull requests:**

- Feature: configurable configmap [\#30](https://github.com/ory/oathkeeper-maester/pull/30) ([paulbdavis](https://github.com/paulbdavis))
- Update CI config [\#29](https://github.com/ory/oathkeeper-maester/pull/29) ([piotrmsc](https://github.com/piotrmsc))

## [v0.0.2-beta-2](https://github.com/ory/oathkeeper-maester/tree/v0.0.2-beta-2) (2019-11-12)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.2-beta.1...v0.0.2-beta-2)

**Closed issues:**

- Update to latest changes in oathkeeper access rule config [\#24](https://github.com/ory/oathkeeper-maester/issues/24)

**Merged pull requests:**

- Refactor retry logic to handle conflict on updates. [\#26](https://github.com/ory/oathkeeper-maester/pull/26) ([Tomasz-Smelcerz-SAP](https://github.com/Tomasz-Smelcerz-SAP))
- Add 'hydrator' mutator to the list of default mutators [\#23](https://github.com/ory/oathkeeper-maester/pull/23) ([kubadz](https://github.com/kubadz))

## [v0.0.2-beta.1](https://github.com/ory/oathkeeper-maester/tree/v0.0.2-beta.1) (2019-08-12)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.1-beta.3...v0.0.2-beta.1)

**Closed issues:**

- Helm charts for the controller [\#4](https://github.com/ory/oathkeeper-maester/issues/4)

**Merged pull requests:**

- Fix [\#22](https://github.com/ory/oathkeeper-maester/pull/22) ([piotrmsc](https://github.com/piotrmsc))
- Support multiple mutators [\#20](https://github.com/ory/oathkeeper-maester/pull/20) ([jakkab](https://github.com/jakkab))
- Update Measter readme [\#19](https://github.com/ory/oathkeeper-maester/pull/19) ([tomekpapiernik](https://github.com/tomekpapiernik))

## [v0.0.1-beta.3](https://github.com/ory/oathkeeper-maester/tree/v0.0.1-beta.3) (2019-07-30)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.1-beta.2...v0.0.1-beta.3)

**Closed issues:**

- Run integration tests on CI [\#14](https://github.com/ory/oathkeeper-maester/issues/14)
- Setup CI/CD [\#3](https://github.com/ory/oathkeeper-maester/issues/3)
- Create controller skeleton  [\#2](https://github.com/ory/oathkeeper-maester/issues/2)
- Define CRD representing access rule [\#1](https://github.com/ory/oathkeeper-maester/issues/1)

**Merged pull requests:**

- Project renaming [\#18](https://github.com/ory/oathkeeper-maester/pull/18) ([piotrmsc](https://github.com/piotrmsc))
- Rename docker image to `oathkeeper-maester` [\#17](https://github.com/ory/oathkeeper-maester/pull/17) ([aeneasr](https://github.com/aeneasr))
- Fix default value for rulesFileName [\#16](https://github.com/ory/oathkeeper-maester/pull/16) ([Demonsthere](https://github.com/Demonsthere))
- Bugfix: retry on get and update ops [\#13](https://github.com/ory/oathkeeper-maester/pull/13) ([jakkab](https://github.com/jakkab))

## [v0.0.1-beta.2](https://github.com/ory/oathkeeper-maester/tree/v0.0.1-beta.2) (2019-07-29)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.1-beta.1...v0.0.1-beta.2)

## [v0.0.1-beta.1](https://github.com/ory/oathkeeper-maester/tree/v0.0.1-beta.1) (2019-07-29)
**Merged pull requests:**

- Add sidecar mode [\#31](https://github.com/ory/oathkeeper-maester/pull/31) ([Demonsthere](https://github.com/Demonsthere))
- Move dataKey to cmd parameters [\#12](https://github.com/ory/oathkeeper-maester/pull/12) ([Demonsthere](https://github.com/Demonsthere))
- Release step in the CI [\#11](https://github.com/ory/oathkeeper-maester/pull/11) ([piotrmsc](https://github.com/piotrmsc))
- Rule controller tests [\#8](https://github.com/ory/oathkeeper-maester/pull/8) ([Tomasz-Smelcerz-SAP](https://github.com/Tomasz-Smelcerz-SAP))
- Add basic controller logic and rule validation [\#7](https://github.com/ory/oathkeeper-maester/pull/7) ([jakkab](https://github.com/jakkab))
- Initial changes to circleci config [\#6](https://github.com/ory/oathkeeper-maester/pull/6) ([piotrmsc](https://github.com/piotrmsc))
- Add controller scaffold [\#5](https://github.com/ory/oathkeeper-maester/pull/5) ([kubadz](https://github.com/kubadz))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*