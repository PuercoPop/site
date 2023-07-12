# How to package a Rails app in Nix

- Use gitlab-rails as an example
- Mastodon
- zammad

after building the rubyEnv we need to create a derivation

What is the installPhase?

```nix
stdenv.mkDerivation {
    pname = ;
    version =;
    src = ./.;
    buildInputs = [
      rubyEnv
      rubyEnv.wrappedRuby
      rubyEnv.bundler
    ]

    RAILS_ENV = "production";

    installPhase = ''
    '';
}
```
