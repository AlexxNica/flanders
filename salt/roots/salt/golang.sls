go-package:
  pkg.installed:
    - pkgs:
      - golang

gopath:
  # File.append searches the file for your text before it appends so it won't append multiple times
  file.append:
    - name: /home/vagrant/.profile
    - text: export GOPATH=/opt/go