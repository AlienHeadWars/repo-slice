# file: .idx/dev.nix
# To learn more about how to use Nix to configure your environment
# see: https://firebase.google.com/docs/studio/customize-workspace
{ pkgs, ... }: {
  # Which nixpkgs channel to use.
  channel = "unstable"; # Use the unstable channel for newer packages

  # Use https://search.nixos.org/packages to find packages
  packages = [
    # Add the specific Go version and linter here
    pkgs.go_1_24
    pkgs.golangci-lint
  ];
  # Sets environment variables in the workspace
  env = {};
  idx = {
    # Search for the extensions you want on https://open-vsx.org/ and use "publisher.id"
    extensions = [];
    # Enable previews
    previews = {
      enable = true;
      previews = {};
    };

    # Workspace lifecycle hooks
    workspace = {
      # Runs when a workspace is first created
      onCreate = {};
      # Runs when the workspace is (re)started
      onStart = {};
    };
  };
}