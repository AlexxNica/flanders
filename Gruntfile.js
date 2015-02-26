module.exports = function(grunt) {
  grunt.initConfig({
    exec: {
      build_go_app: {
        command: 'GOOS=linux go build -o bin/flanders main/main.go',
        stdout: true,
        stderr: true
      },
      build_gui: {
        command: 'grunt build',
        cwd: __dirname + "/web",
      }
    },
    copy: {
      gui: {
        files: [
          // includes files within path and its sub-directories
          {expand: true, cwd:'web/dist/', src: ['**'], dest: 'bin/www'},
        ],
      },
    },
    gitcheckout: {
      develop: {
        options: {
          branch: 'develop'
        }
      },
      master: {
        options: {
          branch: 'master'
        }
      }
    },
    secret: grunt.file.readJSON('secret.json'),
    environments: {
      production: {
        options: {
          host: '<%= secret.production.host %>',
          username: '<%= secret.production.username %>',
          password: '<%= secret.production.password %>',
          deploy_path: '/opt/flanders',
          local_path: 'bin',
          current_symlink: 'current',
          before_deploy: '  echo <%= secret.production.password %> | sudo -S stop flanders',
          after_deploy: '  echo <%= secret.production.password %> | sudo -S start flanders',
        }
      },
    }
  });

  grunt.loadNpmTasks('grunt-exec');
  grunt.loadNpmTasks('grunt-ssh-deploy');
  grunt.loadNpmTasks('grunt-git');
  grunt.loadNpmTasks('grunt-contrib-copy');

  grunt.registerTask('default', ['exec:build_go_app', 'exec:build_gui', 'copy:gui']);
  grunt.registerTask('deploy', 'This deploys app and gui to production', function() {
    grunt.task.run('exec:build_go_app');
    grunt.task.run('exec:build_gui');
    grunt.task.run('copy:gui');
    grunt.task.run('ssh_deploy:production');
  });

  grunt.registerTask('deployApp', 'This builds and deploys the Go app to production', function() {
    grunt.task.run('exec:build_go_app');
    grunt.task.run('ssh_deploy:production');
  })

  grunt.registerTask('deployGui', 'This builds and deploys the Go app to production', function() {
    grunt.task.run('exec:build_gui');
    grunt.task.run('copy:gui');
    grunt.task.run('ssh_deploy:production');
  })

  //grunt.registerTask('default', ['jshint', 'qunit', 'concat', 'uglify']);

};