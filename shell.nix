let
  pkgs = import <nixpkgs> {
    system = "x86_64-linux";
  };

in
pkgs.mkShell {
  buildInputs = [
    pkgs.bashInteractive
    pkgs.etcd
    pkgs.libvirt
    pkgs.virt-manager
  ];

  shellHook = ''
    echo "Entering cthul dev environment (using shell.nix with system nixpkgs)..."

    sudo etcd --name test \
      --listen-client-urls http://127.0.0.1:2379 \
      --advertise-client-urls http://127.0.0.1:2379 \
      --listen-peer-urls http://127.0.0.1:2380 \
      --initial-advertise-peer-urls http://127.0.0.1:2380 \
      --initial-cluster test=http://127.0.0.1:2380 \
      --initial-cluster-state new \
      --initial-cluster-token test \
      --data-dir /var/lib/cthul/etcd/ &
    
    echo "Background services initiated. Check their status if needed."
  '';
}
