## Design Decisions and Question

1. Should a user be able to create a `Client` struct directly or should we
  require them to invoke `NewClient`?
2. Does it make sense to introduce a 3rd-party testing framework and what are
  the tradeoffs?
3. Should `Authenticate` be separate from the other API calls? Should we make an
  optional **auto-authenticate** feature, and how does this change the
  `NewClient` method signature?
4. Review exported/unexported types
5. Refactor all of the `log.Fatal` calls inside of `Authenticate`
6. Consider the ramifications of using `var defaultBaseURL`.
