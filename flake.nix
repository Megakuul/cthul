{
  description = "Cthul system flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.05";
  };

  outputs = { self, nixpkgs, ... }@inputs: {
    nixosModules.default = import "${self}/build/package/default.nix";
  };
}
