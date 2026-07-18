{
  description = "Bluescreen – Ebiten dev shell";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Alles, was Ebiten/GLFW zur Laufzeit dynamisch nachlädt
        runtimeLibs = with pkgs; [
          libGL
          libX11
          libXrandr
          libXcursor
          libXinerama
          libXi
          libXxf86vm
          alsa-lib
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [ go gopls pkg-config ];
          buildInputs = runtimeLibs;

          # CGO findet die Header über buildInputs, aber der Linker/Loader
          # braucht den Pfad zur Laufzeit nochmal explizit
          LD_LIBRARY_PATH = pkgs.lib.makeLibraryPath runtimeLibs;
        };
      });
}
