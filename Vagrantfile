# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "fs" do |v|
    v.vm.provider "docker" do |d|
      d.build_dir = "env"
      d.name = "fs"
      # d.cmd     = ["/sbin/my_init", "--enable-insecure-key"]
      # d.cmd = ["freeswitch"]
      d.has_ssh = true
    end
    v.ssh.port = 22
    v.ssh.username = "root"
    v.ssh.private_key_path = "env/phusion.key"
    v.vm.provision "shell", inline: "echo hello"
  end

  config.vm.define "phusion" do |v|
    v.vm.provider "docker" do |d|
      d.cmd = ["/sbin/my_init", "--enable-insecure-key"]
      d.image = "phusion/baseimage"
      d.name = "dockerizedvm"
      d.has_ssh = true
      #d.force_host_vm = true
    end
    v.ssh.port = 22
    v.ssh.username = "root"
    v.ssh.private_key_path = "phusion.key"
    v.vm.provision "shell", inline: "echo hello"
    #v.vm.synced_folder "./keys", "/vagrant"
  end
  
end