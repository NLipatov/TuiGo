# Changelog

## [0.3.0](https://github.com/NLipatov/TuiGo/compare/v0.2.0...v0.3.0) (2026-06-19)


### ⚠ BREAKING CHANGES

* **core:** remove Frame.CellAt ([#15](https://github.com/NLipatov/TuiGo/issues/15))

### Code Refactoring

* **core:** remove Frame.CellAt ([#15](https://github.com/NLipatov/TuiGo/issues/15)) ([2556489](https://github.com/NLipatov/TuiGo/commit/2556489fbd1d8bba7a2c614cbdbfe909e26f00dd))

## [0.2.0](https://github.com/NLipatov/TuiGo/compare/v0.1.3...v0.2.0) (2026-06-18)


### ⚠ BREAKING CHANGES

* simplify public color and input APIs ([#13](https://github.com/NLipatov/TuiGo/issues/13))

### Code Refactoring

* simplify public color and input APIs ([#13](https://github.com/NLipatov/TuiGo/issues/13)) ([9987f57](https://github.com/NLipatov/TuiGo/commit/9987f5786300bea7eb38b7123d19c503e243147c))

## [0.1.3](https://github.com/NLipatov/TuiGo/compare/v0.1.2...v0.1.3) (2026-06-16)


### Features

* **core:** support variable-width glyphs ([#11](https://github.com/NLipatov/TuiGo/issues/11)) ([5f28744](https://github.com/NLipatov/TuiGo/commit/5f28744370e96d3bf82bc9ff53c6676d245ed174))

## [0.1.2](https://github.com/NLipatov/TuiGo/compare/v0.1.1...v0.1.2) (2026-06-14)


### Bug Fixes

* **deps:** update Go dependencies ([#8](https://github.com/NLipatov/TuiGo/issues/8)) ([84d901f](https://github.com/NLipatov/TuiGo/commit/84d901fb61db9237e2010e579c18ba7fbe4565b5))

## [0.1.1](https://github.com/NLipatov/TuiGo/compare/v0.1.0...v0.1.1) (2026-06-04)


### Features

* **device:** add IsModeChanged ([4919a8c](https://github.com/NLipatov/TuiGo/commit/4919a8c70eca703e47fb4744f445d22de742f033))


### Bug Fixes

* **session:** guard terminal restore lifecycle ([b9a3a6d](https://github.com/NLipatov/TuiGo/commit/b9a3a6d4df9cd7be9419a9dab00137055a71c01d))

## 0.1.0 (2026-06-01)


### Features

* add ansi escape sequence is command method ([30fb01b](https://github.com/NLipatov/TuiGo/commit/30fb01b11fbe226f2bb7a6d9747738b00c6e701a))
* add ansi escape sequences ([84972ac](https://github.com/NLipatov/TuiGo/commit/84972ac291f75f01425aebcbbbf4a9f1177dd698))
* add CSI to ansi escape sequences ([e142568](https://github.com/NLipatov/TuiGo/commit/e142568b08ff03e5cd4b731925877ab3cf150330))
* add domain model ([bfc5c3d](https://github.com/NLipatov/TuiGo/commit/bfc5c3db1a4c8ad5bc3953dc08b751185638ff03))
* add frame.RowAt method ([ac723d8](https://github.com/NLipatov/TuiGo/commit/ac723d8c2a66aa3ee8e135998d02af944f11c9bb))
* add hello example ([4311dbf](https://github.com/NLipatov/TuiGo/commit/4311dbf07b77e0864b697ef03d9e9a37b11970af))
* add new Cell methods: Foreground, Background, Symbol ([09a5f76](https://github.com/NLipatov/TuiGo/commit/09a5f76b368592382aa6895608072d35edbbbf07))
* add new escape sequences ([b722253](https://github.com/NLipatov/TuiGo/commit/b722253f299217f73d96561c619b45700a3ed687))
* add render write err full rerender test ([37746d2](https://github.com/NLipatov/TuiGo/commit/37746d237514cb00319f7dc905c9a9bf66ff8479))
* add screen type ([3249059](https://github.com/NLipatov/TuiGo/commit/3249059dcd6758a12fad1bffc128c0f28c1d4688))
* add session input event handling ([eacdccf](https://github.com/NLipatov/TuiGo/commit/eacdccf81d958b47c07602bedb45b422a6cfa35f))
* add size method to session ([0a6e4fc](https://github.com/NLipatov/TuiGo/commit/0a6e4fcacbf6f958df215e8cf85df15449722eba))
* add style diff ([8a797ff](https://github.com/NLipatov/TuiGo/commit/8a797ffa7755a56dd5b6276cc1982d10501fa47c))
* add terminal device ([0ce73a9](https://github.com/NLipatov/TuiGo/commit/0ce73a9a5341b01fc3e229171db3d6a873e7a038))
* add terminal input listener ([f2d1ee2](https://github.com/NLipatov/TuiGo/commit/f2d1ee24367c40827c4f420ac3ac7f6961295511))
* add terminal session event stream ([e4ab040](https://github.com/NLipatov/TuiGo/commit/e4ab0405f510015ee7c7c26f4f3267f082aeb70b))
* add unix and windows resize event listener implementations ([4246e2b](https://github.com/NLipatov/TuiGo/commit/4246e2bad5ff4e7fa1bd11341f952d349db93472))
* add zero-allocation cell rendering ([8b16b39](https://github.com/NLipatov/TuiGo/commit/8b16b39e4da5ab63f7bbb94d406bad14fcc038bc))
* batch renderer writes ([2cf285e](https://github.com/NLipatov/TuiGo/commit/2cf285e85b0debb8ec6e73d1c57353ac9787a458))
* batch renderer writes ([6f1860d](https://github.com/NLipatov/TuiGo/commit/6f1860d38520af556c610de00ccac8c0224f1f03))
* color.String() ([77983ab](https://github.com/NLipatov/TuiGo/commit/77983abf075d5d4c0f161ec54317867156feed27))
* emit listener error events ([8b78ecf](https://github.com/NLipatov/TuiGo/commit/8b78ecf0e36e315943689ca442580d6ff72e4790))
* go mod ([5106ba1](https://github.com/NLipatov/TuiGo/commit/5106ba19d10910abce082d59a7f855aa5fdf6924))
* grid.CellAt, gridSetCellAt ([818b500](https://github.com/NLipatov/TuiGo/commit/818b5002d1355c5038efaa19228029a1899774d9))
* implement input listener ([e1be5d7](https://github.com/NLipatov/TuiGo/commit/e1be5d77b08523e64acbf4e2ee0729eeb39ff620))
* implement row-based rendering ([904e9d4](https://github.com/NLipatov/TuiGo/commit/904e9d4b498bbc5ef7418f6b1d16f296169d483c))
* main entrypoint ([9f49cf1](https://github.com/NLipatov/TuiGo/commit/9f49cf10d28eeacb2eed0f3fabf8bf8e0229f813))
* track renderer frame state ([55d4184](https://github.com/NLipatov/TuiGo/commit/55d4184c301fd2a91bbf4edcf5dcb179c9a9dbdf))
* wire render into session ([ce78b04](https://github.com/NLipatov/TuiGo/commit/ce78b04b63915825c7eeb7f1f541f1cb188bd2af))


### Bug Fixes

* add style reset to restoreTerminal ([01f2f77](https://github.com/NLipatov/TuiGo/commit/01f2f779e869dce881e1f4125827408a6f671c12))
* fix renderer style ([054aa91](https://github.com/NLipatov/TuiGo/commit/054aa9166915d46dd5f4c3bdd362e88893cfb3e6))
