/**
 * Created by gsgritta on 30/4/17.
 */

var gulp = require('gulp');
var uglify = require('gulp-uglify');
var cleanCSS = require('gulp-clean-css');
var clean = require('gulp-clean');
var rename = require('gulp-rename');

gulp.task('default', ['build']);

gulp.task('build', ['minify-js', 'minify-css']);

gulp.task('clean', function(){
    return gulp.src(['assets/js/*.min.js', 'assets/style/*.min.css'], {read: false})
        .pipe(clean());
});

gulp.task('minify-css', ['clean'], function(){
    return gulp.src(['assets/style/*.css'])
        .pipe(cleanCSS())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('assets/style/'));
});

gulp.task('minify-js', ['clean'], function(){
    return gulp.src(['assets/js/*.js'])
        .pipe(uglify())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('assets/js/'));
});
