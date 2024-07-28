# Changelog

## [0.6.0-rc1](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.5.0...v0.6.0-rc1) (2024-07-28)


### Features

* automatically prep prereleases ([5bf2257](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/5bf2257bcfc0ed7435a9e62e7321778e14799edc))
* upgrade kubebuilder and add conformance tests ([#112](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/112)) ([43e30b6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/43e30b69aec085067a4f858956dae4745671745c))


### Bug Fixes

* **deps:** update kubernetes packages to v0.30.3 ([#110](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/110)) ([1567ec4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1567ec43296c6f67164ccfc45fdd772df9529c95))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.0 ([#117](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/117)) ([543dcfd](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/543dcfd30b9d14889598a078b5172c55ec4932e1))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.1 ([#121](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/121)) ([bb0d3f1](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/bb0d3f1b21380521d615459391edec2d211dda38))
* **deps:** update module github.com/onsi/gomega to v1.34.0 ([#120](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/120)) ([7e03b53](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/7e03b535d8c3a0899cf0f02d83147ea5857d07bb))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.4 ([#116](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/116)) ([dd724cf](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/dd724cf79f3e4566de8dba528b0435984d3121ef))
* pin manifest version in README ([ae6ad29](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/ae6ad290d13887e23276ed3881ee8653f8508d76)), closes [#123](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/123)
* prereleases CI ([3c11c23](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/3c11c230647c4b140c5586846f821b246496f9dc))
* trigger prereleases manually ([358d95d](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/358d95db074376ac95b2c841e629325c2659a45b))

## [0.5.0](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.4.0...v0.5.0) (2024-06-10)


### Features

* refactor CI, add conformance tests ([1d6f200](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1d6f20021922a6d063d6347f4b94fdc2f7cbb506))


### Bug Fixes

* conformance cluster setup and report ([fa3b2cd](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/fa3b2cdc2455f60995b2e05e91fec87c2ac6d397))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.4 ([#104](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/104)) ([62cc730](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/62cc730295f22518d764ef6df748737ea69ae875))
* increment GatewayClass ObservedGeneration ([9e6d2d7](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9e6d2d75852a8dc16d7fa2ffcd08c6821acf441b))
* show manager logs after conformance tests ([684bb91](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/684bb91ca8f4db02594e3904182568b395791b3e))

## [0.4.0](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.8...v0.4.0) (2024-06-01)


### Features

* cloudflare-go v2, reconcile deployment, expose metrics ([#70](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/70)) ([1d970ba](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1d970baad50336e5c0436be525abbacf9e4fe1a0))


### Bug Fixes

* auto-update base image ([73f115e](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/73f115e72d9ae2dcc822b528682a07d9ca761216))
* **deps:** update kubernetes packages to v0.30.1 ([#95](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/95)) ([ad14c70](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/ad14c70a580cccadc3b3792c4832a3c01d160c97))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.96.0 ([#92](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/92)) ([73b2c33](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/73b2c33633c6c0c84a6777ecd8ffd135e09eb77a))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.3 ([#91](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/91)) ([14e20b7](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/14e20b74ce0a19adbda9950e8e9632e47e101820))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.19.0 ([#97](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/97)) ([6eb4078](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/6eb40786629e354911edd9159d8c8bda08040609))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.1 ([#86](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/86)) ([89a3859](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/89a38591bba65ec8c4acad11d5608c42d32ca965))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.2 ([#88](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/88)) ([5fc2232](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/5fc223234df44914050bdfb7f79bb6056166398d))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.3 ([#98](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/98)) ([271dd03](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/271dd03e88d9b3324208c692449540aa7836f78c))
* **deps:** update module sigs.k8s.io/gateway-api to v1.1.0 ([#93](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/93)) ([9fe3c89](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9fe3c897b230fbe30d93ff31813cbe2da071862d))

## [0.3.8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.7...v0.3.8) (2024-05-01)


### Bug Fixes

* Allow root domain has HTTPRoute hostname ([1e4f5c9](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1e4f5c9161600af8a8507633b2e5dae9c6d95f4f))
* **deps:** update kubernetes packages to v0.29.4 ([#71](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/71)) ([cbbc50a](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/cbbc50a8f47853a4fe7d48ff4040ec955975ca0d))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.92.0 ([#63](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/63)) ([a9f4263](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/a9f42638324e1276c00de21612fed585492afb30))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.93.0 ([#68](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/68)) ([79a5a45](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/79a5a455bfbf1885b2150b301daa3829f1c306ac))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.94.0 ([#75](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/75)) ([9bf659f](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9bf659fc6ea5593dbad354c69f04e6a78157d8fe))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.2 ([#78](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/78)) ([c30685c](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/c30685ce30716cb290d5cab4a7ebbd94ae2d91bd))
* **deps:** update module github.com/onsi/gomega to v1.33.0 ([#73](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/73)) ([fd91cd2](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/fd91cd28f3cf633fb480c4ed6cb40e9f65718023))
* **deps:** update module github.com/onsi/gomega to v1.33.1 ([#79](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/79)) ([d34f295](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/d34f295f91e369a3f7e2de27f9ed3a4d1cca9d3c))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.3 ([#67](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/67)) ([66d9b44](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/66d9b44076982d897cc6c36725e813e4bb8a7e93))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.18.0 ([#76](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/76)) ([1cda9da](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1cda9dac3f77861b9062e89b0b403f37f83e1378))
* Handle HTTPRoute without parentRefs[].namespace ([8477c93](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/8477c935e57802a832c18c6745697af13b7c7e3d))
* Remove /api from build ([8629959](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/8629959b3838b87fd216fb3f79257efccf3db219))
* Specify registry for golang container ([4ed277e](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4ed277e4a2ecdc643c3b36b6744b1f6dc3e2904a))
* Strip spaces from secret data ([9a37c06](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9a37c069e1ad57b9029b8ab37157e42e6dbc52e0))

## [0.3.7](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.6...v0.3.7) (2024-03-24)


### Bug Fixes

* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.1 ([#61](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/61)) ([df981c2](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/df981c20f09287aab0cee185a62bcb1df352fe13))

## [0.3.6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.5...v0.3.6) (2024-03-22)


### Bug Fixes

* **deps:** update module github.com/cloudflare/cloudflare-go to v0.91.0 ([#60](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/60)) ([67bf9f9](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/67bf9f94e2e38f74b245e38271e070568e0d3f6e))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.17.0 ([#57](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/57)) ([5734170](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/57341705ae2f9b28135f9a143fd05fc55ae2ae0e))
* **deps:** update module github.com/onsi/gomega to v1.32.0 ([#56](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/56)) ([08d3f39](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/08d3f394986765fa1a012e55d054d025417f6278))

## [0.3.5](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.4...v0.3.5) (2024-03-16)


### Bug Fixes

* **deps:** update kubernetes packages to v0.29.3 ([#53](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/53)) ([385fb65](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/385fb656f65f3116c08fe7e4149ee8f05ebeadc9))

## [0.3.4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.3...v0.3.4) (2024-03-13)


### Bug Fixes

* **deps:** update module github.com/cloudflare/cloudflare-go to v0.90.0 ([#51](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/51)) ([acb3602](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/acb36026538c22dbad807cfd314b2256b90114a1))

## [0.3.3](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.2...v0.3.3) (2024-03-05)


### Bug Fixes

* **deps:** update module github.com/onsi/ginkgo/v2 to v2.16.0 ([#49](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/49)) ([b84acfc](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/b84acfc9dac8bd9b8495777135d59671130dc23b))

## [0.3.2](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.1...v0.3.2) (2024-02-29)


### Bug Fixes

* **deps:** update kubernetes packages to v0.29.2 ([#45](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/45)) ([efade8a](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/efade8af639839e4e9c20ac2a7b649764488bf5f))
* **deps:** update module github.com/cloudflare-go to v0.87.0 ([#37](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/37)) ([c182615](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/c182615703e6b197e7fe87276d1248d957431e88))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.88.0 ([#44](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/44)) ([2308518](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2308518c587b77d9267c838feb4e0f642d3f93df))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.89.0 ([#48](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/48)) ([f93f20e](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/f93f20edb3561a8cfdb5d624918add508ca1002e))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.1 ([#42](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/42)) ([fff7616](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/fff76165cf1ee4d9252555f50fb87d1581b7eb2a))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.2 ([#46](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/46)) ([5c2dac4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/5c2dac4ae629524fe6e071c97cb5ab01ba92c815))

## [0.3.1](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.3.0...v0.3.1) (2024-01-23)


### Bug Fixes

* **deps:** add renovate config for cloudflared ([#32](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/32)) ([a818429](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/a818429ba00c50d413522e6b2b8eb12f979c3c6d))

## [0.3.0](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.2.1...v0.3.0) (2024-01-20)


### Features

* add GoReleaser ([cecb8a7](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/cecb8a713f39ae559373112c0d39f3954b832cff)), closes [#18](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/18)
* automatic semver releases ([3b88861](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/3b8886144fee041a22271dcf6152608c475aa94e))
* implement HTTPRoute DNS records ([15063c5](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/15063c567a8fdba020053b0b62ba9b3672a19839))
* initial Gateway implementation ([9b84572](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9b84572eb4b1745cc6b19e42287b348fae25940c))
* initial HTTPRoute implementation ([cfa93af](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/cfa93afe4870d34d60c48e05603c205d3e22ed66))
* reconcile deleted HTTPRoutes ([08d45d8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/08d45d8647ccac8b4cc0c14c3e4d2516c7c4cee1))
* reconcile sibling HTTPRoutes ([5cfc97d](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/5cfc97dada78404c59efa97a864d1f510b392f16))
* scaffolding ([60ed23e](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/60ed23ed261779bf2b0bd4d55c0fcf529f00b686))
* support all ko platforms, enable seccomp, fix entrypoint ([7dd0822](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/7dd0822bc5d925bdbdb46b2018325fefe5526b0f))
* upgrade to Gateway v1, implement GatewayClass ([13c0c7e](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/13c0c7e37792ff7d24a642b8782f7f8deacfabde))


### Bug Fixes

* add license ([72d67d4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/72d67d4144d307ec9f28fca62fe61e141dd004bd))
* CI build target ([a0b0c6b](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/a0b0c6b9334cd98459028ed5d471c8684ee6411a))
* CI go version ([91ac525](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/91ac525e09654717a176e00a61129202ce33837b))
* controller runAsUser ([4e83dfb](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4e83dfbe09105380ae80cbb91b4a7247790fc505))
* **deps:** update kubernetes packages to v0.28.4 ([#13](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/13)) ([4d3cdd8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4d3cdd856488d59889559f74b0ada9de33d08128))
* **deps:** update kubernetes packages to v0.29.0 ([#17](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/17)) ([c02326b](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/c02326bec556f5092e1d59399764bd0edf68f8bd))
* **deps:** update kubernetes packages to v0.29.1 ([#25](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/25)) ([d8c8a44](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/d8c8a443a39471b825655a4ffaf2ca8e5a970a68))
* **deps:** update module cloudflare-go to v0.86.0 ([#23](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/23)) ([d1309c6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/d1309c6f037d6fb276d3406f9b0029b4fd0ead4e))
* **deps:** update module cloudflare/cloudflare-go to v0.85.0 ([#19](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/19)) ([1050a0a](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/1050a0ad22d2eafed91b400b5e55cac6fb53c0ed))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.81.0 ([#10](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/10)) ([eb8bcf8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/eb8bcf86842210800fbb5f435f36cea01402c4e4))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.82.0 ([#14](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/14)) ([792d154](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/792d1545099f1a258be033d595c096b899869c18))
* **deps:** update module github.com/cloudflare/cloudflare-go to v0.84.0 ([#16](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/16)) ([fdda5a1](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/fdda5a1640ceef1e25d492898585e9b481562fb4))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.0 ([#6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/6)) ([9b1ad65](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9b1ad651a7853d109c0cb11eb890264598971e71))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.1 ([#12](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/12)) ([c832e5d](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/c832e5deccf25bebf4c1d29d9510faed57214182))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.13.2 ([#15](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/15)) ([2c8ad88](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2c8ad88cdc113a579cde1bfde00693eff1a06339))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.14.0 ([#21](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/21)) ([a152934](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/a1529340bf3a0e9bfebdb7ca981576f3637f4d79))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.15.0 ([#26](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/26)) ([f8af220](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/f8af220b7adfd9ceb627f9c89510cd15cbd901f7))
* **deps:** update module github.com/onsi/gomega to v1.29.0 ([#7](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/7)) ([39ea399](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/39ea39979b0c453fd6c0eaa1d9910dad3107741f))
* **deps:** update module github.com/onsi/gomega to v1.30.0 ([#11](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/11)) ([3f2ebfd](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/3f2ebfd17562eb4a57270d69e650a1445de0a052))
* **deps:** update module github.com/onsi/gomega to v1.31.1 ([#28](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/28)) ([6d9f723](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/6d9f723b231107ef74ac27b883b68b01aa4aaa2b))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.3 ([e36d2b6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/e36d2b66b3dd02b64804d06ed718a60aa312430c))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.16.3 ([0e5ad53](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/0e5ad53e1691da7af3154c7bb3cce441124d41e0))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.0 ([#22](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/22)) ([9397b0a](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9397b0a25f79aced2167107c94aa9690b9818d46))
* DNS record comment, supported platforms ([e0f31e5](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/e0f31e51fc18988627c6170b70ac35005ecd5fac))
* GatewayClass controller name ([f8192fe](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/f8192fe00c75f7464ea39e50d42dcc3235d59b72))
* graceful exit on incorrect permission edge cases ([dc256fc](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/dc256fc8a51206c5bef74f0a069147d7f2e1118f))
* linting ([96166e2](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/96166e20d48aa4de8fb6fee9527913a4123051a5))
* manager kustomization tag, prep for v0.1.4 ([2e1fddd](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2e1fddd1b75e1e5446f6a69f43edf9305bd203c2))
* panic if Gateway hasn't created tunnel yet ([459ac81](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/459ac81887fec8ea52e8c31858fda5d4c929067e))
* pin devcontainer version ([ac45335](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/ac4533519c00af6145cb3786e13fcb2c3fa9135a))
* prep for 0.2.0 release ([4c452b8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4c452b88a2889000f1d2f92f61e84599b41283e4))
* README example apiVersions ([2c7b2f4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2c7b2f43aa2b3b3f6be5c63a0eb602e5c15d8f60))
* rename namespace to cloudflare-gateway ([b0f6f80](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/b0f6f80a321de094e312098cc52b8fc80dcc4e6e))
* various bugs from testing ([2b2a21f](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2b2a21f5fd1a9b4a47f762f887e44b6da07f8c44))
* Windows support with golang base image ([106120f](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/106120fadd4377ab561088f906c14c43c16f5438))

## [0.2.1](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.2.0...v0.2.1) (2024-01-20)


### Bug Fixes

* **deps:** update kubernetes packages to v0.29.1 ([#25](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/25)) ([d8c8a44](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/d8c8a443a39471b825655a4ffaf2ca8e5a970a68))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.15.0 ([#26](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/26)) ([f8af220](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/f8af220b7adfd9ceb627f9c89510cd15cbd901f7))
* **deps:** update module github.com/onsi/gomega to v1.31.1 ([#28](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/28)) ([6d9f723](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/6d9f723b231107ef74ac27b883b68b01aa4aaa2b))

## [0.2.0](https://github.com/pl4nty/cloudflare-kubernetes-gateway/compare/v0.1.4...v0.2.0) (2024-01-17)


### Features

* automatic semver releases ([3b88861](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/3b8886144fee041a22271dcf6152608c475aa94e))


### Bug Fixes

* controller runAsUser ([4e83dfb](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4e83dfbe09105380ae80cbb91b4a7247790fc505))
* **deps:** update module cloudflare-go to v0.86.0 ([#23](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/23)) ([d1309c6](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/d1309c6f037d6fb276d3406f9b0029b4fd0ead4e))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.14.0 ([#21](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/21)) ([a152934](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/a1529340bf3a0e9bfebdb7ca981576f3637f4d79))
* **deps:** update module sigs.k8s.io/controller-runtime to v0.17.0 ([#22](https://github.com/pl4nty/cloudflare-kubernetes-gateway/issues/22)) ([9397b0a](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/9397b0a25f79aced2167107c94aa9690b9818d46))
* pin devcontainer version ([ac45335](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/ac4533519c00af6145cb3786e13fcb2c3fa9135a))
* prep for 0.2.0 release ([4c452b8](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/4c452b88a2889000f1d2f92f61e84599b41283e4))
* README example apiVersions ([2c7b2f4](https://github.com/pl4nty/cloudflare-kubernetes-gateway/commit/2c7b2f43aa2b3b3f6be5c63a0eb602e5c15d8f60))
