{ config, pkgs, lib, ... }:

{
  options.cthul = {
    enable = lib.mkEnableOption "enable cthul system";
    node = lib.mkOption {
      type = lib.types.str;
      default = "node001.cluster.cthul.io";
      description = "id used to identify this node";
    };

    rune = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable rune service";
      };
    };
    wave = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable wave service";
      };
    };
    granit = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable granit service";
      };
    };
    proton = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable proton service";
      };
    };
    etcd = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable etcd database";
      };
    };
    libvirt = {
      enable = lib.mkOption {
        type = lib.types.bool;
        default = config.cthul.enable;
        description = "enable libvirt service";
      };
    };
  };

  config = {
    services.etcd = lib.mkIf config.cthul.etcd.enable {
      enable = true;
      name = config.cthul.node;
      dataDir = "/var/lib/cthul/etcd";
      initialCluster = ["${config.cthul.node}=http://127.0.0.1:2380"];
      initialAdvertisePeerUrls = ["http://127.0.0.1:2380"];
      listenPeerUrls = ["http://127.0.0.1:2380"];
      advertiseClientUrls = ["http://127.0.0.1:2379"];
      listenClientUrls = ["http://127.0.0.1:2379"];
    };

    virtualisation.libvirtd = lib.mkIf config.cthul.libvirt.enable {
      enable = true;
    };
  }; 
}
