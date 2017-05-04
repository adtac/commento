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
    return gulp.src(['assets/**/*.min.js', 'assets/**/*.min.css'], {read: false})
        .pipe(clean());
});

gulp.task('minify-css', ['clean'], function(){
    return gulp.src(['assets/**/*.css'])
        .pipe(cleanCSS())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('assets/'));
});

gulp.task('minify-js', ['clean'], function(){
    return gulp.src(['assets/**/*.js'])
        .pipe(uglify())
        .pipe(rename({suffix: '.min'}))
        .pipe(gulp.dest('assets/'));
});