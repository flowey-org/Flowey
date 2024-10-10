{
  inputs = {
    nixpkgs.url = "github:paveloom/nixpkgs/system";
  };

  outputs =
    { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        name = "flowey-shell";

        nativeBuildInputs = with pkgs; [
          bashInteractive
          nil
          nixfmt-rfc-style

          nodejs_latest
          typescript-language-server
          vscode-langservers-extracted

          go_1_23
          (gopls.override {
            buildGoModule = pkgs.buildGo123Module;
          })

          dockerfile-language-server-nodejs
          hadolint
          yamlfmt
          yamllint
        ];
      };
    };
}
