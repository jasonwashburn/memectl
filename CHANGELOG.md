# Changelog

## [0.4.0](https://github.com/jasonwashburn/memectl/compare/v0.3.0...v0.4.0) (2026-08-26)


### ⚠ BREAKING CHANGES

* `memectl create meme` now requires a local name and `--template <template-id>`. Positional template IDs are no longer supported.

### Features

* add delete meme command ([#35](https://github.com/jasonwashburn/memectl/issues/35)) ([75510f4](https://github.com/jasonwashburn/memectl/commit/75510f45a87025e02a58016bd805b253fa41b73c))
* add describe meme command ([#37](https://github.com/jasonwashburn/memectl/issues/37)) ([eb89177](https://github.com/jasonwashburn/memectl/commit/eb891771d5c13b2893bfd63c446b23fb374e2120))
* add managed meme inventory ([#31](https://github.com/jasonwashburn/memectl/issues/31)) ([db149ff](https://github.com/jasonwashburn/memectl/commit/db149ff364bbf97ffecba47c6412d85d664baf82))


### Continuous Integration

* add test workflow ([#33](https://github.com/jasonwashburn/memectl/issues/33)) ([b7fbdfe](https://github.com/jasonwashburn/memectl/commit/b7fbdfe13ce212361fad30604e839986b6e48c16))


### Specifications

* archive add-delete-meme ([#36](https://github.com/jasonwashburn/memectl/issues/36)) ([a832733](https://github.com/jasonwashburn/memectl/commit/a8327338c0791c9d44adb0248ae8fbd1f254d455))
* archive add-describe-meme ([#38](https://github.com/jasonwashburn/memectl/issues/38)) ([3208b56](https://github.com/jasonwashburn/memectl/commit/3208b5644010518778eedf037e6f3cd57cdd8938))
* archive add-get-templates-wide-output ([#29](https://github.com/jasonwashburn/memectl/issues/29)) ([611a4f7](https://github.com/jasonwashburn/memectl/commit/611a4f786be75877ead51baeb7d16c8638bed65e))
* archive add-managed-meme-inventory ([#32](https://github.com/jasonwashburn/memectl/issues/32)) ([0113daa](https://github.com/jasonwashburn/memectl/commit/0113daacb475f650d48d07ae8bf9286eacb02505))
* archive add-testing-in-ci ([#34](https://github.com/jasonwashburn/memectl/issues/34)) ([f0ffbd4](https://github.com/jasonwashburn/memectl/commit/f0ffbd44475bb514b8c4424629ce8e47ac98e53c))

## [0.3.0](https://github.com/jasonwashburn/memectl/compare/v0.2.0...v0.3.0) (2026-08-23)


### Features

* add wide output format for get templates command ([#28](https://github.com/jasonwashburn/memectl/issues/28)) ([e5c3971](https://github.com/jasonwashburn/memectl/commit/e5c397109093ba5c24478203b832cf2fe469d4b6))


### Documentation

* add CONTRIBUTING.md and update README.md ([#26](https://github.com/jasonwashburn/memectl/issues/26)) ([8791938](https://github.com/jasonwashburn/memectl/commit/8791938a71453ec4ea464bf7ea6287b48b65217c))


### Tests

* use testify ([#25](https://github.com/jasonwashburn/memectl/issues/25)) ([3716bca](https://github.com/jasonwashburn/memectl/commit/3716bca76460314c89bc32f9427ff4789d85a66b))


### Continuous Integration

* add dependency cooldowns to dependabot config ([#22](https://github.com/jasonwashburn/memectl/issues/22)) ([185f740](https://github.com/jasonwashburn/memectl/commit/185f740fc6cf3323567519639f12a147c3c323f7))
* improve zizmor, mise, and hk configs ([185f740](https://github.com/jasonwashburn/memectl/commit/185f740fc6cf3323567519639f12a147c3c323f7))


### Maintenance

* **deps:** bump goreleaser/goreleaser-action from 6.4.0 to 7.2.3 ([#24](https://github.com/jasonwashburn/memectl/issues/24)) ([d3a0868](https://github.com/jasonwashburn/memectl/commit/d3a08688cbcfd0d415a49a01a8e3f7fb5a6444b8))


### Specifications

* archive improve-documentation ([#27](https://github.com/jasonwashburn/memectl/issues/27)) ([ed69f2c](https://github.com/jasonwashburn/memectl/commit/ed69f2c3270020f3473235dc3c022a5440719f57))

## [0.2.0](https://github.com/jasonwashburn/memectl/compare/v0.1.0...v0.2.0) (2026-08-23)


### Features

* add create meme command ([#20](https://github.com/jasonwashburn/memectl/issues/20)) ([33ca0f1](https://github.com/jasonwashburn/memectl/commit/33ca0f1f1e320e0fb33dd61f9d1fe777e44e1a79))


### Specifications

* archive add-automated-release-builds ([#18](https://github.com/jasonwashburn/memectl/issues/18)) ([dd40111](https://github.com/jasonwashburn/memectl/commit/dd40111a1613d7c79d8bfa62f9d12e6cd5e71a37))
* archive create-captioned-meme ([#21](https://github.com/jasonwashburn/memectl/issues/21)) ([6830cae](https://github.com/jasonwashburn/memectl/commit/6830cae05051a3260fc3eda56005602d11c7d391))

## [0.1.0](https://github.com/jasonwashburn/memectl/compare/v0.0.1...v0.1.0) (2026-08-22)


### Features

* add get templates command ([#6](https://github.com/jasonwashburn/memectl/issues/6)) ([006eb05](https://github.com/jasonwashburn/memectl/commit/006eb05c082fa58f04922b28a8d9e533adbf016d))


### Documentation

* update readme ([#4](https://github.com/jasonwashburn/memectl/issues/4)) ([9567643](https://github.com/jasonwashburn/memectl/commit/956764353d699c90d55c94b1a4892b262d92b661))


### Tests

* add mise coverage task ([#9](https://github.com/jasonwashburn/memectl/issues/9)) ([70350e7](https://github.com/jasonwashburn/memectl/commit/70350e7ce8e2d4c14e6100e0d25551cff3dc8f1e))


### Continuous Integration

* add automated release builds ([#16](https://github.com/jasonwashburn/memectl/issues/16)) ([e27c5b2](https://github.com/jasonwashburn/memectl/commit/e27c5b23f99de8f20a4605e25c4d87914d9dfbe4))
* add detail to automated changelogs ([#17](https://github.com/jasonwashburn/memectl/issues/17)) ([92e9c84](https://github.com/jasonwashburn/memectl/commit/92e9c84aa12005d764a4aa1e7364da6ca3f2455f))
* add zizmor ([#15](https://github.com/jasonwashburn/memectl/issues/15)) ([03170e8](https://github.com/jasonwashburn/memectl/commit/03170e8fbcfb24d4ebfcf4e9884f5025fbaa2028))
* configure dependabot for gomod and github actions ([#14](https://github.com/jasonwashburn/memectl/issues/14)) ([8ecfcf0](https://github.com/jasonwashburn/memectl/commit/8ecfcf0ec6e5e5a67b75e773291a081e6a090717))


### Maintenance

* add hk-based pre-commit configuration and CI hook validation ([#10](https://github.com/jasonwashburn/memectl/issues/10)) ([4046994](https://github.com/jasonwashburn/memectl/commit/40469944d1326950264c5abeb0f973db0121f4d3))
* add release-please workflow and configuration ([#8](https://github.com/jasonwashburn/memectl/issues/8)) ([232a9b4](https://github.com/jasonwashburn/memectl/commit/232a9b4a2abf9ccae5bc3070fcce071201b2bd20))
* add worktrunk post-start copy ([#3](https://github.com/jasonwashburn/memectl/issues/3)) ([871ceca](https://github.com/jasonwashburn/memectl/commit/871ceca0d0614cda568d3411ab8bb81801cde1e9))
* bootstrap openspec ([#5](https://github.com/jasonwashburn/memectl/issues/5)) ([f295fb3](https://github.com/jasonwashburn/memectl/commit/f295fb37cfcfa857ead556d30616d713d9f8373a))
* initial project skeleton setup ([#1](https://github.com/jasonwashburn/memectl/issues/1)) ([48ed8a0](https://github.com/jasonwashburn/memectl/commit/48ed8a064c83687207d522e8b736dcbe133bfac4))
* set up worktrunk for local mise toml ([#2](https://github.com/jasonwashburn/memectl/issues/2)) ([cf1a3cf](https://github.com/jasonwashburn/memectl/commit/cf1a3cfb4941c70abdbac8d28f4ff7002d60ec5d))
* update release-please manifest ([#13](https://github.com/jasonwashburn/memectl/issues/13)) ([2724519](https://github.com/jasonwashburn/memectl/commit/272451973b8dfb8c8aeaa451477871e997725114))


### Specifications

* archive completed add-get-templates ([#7](https://github.com/jasonwashburn/memectl/issues/7)) ([2e59b9c](https://github.com/jasonwashburn/memectl/commit/2e59b9c22166c3ddd73c194031e156bc74e94d3e))
