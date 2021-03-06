# Design

Design notes and questions for AuthPlz.


## Overview

- All endpoints should return 400 bad request if required parameters or fields in the request body are not issued
- Authentication endpoints will return 200 success/201 partial or 403 unauthorized
- Internal errors will result in a 401 internal error with no (or a generic) error message to avoid leaking internal
- Other endpoints will return JSON formatted API messages

### Modules

Modules consist of a set of interfaces, defining the dependencies of the module, a controller that uses these interfaces to implement the functionality of the module, and an api that wraps the controller in HTTP endpoints.

All data structures returned from controllers should be safe for API use (ie. no internal structures may be returned, explicitly wrap / translate everything).


## Flows

### Account Creation

1. post username, email, password to /api/users/create
2. server sends account activation email
3. execute activation link
4. server sends account activated email


### Basic Login

1. post email, password to /api/login
2. server responds with 200 success, 201 partial (2fa) or 403 unauthorized


### Account Unlock

1. post email, password to /api/login
2. server responds with 403 unauthorized, sends unlock email to registered address, caches credentials
3. get /api/flash returns "Account locked"
4. user clicks unlock link to /api/action?token=TOKEN
5. server redirects to login page
6. post email, password to /api/login (unless credentials are already cached)
7. server executes unlock token

Seems like this could be more efficient / remove the need for the second login if the user clicks the unlock link with a partially formed session (valid username and password).


### Password Change 

1. user logs in as above
2. user submits old, new passwords to /api/users/account
3. server validates, responds with 200 success or 400 bad request


### U2F enrolment

1. user logs in as above
2. user submits token name to /api/u2f/register
3. server responds with registration challenge
4. browser executes challenge, posts response
5. server validates registration response
6. server responds with 200 success or 401 error


### U2F Login

1. post email, password to /api/login
2. server responds with 201 partial (2fa) and available factors object ({u2f: true})
3. browser fetches challenge from /api/u2f/authenticate
4. browser executes challenge, posts response
5. server responds with 200 success or 403 unauthorized

### TOTP enrolment

1. user logs in as above
2. user submits token name to /api/totp/enrol
3. server responds with registration challenge (string and image)
4. user loads totp onto device, posts a valid code
5. server validates registration response
6. server responds with 200 success or 401 error


### TOTP Login

1. post email, password to /api/login
2. server responds with 201 partial (2fa) and available factors object ({totp: true})
3. user gets code from totp app
4. browser posts code to /api/totp/authenticate
5. server responds with 200 success or 403 unauthorized

### Password Reset

1. post email account to /api/recovery
2. server sends recovery token to user email
3. token submitted to /api/recovery (could be /api/token, but different process required so easier to split)
4. if 2fa, require 2fa to validate recovery session. If lost, sms or recovery codes.
5. user submits new password to /api/reset
6. server responds 200 success or 400 bad request
7. server sends alert email to user

This requires that all stages be undertaken from the same session. Backup codes are treated just another 2fa provider.


### OAuth Clients

A variety of clients can be enrolled based on user account priviledges. Admins can enrol all OAuth client types, users can enrol Client Credential (for end devices) and Implicit (no secret storage) client types.
Allowed authorizations for each account type can be set in the configuration file.

#### Authorisation Code (Explicit) Grant
For trusted services, created by administrators, available to all users.

#### Authorisation Code (Implicit) Grant
For services that do not have secret storage, created by and available to individual users.

#### Client Credentials Grant
For end devices, created by and available to individual users.

#### Introspection
Explicit grants can be provided with the "introspection" scope, allowing introspection of other tokens using these credentials.
This allows trusted services to evaluate the validity of credentials for broker-like behaviour.


#### Refresh Token Grant

Allows tokens to be refreshed / reissued. Available with both Authorization Code grant types.

## Questions

- How can you enrol / remove tokens, what is required?
- How do plugins require further login validation (ie. "this doesn't look right, click token in email to validate")?
- How do we run multiple OAuth schemes for different clients?
  - Guess user interaction is going to be important here as to what is granted

## Discussions

What if instead of imposing a security level on users, we informed them and let them pick?
Users could then be given a security score on their account dashboard to gamify improving it.
For example:
- You only have password set, password resets and account recovery will currently require only your email address, register a phone number or 2fa token to improve this
- Good work registering 2fa! Password resets will now require this 2fa token. For account recovery purposes you must now either register a phone number or create recovery codes

Other ideas:

- Testing that users have recovery keys etc to keep people in practice. Users could be (optionally) prompted to enter a backup key name at a given interval, on failure given the option of regenerating backup keys.

