# Akismet Integration

Akismet is a leading spam filtering service that's used by millions of websites. For non-commercial websites, you can obtain a free API key (you can also pay as much as you want). Commento is not associated with Akismet in anyway, and Akismet integration is completely optional; if you do not want to use their service, comments will not be filtered. However, if you do choose to enable Akismet integration, please be advised that *all* comments will be sent to Akismet's servers for filetering.

To enable Akismet integration in Commento, go to [their website](https://akismet.com) and signup for an account. Choose your [appropriate plan](https://akismet.com/account/upgrade), and then go to your dashboard. Here, you'll find your 12-digit hexadecimal API key. To make Commento use this key, either set an environment variable though `export AKISMET_KEY=abcdef012345`, or use a `.env` file (see README for more details on how to use `.env` files).
