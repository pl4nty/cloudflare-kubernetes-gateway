# Changelog

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
