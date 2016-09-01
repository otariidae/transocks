# Change Log

All notable changes to this project will be documented in this file.

## [Unreleased]

## [1.0.0] - 2016-09-01
### Added
- transocks now adopts [github.com/cybozu-go/cmd][cmd] framework.  
  As a result, it implements [the common spec][spec] including graceful restart.

### Changed
- The default configuration file path is now `/etc/transocks.toml`.
- "listen" config option becomes optional.  Default is "localhost:1081".
- Configuration items for logging is changed.

[cmd]: https://github.com/cybozu-go/cmd
[spec]: https://github.com/cybozu-go/cmd/blob/master/README.md#specifications
[Unreleased]: https://github.com/cybozu-go/transocks/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/cybozu-go/transocks/compare/v0.1...v1.0.0
