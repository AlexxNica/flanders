mongo-deb:
  pkgrepo.managed:
    - humanname: MongoDB PPA
    - name: deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen
    - dist: dist
    - file: /etc/apt/sources.list.d/mongodb.list
    - keyserver: hkp://keyserver.ubuntu.com:80
    - keyid: 7F0CEB10
    - require_in:
      - pkg: mongodb-org

  pkg.latest:
    - name: mongodb-org
    - refresh: True