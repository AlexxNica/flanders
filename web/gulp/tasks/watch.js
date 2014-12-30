/*jshint indent: 2, node: true, nomen: true, browser: true*/
/*global gulp */
var browserSync = require('browser-sync');

gulp.task('watch', ['setWatch', 'browserSync'], function () {
  gulp.watch('app/scripts/**/*.jsx', ['react', browserSync.reload]);
  gulp.watch('app/scripts/**/*.js', ['copy.scripts', browserSync.reload]);
  gulp.watch('app/less/**/*.less',['less', browserSync.reload]);
  gulp.watch('app/images/**', ['images', browserSync.reload]);
  gulp.watch('app/**/*.html', ['copy.html', browserSync.reload]);
  // Note: The browserify task handles js recompiling with watchify
});
