/**
 * Project: flanders
 * Module:
 * User: tyson
 * Date: 10/27/14
 * Time: 4:43 PM
 */
var jshint = require('gulp-jshint');

gulp.task('lint', function() {
  return gulp.src('compiled/**/*.js')
    .pipe(jshint())
    .pipe(jshint.reporter('jshint-stylish'));
});
