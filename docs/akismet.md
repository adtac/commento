# Akismet Integration

Akismet is a leading spam filtering service that's used by millions of websites. For non-commercial websites, you can obtain a free API key (you can also pay as much as you want). Akismet integration is completely optional in Commento; if you do not want to use their service, comment will not be filtered. However, if you do choose to enable Akismet integration, be advised that *all* comments will be sent to Akismet's servers for spam detection.

To do this, go to [their website](https://akismet.com) and signup for an account. Choose your [appropriate plan](https://akismet.com/account/upgrade), and then go to your dashboard. Here, you'll find your 12-digit hexadecimal API key. To make Commento use this key, either set an environment variable though `export AKISMET_KEY=abcdef012345`, or use a `.env` file (see README for more details).
