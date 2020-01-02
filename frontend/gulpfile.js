"use strict";

const gulp = require("gulp");
const sass = require("gulp-sass");
const sourcemaps = require("gulp-sourcemaps");
const cleanCss = require("gulp-clean-css");
const htmlMinifier = require("gulp-html-minifier");
const uglify = require("gulp-uglify");
const concat = require("gulp-concat");
const rename = require("gulp-rename");
const eslint = require("gulp-eslint");

const develPath = "build/devel/";
const prodPath = "build/prod/";
const scssSrc = "./sass/*.scss";
const cssDir = "css/";
const fontsDir = "fonts/";
const fontsGlob = fontsDir + "**/*";
const imagesDir = "images/";
const imagesGlob = imagesDir + "**/*";
const jsDir = "js/";
const jsGlob = jsDir + "*.js";
const htmlGlob = "./*.html";

const jsCompileMap = {
  "js/jquery.js": ["node_modules/jquery/dist/jquery.min.js"],
  "js/vue.js": ["node_modules/vue/dist/vue.min.js"],
  "js/highlight.js": ["node_modules/highlightjs/highlight.pack.min.js"],
  "js/chartist.js": ["node_modules/chartist/dist/chartist.min.js"],
  "js/login.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/auth-common.js",
    "js/login.js"
  ],
  "js/forgot.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/forgot.js"
  ],
  "js/reset.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/reset.js"
  ],
  "js/signup.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/auth-common.js",
    "js/signup.js"
  ],
  "js/dashboard.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/errors.js",
    "js/self.js",
    "js/dashboard.js",
    "js/dashboard-setting.js",
    "js/dashboard-domain.js",
    "js/dashboard-installation.js",
    "js/dashboard-general.js",
    "js/dashboard-moderation.js",
    "js/dashboard-statistics.js",
    "js/dashboard-import.js",
    "js/dashboard-danger.js",
    "js/dashboard-export.js",
  ],
  "js/settings.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/errors.js",
    "js/self.js",
    "js/settings.js"
  ],
  "js/logout.js": [
    "js/constants.js",
    "js/utils.js",
    "js/logout.js"
  ],
  "js/commento.js": ["js/commento.js"],
  "js/count.js": ["js/count.js"],
  "js/unsubscribe.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/unsubscribe.js",
  ],
  "js/profile.js": [
    "js/constants.js",
    "js/utils.js",
    "js/http.js",
    "js/profile.js",
  ],
};

gulp.task("scss-devel", function (done) {
  let res = gulp.src(scssSrc)
    .pipe(sourcemaps.init())
    .pipe(sass({outputStyle: "expanded"}).on("error", sass.logError))
    .pipe(sourcemaps.write())
    .pipe(gulp.dest(develPath + cssDir));
  done();
  return res;
});

gulp.task("scss-prod", function (done) {
  let res = gulp.src(scssSrc)
    .pipe(sass({outputStyle: "compressed"}).on("error", sass.logError))
    .pipe(cleanCss({compatibility: "ie8", level: 2}))
    .pipe(gulp.dest(prodPath + cssDir));
  done();
  return res;
});

gulp.task("html-devel", function (done) {
  gulp.src([htmlGlob]).pipe(gulp.dest(develPath));
  done();
});

gulp.task("html-prod", function (done) {
  gulp.src(htmlGlob)
    .pipe(htmlMinifier({collapseWhitespace: true, removeComments: true}))
    .pipe(gulp.dest(prodPath))
  done();
});

gulp.task("fonts-devel", function (done) {
  gulp.src([fontsGlob]).pipe(gulp.dest(develPath + fontsDir));
  done();
});

gulp.task("fonts-prod", function (done) {
  gulp.src([fontsGlob]).pipe(gulp.dest(prodPath + fontsDir));
  done();
});

gulp.task("images-devel", function (done) {
  gulp.src([imagesGlob]).pipe(gulp.dest(develPath + imagesDir));
  done();
});

gulp.task("images-prod", function (done) {
  gulp.src([imagesGlob]).pipe(gulp.dest(prodPath + imagesDir));
  done();
});

gulp.task("js-devel", function (done) {
  for (let outputFile in jsCompileMap) {
    gulp.src(jsCompileMap[outputFile])
      .pipe(sourcemaps.init())
      .pipe(concat(outputFile))
      .pipe(rename(outputFile))
      .pipe(sourcemaps.write())
      .pipe(gulp.dest(develPath))
  }
  done();
});

gulp.task("js-prod", function (done) {
  for (let outputFile in jsCompileMap) {
    gulp.src(jsCompileMap[outputFile])
      .pipe(concat(outputFile))
      .pipe(rename(outputFile))
      .pipe(uglify())
      .pipe(gulp.dest(prodPath))
  }
  done();
});

gulp.task("lint", function (done) {
  let res = gulp.src(jsGlob)
    .pipe(eslint())
    .pipe(eslint.failAfterError());
  done();
  return res;
});

gulp.task("devel", gulp.parallel("scss-devel", "html-devel", "fonts-devel", "images-devel", "lint", "js-devel"));
gulp.task("prod", gulp.parallel("scss-prod", "html-prod", "fonts-prod", "images-prod", "lint", "js-prod"));
