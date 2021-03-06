# authplz

A simple Authentication and User Management microservice, designed to help build secure user-accessible services, and to avoid having to write another authentication and user management service (ever again).

This is heavily inspired by the way Github manage user accounts, two factor authentication, authorized devices etc., and is intended to provide common user management features required for a web application (or web application suite) to handle authentication of web, native and mobile applications, as well end devices, with the minimum possible complexity for developers.

Systems using this service will use OAuth grants with token introspection to validate user credentials so that users can create third party applications that utilise the same APIs.

This provides an alternative to hosted solutions such as [StormPath](https://stormpath.com/) and [AuthRocket](https://authrocket.com/) for companies that prefer (or require) self hosted identity providers. 
For a well supported self hosted alternative for Single Sign On (SSO) you may wish to investigate [gluu](https://www.gluu.org), as well as wikipedia's [List of SSO implementations](https://en.wikipedia.org/wiki/List_of_single_sign-on_implementations).
If you already have user management infrastructure, you may be interested in [coreos/dex](https://github.com/coreos/dex) as an OAuth extension.

If you would like to be involved with this project, please first read (and agree to abide by) the [Code of Conduct](https://github.com/authplz/authplz-core/blob/master/CONDUCT.md), then go ahead and join the chat on [Gitter](https://gitter.im/authplz/Lobby) or [open an issue](https://github.com/authplz/authplz-core/issues/new).

## Goals

- Developers shouldn't have to write any user management / authorization code
- Users should be able to manage their accounts & authorizations (and create third party apps using these)
- Admins should be able to manage user accounts and create integrations
- Users should be able to make informed security decisions about their account

Check out [design.md](design.md) for more.

## Status

WIP. Most components working.

[![GitHub tag](https://img.shields.io/github/tag/authplz/authplz-core.svg)](https://github.com/authplz/authplz-core)
[![Build Status](https://travis-ci.org/authplz/authplz-core.svg)](https://travis-ci.org/authplz/authplz-core/branches)
[![Documentation](https://img.shields.io/badge/docs-godoc-blue.svg)](https://godoc.org/github.com/authplz/authplz-core)
[![Chat](https://img.shields.io/gitter/room/gitterHQ/gitter.svg)](https://gitter.im/authplz/Lobby)

Check out the current build at [https://authplz.github.io](https://authplz.github.io).

## Usage

Frontend components and templates are now in a [authplz/authplz-ui](https://github.com/authplz/authplz-ui) project (and have been grossly neglected).
Configuration is now via a yaml [configuration file](authplz.yml) which supports explicit loading of environmental variables and explains all the required config options. Use `./authplz --help` to display options to specify the config file or environment prefix.

### Dependencies

- Golang (for building)
- Docker (for building/running docker images and the dev environment)
- Postgres (for backing data storage)

### Development

If you have contributor access to the repository, changes should be created in branches and pull requests opened to merge (as is enforced by the repository settings and Travis-CI). If you don't have access, please follow the normal fork/pull-request flow (though be aware that forking GO projects can be a little interesting due GOPATH).

All features must be implemented with tests to demonstrate the correct (and incorrect) behaviours of the feature.

1. `go get github.com/authplz/authplz-core.git` to fetch the core repo into your GOPATH
2. `cd $GOPATH/src/github.com/authplz/authplz-core` to switch to the repo
3. (For out of tree development) `git remote add upstream github.com/authplz/authplz-core.git` to add the root as an upstream
4. (For out of tree development) `git remote set-url github.com/YOURNAME/authplz-core.git` to set the master to your fork
5. `./gencert.sh` to generate self signed TLS certificates
6. `make build-env` and `make start-env` to build and run dependent services (eg. the database)
7. `git checkout -b "feature/my-new-feature"` to create a new branch
8. Do some work...
9. `make test` to run repository tests
10. `make run` to run the application if required
11. Once tests pass you can `git commit` and `git push` your changes

For frontend development it is useful to run a local AuthPlz instance that is then proxied by the create-react-app development runner.
This is automatic when calling `npm start` from the authplz-ui project.

### Running

1. `make install` to install dependencies
2. `./gencert.sh` to generate self signed TLS certificates
3. `make build-env` and `make start-env` to build and run dependencies
4. `./authplz` to launch the app

`./authplz --help` will list available configuration options.

## Features

- [X] Account creation
- [X] Account activation
- [X] User login
- [ ] User administration
  - [ ] Account Unlock / Password Reset
  - [ ] Account enable / disable
- [X] Account locking (and token + password based unlocking)
- [X] User logout
- [X] User password update
- [X] User Password reset
- [X] Email notifications
- [X] Audit / Event logging
- [X] 2FA token enrolment
  - [X] TOTP
  - [X] FIDO
  - [X] BACKUP
- [X] 2FA token validation
  - [X] TOTP
  - [X] FIDO
  - [X] BACKUP
- [X] 2FA token management
  - [X] TOTP
  - [X] FIDO
  - [X] BACKUP
- [-] OAuth2
  - [X] Authorization Code grant type
  - [X] Implicit grant type
  - [ ] User client management
  - [ ] User token management
- [X] ACLs (based on fosite heirachicle ie. `public.something.read`)
- [ ] Account linking (google, facebook, github)
- [ ] Plugin Support
  - [ ] IP based rate limiting
  - [ ] Webhooks
  - [ ] Distributed Synchronisation
- [X] Test Server
  - [X] Deployment to https://authplz.herokuapp.com
  - [X] Deployment of frontend assets

## Project Layout

Checkout [DESIGN.md](DESIGN.md) for design notes and API interaction flows.

- [cmd/authplz/main.go](cmd/authplz/main.go) contains the launcher for the AuthPlz server
- [lib/api](lib/api) contains internal and external API definitions
- [lib/app](lib/app) contains the overall application including configuration and wiring (as well as integration tests)
- [lib/appcontext](lib/appcontext) contains the base application context (shared across all API modules)
- [lib/controllers](lib/controllers) contains controllers that can be shared across API modules
  - [lib/datastore](lib/datastore) contains the data storage module and implements the interfaces required by other modules
  - [lib/token](lib/controllers/token) contains a token generator and validator
- [lib/modules](lib/modules) contains functional modules that can be bound into the system (including interface, controller and API)
  - [lib/core](lib/modules/core) contains the core login/logout/action endpoints that further modules are bound into. Checkout this module for information on what components / bindings are available.
  - [lib/user](lib/modules/user) contains the user account management module and API
  - [lib/2fa](lib/modules/2fa) contains 2fa implementations
  - [lib/user](lib/modules/audit) contains the account action / auditing API
  - [lib/user](lib/modules/oauth) contains oauth and OpenID functionality
- [lib/templates](lib/templates) contains default template files used by components (ie. mailer)
- [lib/test](lib/test) contains test helpers (and maybe one day integration tests)

Modules are self-binding and should define interfaces required to function rather than including any (non api or appcontext) other modules.

Each module should define the interfaces required, a controller for interaction / data processing, and an API if required by the module. For an example, checkout [lib/modules/2fa/u2f](lib/modules/2fa/u2f).


------

If you have any questions, comments, or suggestions, feel free to contact us (uhh, me) on gitter or to open an issue or a pull request.
