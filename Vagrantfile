# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

  # Convenience VM for some easy testing of Flanders
  config.vm.define "fs" do |v|
    # This box is Ubuntu 14.04 LTS with Salt pre-installed
    v.vm.box = "inflection/ubuntu-1404-salt"

    v.vm.provider "virtualbox" do |vb|
      vb.customize ["modifyvm", :id, "--memory", "2048"]
      vb.name = "freeswitch"
    end

    # Create a private network, which allows host-only access to the machine
    # using a specific IP.
    v.vm.network "private_network", ip: "12.0.0.2"

    # Bridged networks make the machine appear as another physical device on
    # your network.
    # v.vm.network "public_network"

    v.vm.synced_folder "salt/roots/salt/", "/srv/salt/"
    v.vm.synced_folder "./", "/opt/go/src/lab.getweave.com/weave/flanders"
    v.vm.provision :salt do |salt|
      salt.minion_config = "salt/minion"
      salt.pillar({
        "freeswitch" => {
          "version" => "v1.4.6",
          "use_debs" => false,
        }
      })
      salt.log_level = 'all'
      salt.run_highstate = true
    end 
    
  end
  
end