<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Change Log](#change-log)
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

## [v0.0.6](https://github.com/ory/oathkeeper-maester/tree/v0.0.6) (2019-12-20)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.5...v0.0.6)

**Closed issues:**

- Broken handling of rules across namespaces [\#36](https://github.com/ory/oathkeeper-maester/issues/36)
- Add leader election flag to the controller [\#35](https://github.com/ory/oathkeeper-maester/issues/35)
- High memory consumption in clusters with many ConfigMaps [\#32](https://github.com/ory/oathkeeper-maester/issues/32)

**Merged pull requests:**

- Filter only if CM field is set [\#37](https://github.com/ory/oathkeeper-maester/pull/37) ([Demonsthere](https://github.com/Demonsthere))

## [v0.0.5](https://github.com/ory/oathkeeper-maester/tree/v0.0.5) (2019-12-11)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.4...v0.0.5)

**Merged pull requests:**

- Don't watch ConfigMaps [\#33](https://github.com/ory/oathkeeper-maester/pull/33) ([Tomasz-Smelcerz-SAP](https://github.com/Tomasz-Smelcerz-SAP))

## [v0.0.4](https://github.com/ory/oathkeeper-maester/tree/v0.0.4) (2019-12-11)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.3...v0.0.4)

**Merged pull requests:**

- Add sidecar mode [\#31](https://github.com/ory/oathkeeper-maester/pull/31) ([Demonsthere](https://github.com/Demonsthere))

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
- Enable integration tests in the CI [\#15](https://github.com/ory/oathkeeper-maester/pull/15) ([piotrmsc](https://github.com/piotrmsc))
- Bugfix: retry on get and update ops [\#13](https://github.com/ory/oathkeeper-maester/pull/13) ([jakkab](https://github.com/jakkab))

## [v0.0.1-beta.2](https://github.com/ory/oathkeeper-maester/tree/v0.0.1-beta.2) (2019-07-29)
[Full Changelog](https://github.com/ory/oathkeeper-maester/compare/v0.0.1-beta.1...v0.0.1-beta.2)

## [v0.0.1-beta.1](https://github.com/ory/oathkeeper-maester/tree/v0.0.1-beta.1) (2019-07-29)
**Merged pull requests:**

- Move dataKey to cmd parameters [\#12](https://github.com/ory/oathkeeper-maester/pull/12) ([Demonsthere](https://github.com/Demonsthere))
- Release step in the CI [\#11](https://github.com/ory/oathkeeper-maester/pull/11) ([piotrmsc](https://github.com/piotrmsc))
- Rule controller tests [\#8](https://github.com/ory/oathkeeper-maester/pull/8) ([Tomasz-Smelcerz-SAP](https://github.com/Tomasz-Smelcerz-SAP))
- Add basic controller logic and rule validation [\#7](https://github.com/ory/oathkeeper-maester/pull/7) ([jakkab](https://github.com/jakkab))
- Initial changes to circleci config [\#6](https://github.com/ory/oathkeeper-maester/pull/6) ([piotrmsc](https://github.com/piotrmsc))
- Add controller scaffold [\#5](https://github.com/ory/oathkeeper-maester/pull/5) ([kubadz](https://github.com/kubadz))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*