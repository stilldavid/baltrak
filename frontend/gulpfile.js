var gulp  = require('gulp'),
    gutil = require('gulp-util')

    jshint     = require('gulp-jshint'),
    sass       = require('gulp-sass'),
    concat     = require('gulp-concat'),
    //sourcemaps = require('gulp-sourcemaps');

    input  = {
      'sass': 'source/scss/**/*.scss',
      'javascript': [
        'bower_components/jquery/dist/jquery.js',
        'bower_components/moment/moment.js',
        'bower_components/leaflet/dist/leaflet.js',
        'bower_components/Leaflet.Coordinates/dist/Leaflet.Coordinates-0.1.5.src.js',
        'source/js/**/*.js',
      ],
      'html': 'source/**/*.html',
      'vendorjs': 'public/assets/js/vendor/**/*.js'
    },

    output = {
      'stylesheets': 'public/assets/css',
      'javascript': 'public/assets/js'
    };

/* run the watch task when gulp is called without arguments */
gulp.task('default', ['watch']);

/* run javascript through jshint */
gulp.task('jshint', function() {
  return gulp.src(input.javascript)
    .pipe(jshint())
    .pipe(jshint.reporter('jshint-stylish'));
});

/* compile scss files */
gulp.task('build-css', function() {
  return gulp.src(input.sass)
    .pipe(sass().on('error', sass.logError))
    .pipe(gulp.dest(output.stylesheets));
});

/* concat javascript files, minify if --type production */
gulp.task('build-js', function() {
  return gulp.src(input.javascript)
    .pipe(concat('bundle.js'))
    //only uglify if gulp is ran with '--type production'
    .pipe(gutil.env.type === 'production' ? uglify() : gutil.noop()) 
    .pipe(gulp.dest(output.javascript));
});

gulp.task('copy', function() {
  // copy any html files in source/ to public/
  gulp.src('source/*.html').pipe(gulp.dest('public'));

  gulp.src('bower_components/leaflet/dist/images/*.png').pipe(gulp.dest('public/assets/img'));
});

/* Watch these files for changes and run the task on update */
gulp.task('watch', function() {
  gulp.watch(input.javascript, ['jshint', 'build-js']);
  gulp.watch(input.html, ['copy']);
  gulp.watch(input.sass, ['build-css']);
});
