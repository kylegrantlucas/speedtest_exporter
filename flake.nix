{
  description = "Prometheus exporter for the Ookla Speedtest CLI";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "aarch64-darwin"
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs systems (
          system:
          f (
            import nixpkgs {
              inherit system;
              # ookla-speedtest is unfree.
              config.allowUnfreePredicate = pkg: nixpkgs.lib.getName pkg == "ookla-speedtest";
            }
          )
        );
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "speedtest_exporter";
          version = self.shortRev or self.dirtyShortRev or "dev";
          src = ./.;
          vendorHash = "sha256-nIJE8A+/nrazsxymktYkzuN+uFnixRkZpfq7h1PiJHM=";
          env.CGO_ENABLED = 0;
          ldflags = [
            "-s"
            "-w"
          ];
          nativeBuildInputs = [ pkgs.makeWrapper ];
          postInstall = ''
            wrapProgram $out/bin/speedtest_exporter \
              --prefix PATH : ${pkgs.lib.makeBinPath [ pkgs.ookla-speedtest ]}
          '';
          meta.mainProgram = "speedtest_exporter";
        };
      });

      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            gopls
            golangci-lint
            ookla-speedtest
          ];
        };
      });
    };
}
