# -*- mode: ruby -*-
# vi: set ft=ruby :

# Require the reboot plugin.

$install_upstart = <<SCRIPT
  apt-get update
  apt-get install -y upstart-sysv
  update-initramfs -u
SCRIPT

$bootstrap = <<SCRIPT
  apt-get purge -y systemd 
  apt-get autoremove
  apt-get install -y curl git-core vim
  apt-get install -y fabric
  apt-get install -y bridge-utils
  apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
  apt-add-repository 'deb https://apt.dockerproject.org/repo ubuntu-xenial main'
  apt-get update
  apt-cache policy docker-engine
  apt-get install -y docker-engine
  apt-get install util-linux
  usermod -a -G docker ubuntu
  docker pull alpine:3.3
SCRIPT

$postinstall = <<SCRIPT
  iptables --policy FORWARD ACCEPT
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/xenial64"
  node_vm_name = "kjtest"
  #use common public key
  config.ssh.insert_key = false
  config.ssh.private_key_path = File.expand_path('~/.ssh/id_rsa')
  config.ssh.forward_agent = true

  config.vm.define node_vm_name do |node|
    node.vm.hostname = node_vm_name

      node.vm.provision "shell" do |s|
        ssh_pub_key = File.readlines("#{Dir.home}/.ssh/id_rsa.pub").first.strip
        s.inline = <<-SHELL
          echo #{ssh_pub_key} >> /home/ubuntu/.ssh/authorized_keys
          echo #{ssh_pub_key} >> /root/.ssh/authorized_keys
        SHELL
      end

    #apt-get
    node.vm.provision "install_upstart", type: "shell", run: "once", privileged: true, inline: $install_upstart  
    # comment by Jian
    #node.vm.provision :unix_reboot
    node.vm.provision "bootstrap", type: "shell", run: "once", privileged: true, inline: $bootstrap
    node.vm.provision "postinstall", type: "shell", run: "once", privileged: true, inline: $postinstall
  end
end