## 3.1.0
- Support Apple Silicon chip (M1)
- Improve installation/update scripts
- Bug fixes
- Update package dependency 
- Move from CircleCI to github action

## 3.0.15
- Support deploying a new darknode using an existing config file.

## 3.0.14
- Fix bug when fetching latest darknode-release during darknode installation
- Remove `jq` from darknode installation

## 3.0.13
- Fix pagination issue when fetching latest darknode releases from github.

## 3.0.12
- Fix connection issue when withdrawing gas 

## 3.0.11
- Use suggest gas price for withdrawing ETH.

## 3.0.10
- Improve `update` command to support suffix tags #61 
- Show darknodes versions when from the `list` command.

## 3.0.9
- Update bootstrap node's multi-address

## 3.0.8
- Update the terraform version when installing the CLI. 
- Update protocol contract addresses 
- Support mainnet deployment

## 3.0.7
- Check existence of the Darknode when running commands which darknode name as parameter #55 
- `up` command gets latest release from github, instead of ipfs. 
- Remove auto-updater and modify the `update` command #56 
- Make sure the CI script does not overwrite the previous release

## 3.0.6
Minor fix and improvement #52 
- Fix concurrent writes issue when running `darknode list`
- Hide the error message when checking version against github.
- Ignore the `ECDSADistKeyShare` field when reading the config.

## 3.0.5 
Issue #49 
- Show messages when starting executing the script.

Issue #51 
- Update config format to the latest version in darknode.
- Remove `--config` flag form the `update` command to avoid conflicts in the config file. 
- Show messages to tell if nodes have successfully been started/stopped/restarted.
- Minor tweak with the darknode installation script.
- Fix incorrect bootstrap addresses on testnet and devnet.

## 3.0.4
- Add Google Cloud Platform support
Thanks to @Pega88 for contributing to this release

## 3.0.3
- Fix an issue which could cause fail deployment on digital ocean. 
- Successful deployment on Windows(WSL) will automatically open the registering link.
- Hide the red error message when cannot find any command to open the registering link. 
 
## 3.0.2
- set `DisablePeerDiscovery` to true when initializing a new darknode config. 
- manually install the metrics agent for digital ocean to enable detail monitoring

## 3.0.1
- Build binary for `arm` architecture and include it in the release.
- Update bindings for darknode registry contract to the latest version.
- Not showing any error when failed to open registering url in a browser.
- Update prompt message to show the correct command when a new release is detected.

## 3.0.0

New CLI for Chaosnet darknodes deployment. 

> Since the config format for darknode is changed and we move from terraform 0.11.x to 0.12.x. New cli will not be backwards compatible with old darknodes. Users need to deregister and destroy old nodes first before updating to this version.

Changelog:
- Windows support with WSL
- Use updated `terraform` version `v0.12.12`
- Release check before running any commands, it will warn users if they are using an old version.
- Support `start/stop/restart` commands to start/stop/restart the service on darknode.
- Add `exec` command which allow users to run script/file on remote instance.
- Add `register` command which shows the link to register a particular node and try opening it using the default browser. 
- Show provider information when running `darknode list`
- `resize` command will not increase the disk size for digital ocean. (This allows users to downgrade their droplets after upgrading the plan) 
- Detail monitoring is enabled for all providers. 
- Code has been refactored to be self-contained 
