# Commento User Installation Guide

## Table of Contents
- [Overview](#overview)
- [Option 1: Using Hosted Commento (Recommended for Most Users)](#option-1-using-hosted-commento-recommended-for-most-users)
- [Option 2: Self-Hosting Commento](#option-2-self-hosting-commento)
- [Adding Your Website](#adding-your-website)
- [Installing the Widget](#installing-the-widget)
- [Configuration Options](#configuration-options)
- [Testing Your Installation](#testing-your-installation)
- [Advanced Setup](#advanced-setup)
- [Troubleshooting](#troubleshooting)
- [Migration from Other Platforms](#migration-from-other-platforms)

---

## Overview

Commento is a privacy-focused commenting platform that you can add to your blog, website, or documentation site. This guide will walk you through the process of getting Commento up and running on your website.

### What You'll Need
- A website or blog where you want to add comments
- Ability to edit your website's HTML
- 10-15 minutes of setup time

### Two Deployment Options

1. **Hosted Service** (Easiest) - Sign up at a hosted Commento service, pay a subscription, and get started immediately
2. **Self-Hosted** (Advanced) - Run Commento on your own server with full control

---

## Option 1: Using Hosted Commento (Recommended for Most Users)

### Step 1: Sign Up for an Account

1. **Visit the Hosted Service**
   - Go to your Commento provider's website (e.g., `https://commento.io` or your custom instance)
   - Click on **Sign Up** in the navigation bar

2. **Create Your Account**

   Fill in the registration form:

   ```
   Email Address: your-email@example.com
   Full Name: John Doe
   Password: [Create a strong password]
   Confirm Password: [Re-enter your password]
   ```

3. **Choose Your Plan** (If Applicable)

   Some hosted services offer different pricing tiers:

   | Plan | Price | Features |
   |------|-------|----------|
   | **Starter** | $5/month | 1 domain, unlimited comments |
   | **Professional** | $15/month | 5 domains, priority support |
   | **Business** | $50/month | Unlimited domains, custom branding |

   *Note: Pricing may vary by provider. Self-hosted is always free.*

4. **Complete Email Verification**

   - Check your email inbox for a verification link
   - Click the link to confirm your account
   - You'll be redirected to the dashboard

### Step 2: Add Payment Information (Hosted Services)

1. **Navigate to Billing Settings**
   - Click on your name in the top right
   - Select **Settings** → **Billing**

2. **Enter Payment Details**

   Depending on the service provider, you may see:

   - Credit card form (Stripe integration)
   - PayPal option
   - Annual billing discount option

3. **Confirm Subscription**
   - Review your selected plan
   - Click **Subscribe** to activate

### Step 3: Add Your Website Domain

Now you're ready to add your website to Commento!

1. **Access the Dashboard**
   - After logging in, you'll see the Commento dashboard
   - On the left sidebar, click the **"+ New Domain"** button

   ![Dashboard view](https://via.placeholder.com/800x400?text=Dashboard+with+New+Domain+button)

2. **Fill in Domain Information**

   A modal will appear with two fields:

   ```
   Website Name: My Awesome Blog
   Website Domain: https://myblog.com
   ```

   **Important Notes:**
   - Use the **full domain** including `https://` protocol
   - Include `www.` if your site uses it (e.g., `https://www.myblog.com`)
   - Do NOT include trailing slash
   - You CANNOT change the domain after creation (only the name)

   **Examples:**
   ```
   ✅ Correct:
   - https://example.com
   - https://www.example.com
   - https://blog.example.com

   ❌ Incorrect:
   - example.com (missing protocol)
   - https://example.com/ (trailing slash)
   - http://example.com (wrong protocol if your site uses HTTPS)
   ```

3. **Submit**
   - Click **Add Domain**
   - Your domain will appear in the left sidebar

### Step 4: Get Your Installation Code

1. **Select Your Domain**
   - Click on your newly added domain in the left sidebar
   - It will be highlighted in blue when selected

2. **Go to Installation Section**
   - In the middle panel, click on **"Installation"**
   - You'll see a code snippet that looks like this:

   ```html
   <div id="commento"></div>
   <script defer
     src="https://commento.example.com/js/commento.js">
   </script>
   ```

3. **Copy the Code**
   - Click anywhere in the code box to select all
   - Press `Ctrl+C` (Windows/Linux) or `Cmd+C` (Mac) to copy
   - Or click the **Copy** button if available

---

## Option 2: Self-Hosting Commento

### Prerequisites

Before self-hosting, ensure you have:

- A Linux server (VPS, cloud instance, or dedicated server)
- Docker and Docker Compose installed
- A domain name pointed to your server
- Basic command-line knowledge

### Step 1: Server Setup

#### Option A: Using Docker Compose (Recommended)

1. **SSH into Your Server**
   ```bash
   ssh user@your-server.com
   ```

2. **Create a Directory for Commento**
   ```bash
   mkdir commento
   cd commento
   ```

3. **Create docker-compose.yml**
   ```bash
   nano docker-compose.yml
   ```

4. **Paste the Following Configuration**
   ```yaml
   version: '3.7'

   services:
     server:
       image: registry.gitlab.com/commento/commento:latest
       ports:
         - "8080:8080"
       environment:
         COMMENTO_ORIGIN: https://commento.yourdomain.com
         COMMENTO_POSTGRES: postgres://postgres:yourpassword@db:5432/commento?sslmode=disable
         COMMENTO_PORT: 8080

         # SMTP Configuration (for email notifications)
         COMMENTO_SMTP_HOST: smtp.gmail.com
         COMMENTO_SMTP_PORT: 587
         COMMENTO_SMTP_USERNAME: your-email@gmail.com
         COMMENTO_SMTP_PASSWORD: your-app-password
         COMMENTO_SMTP_FROM_ADDRESS: noreply@yourdomain.com

         # OAuth (Optional - enable social login)
         # COMMENTO_GOOGLE_KEY: your-google-client-id
         # COMMENTO_GOOGLE_SECRET: your-google-client-secret
         # COMMENTO_GITHUB_KEY: your-github-client-id
         # COMMENTO_GITHUB_SECRET: your-github-client-secret

         # Spam Detection (Optional - requires Akismet API key)
         # COMMENTO_AKISMET_KEY: your-akismet-key

       depends_on:
         - db
       restart: unless-stopped

     db:
       image: postgres:13-alpine
       environment:
         POSTGRES_DB: commento
         POSTGRES_USER: postgres
         POSTGRES_PASSWORD: yourpassword
       volumes:
         - postgres-data:/var/lib/postgresql/data
       restart: unless-stopped

   volumes:
     postgres-data:
   ```

5. **Important: Update the Configuration**

   Replace these values:
   - `COMMENTO_ORIGIN` → Your Commento URL (e.g., `https://commento.yourdomain.com`)
   - Database password → A strong random password
   - SMTP settings → Your email provider's settings
   - OAuth keys → If you want social login (see [OAuth Setup](#oauth-setup))

6. **Start Commento**
   ```bash
   docker-compose up -d
   ```

7. **Check Logs**
   ```bash
   docker-compose logs -f
   ```

   You should see:
   ```
   server_1  | Commento starting...
   server_1  | connected to database
   server_1  | starting server on 0.0.0.0:8080
   ```

#### Option B: Manual Installation

If you prefer not to use Docker, see the detailed installation instructions in the [ARCHITECTURE.md](./ARCHITECTURE.md#deployment-architecture) document.

### Step 2: Configure Reverse Proxy (Nginx)

1. **Install Nginx**
   ```bash
   sudo apt update
   sudo apt install nginx certbot python3-certbot-nginx
   ```

2. **Create Nginx Configuration**
   ```bash
   sudo nano /etc/nginx/sites-available/commento
   ```

3. **Add Configuration**
   ```nginx
   server {
       listen 80;
       server_name commento.yourdomain.com;

       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
           proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
           proxy_set_header X-Forwarded-Proto $scheme;
       }
   }
   ```

4. **Enable Site**
   ```bash
   sudo ln -s /etc/nginx/sites-available/commento /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl restart nginx
   ```

5. **Get SSL Certificate**
   ```bash
   sudo certbot --nginx -d commento.yourdomain.com
   ```

### Step 3: Create Your Admin Account

1. **Visit Your Commento Instance**
   - Open your browser and go to `https://commento.yourdomain.com`
   - Click **Sign Up**

2. **Create Admin Account**
   ```
   Email: admin@yourdomain.com
   Name: Admin
   Password: [Strong password]
   ```

3. **Verify Email** (if SMTP is configured)
   - Check your email and click the verification link
   - Or skip verification by manually updating the database:
   ```bash
   docker-compose exec db psql -U postgres commento
   UPDATE owners SET confirmedEmail = true WHERE email = 'admin@yourdomain.com';
   \q
   ```

---

## Adding Your Website

Whether you chose hosted or self-hosted, the process to add your website is the same:

### 1. Access Dashboard

- Log in to your Commento instance
- You'll see the main dashboard interface

### 2. Click "New Domain"

- Look for the **"+ New Domain"** button on the left sidebar
- Click it to open the domain creation modal

### 3. Enter Website Details

**Website Name**: A friendly name for your site (can be changed later)
```
Example: "My Tech Blog"
```

**Website Domain**: The exact domain where you'll embed Commento (CANNOT be changed later)
```
Example: "https://www.mytechblog.com"
```

**⚠️ Domain Must Match Exactly**

The domain you enter here MUST match the domain where your website is hosted:

| Your Website URL | Enter in Commento |
|------------------|-------------------|
| `http://example.com` | `http://example.com` |
| `https://example.com` | `https://example.com` |
| `https://www.example.com` | `https://www.example.com` |
| `https://blog.example.com` | `https://blog.example.com` |

### 4. Save Domain

- Click **Add Domain**
- Your domain will appear in the left sidebar
- Click on it to select it

---

## Installing the Widget

### Step 1: Get Your Embed Code

1. **Select Your Domain** (left sidebar)
2. **Click "Installation"** (middle panel)
3. **Copy the Code Snippet**

Your embed code will look like this:

```html
<div id="commento"></div>
<script defer
  src="https://commento.example.com/js/commento.js">
</script>
```

### Step 2: Add to Your Website

The installation method depends on your platform:

#### For Static HTML Sites

Open the HTML file where you want comments and add the code before the closing `</body>` tag:

```html
<!DOCTYPE html>
<html>
<head>
    <title>My Blog Post</title>
</head>
<body>
    <article>
        <h1>My Blog Post Title</h1>
        <p>Your blog post content goes here...</p>
    </article>

    <!-- Commento Comments -->
    <div id="commento"></div>
    <script defer
      src="https://commento.example.com/js/commento.js">
    </script>
</body>
</html>
```

#### For WordPress

**Method 1: Using a Plugin (Easiest)**

1. Install the "Insert Headers and Footers" plugin
2. Go to **Settings** → **Insert Headers and Footers**
3. Paste the Commento code in the **Footer** section
4. Save changes

**Method 2: Manual Theme Edit**

1. Go to **Appearance** → **Theme Editor**
2. Open `single.php` or `content.php`
3. Find the location after the post content
4. Paste the Commento embed code
5. Update file

**Method 3: Using a Custom HTML Block**

1. Edit your blog post
2. Add a **Custom HTML** block at the bottom
3. Paste the Commento code
4. Publish

#### For Hugo (Static Site Generator)

1. **Create a Partial Template**

   Create `layouts/partials/commento.html`:

   ```html
   {{ if not .Params.nocomments }}
   <div id="commento"></div>
   <script defer
     src="https://commento.example.com/js/commento.js">
   </script>
   {{ end }}
   ```

2. **Add to Single Post Template**

   Edit `layouts/_default/single.html`:

   ```html
   {{ define "main" }}
   <article>
     <h1>{{ .Title }}</h1>
     {{ .Content }}
   </article>

   {{ partial "commento.html" . }}
   {{ end }}
   ```

3. **Rebuild Site**
   ```bash
   hugo
   ```

#### For Jekyll

1. **Create an Include**

   Create `_includes/commento.html`:

   ```html
   <div id="commento"></div>
   <script defer
     src="https://commento.example.com/js/commento.js">
   </script>
   ```

2. **Add to Post Layout**

   Edit `_layouts/post.html`:

   ```html
   ---
   layout: default
   ---
   <article>
     <h1>{{ page.title }}</h1>
     {{ content }}
   </article>

   {% include commento.html %}
   ```

3. **Build Site**
   ```bash
   jekyll build
   ```

#### For Ghost

1. **Go to Settings**
   - Navigate to **Settings** → **Code Injection**

2. **Add to Post Footer**

   Paste in the **Post Footer** section:

   ```html
   <div id="commento"></div>
   <script defer
     src="https://commento.example.com/js/commento.js">
   </script>
   ```

3. **Save**

#### For React/Next.js

1. **Create a Component**

   Create `components/Commento.js`:

   ```jsx
   import { useEffect } from 'react';

   export default function Commento() {
     useEffect(() => {
       const script = document.createElement('script');
       script.src = 'https://commento.example.com/js/commento.js';
       script.defer = true;
       document.body.appendChild(script);

       return () => {
         document.body.removeChild(script);
       };
     }, []);

     return <div id="commento" />;
   }
   ```

2. **Use in Your Post Component**

   ```jsx
   import Commento from '../components/Commento';

   export default function BlogPost({ post }) {
     return (
       <article>
         <h1>{post.title}</h1>
         <div dangerouslySetInnerHTML={{ __html: post.content }} />

         <Commento />
       </article>
     );
   }
   ```

#### For Gatsby

1. **Create a Component**

   Create `src/components/Commento.js`:

   ```jsx
   import React, { useEffect } from 'react';

   const Commento = () => {
     useEffect(() => {
       const script = document.createElement('script');
       script.src = 'https://commento.example.com/js/commento.js';
       script.async = true;

       const commento = document.getElementById('commento');
       if (commento) {
         commento.appendChild(script);
       }

       return () => {
         const commento = document.getElementById('commento');
         if (commento) {
           while (commento.firstChild) {
             commento.removeChild(commento.firstChild);
           }
         }
       };
     }, []);

     return <div id="commento" />;
   };

   export default Commento;
   ```

2. **Use in Blog Template**

   Edit `src/templates/blog-post.js`:

   ```jsx
   import Commento from '../components/Commento';

   const BlogPostTemplate = ({ data }) => {
     return (
       <Layout>
         <article>
           <h1>{data.markdownRemark.frontmatter.title}</h1>
           <div dangerouslySetInnerHTML={{ __html: data.markdownRemark.html }} />
         </article>

         <Commento />
       </Layout>
     );
   };
   ```

---

## Configuration Options

### Basic Configuration

You can customize the widget behavior using data attributes:

```html
<script defer
  src="https://commento.example.com/js/commento.js"
  data-auto-init="true"
  data-css-override="/custom-commento-styles.css"
  data-no-fonts="false"
  data-hide-deleted="false"
  data-page-id="custom-page-identifier">
</script>
```

| Attribute | Default | Description |
|-----------|---------|-------------|
| `data-auto-init` | `true` | Automatically load comments on page load |
| `data-css-override` | (none) | URL to custom CSS file |
| `data-no-fonts` | `false` | Don't load custom fonts |
| `data-hide-deleted` | `false` | Hide deleted comments |
| `data-page-id` | Current URL path | Custom page identifier |

### Custom Styling

Create a custom CSS file to override default styles:

**custom-commento.css:**
```css
/* Change primary color */
.commento-root {
  font-family: 'Your Custom Font', sans-serif;
}

/* Style comment cards */
.commento-comment-card {
  background: #f9f9f9;
  border-left: 4px solid #4CAF50;
  border-radius: 8px;
}

/* Style buttons */
.commento-button {
  background: #4CAF50;
  color: white;
  border-radius: 4px;
}

.commento-button:hover {
  background: #45a049;
}

/* Style vote buttons */
.commento-upvote:hover {
  color: #4CAF50;
}

.commento-downvote:hover {
  color: #f44336;
}
```

Then reference it in your embed code:

```html
<script defer
  src="https://commento.example.com/js/commento.js"
  data-css-override="https://yourdomain.com/custom-commento.css">
</script>
```

### Dashboard Settings

Configure your domain settings in the dashboard:

#### General Settings

1. **Click your domain** → **General** → **General tab**

   - **Website Name**: Change display name
   - **Authentication Options**:
     - ☑ Anonymous comments
     - ☑ Email address login
     - ☑ Google login (if configured)
     - ☑ GitHub login (if configured)
     - ☑ Twitter login (if configured)
     - ☑ Single sign-on (advanced)

#### Moderation Settings

1. **Click your domain** → **Moderation** → **General tab**

   **Comment Filtering:**
   - ☑ Automatic spam filtering (uses Akismet)
   - ☐ Require all comments to be approved manually
   - ☐ Require anonymous comments to be approved manually

   **Email Schedule:**
   - ◉ Whenever a new comment is created
   - ○ Only for comments pending moderation
   - ○ Do not email moderators

   **Comment Sorting:**
   - ◉ Most upvoted first
   - ○ Newest first
   - ○ Oldest first

#### Add Moderators

1. **Click your domain** → **Moderation** → **Add/Remove Moderators tab**
2. Enter moderator's email address
3. Click **Add moderator**

Moderators can:
- Approve comments
- Delete comments
- Sticky comments (pin to top)
- Lock threads

Moderators CANNOT:
- Access the dashboard
- Change settings
- Add/remove other moderators

---

## Testing Your Installation

### Step 1: Visit Your Website

1. Open your website in a browser
2. Navigate to a page where you added the Commento code
3. Scroll down to where comments should appear

### Step 2: Verify Widget Loads

You should see:
- A comment submission form
- Login options (if not logged in)
- Any existing comments

**If you see nothing:**
- Check browser console for errors (F12 → Console tab)
- Verify the domain in your Commento dashboard matches your website URL exactly
- Ensure the embed script `src` URL is correct

### Step 3: Post a Test Comment

1. **Create a Test Commenter Account**
   - Click **Sign up** in the comment widget
   - Use a test email
   - Create a password

2. **Write a Comment**
   ```
   This is a test comment!
   ```

3. **Submit**
   - Click **Add Comment**
   - Your comment should appear immediately (unless moderation is enabled)

### Step 4: Test Voting

1. Click the ↑ upvote button
2. Score should increase
3. Click the ↓ downvote button
4. Score should decrease

### Step 5: Test Threading

1. Click **Reply** on your test comment
2. Write a reply
3. Submit
4. Reply should appear indented below the parent

### Step 6: Check Dashboard

1. Go to your Commento dashboard
2. Click on your domain
3. Click **Statistics**
4. You should see:
   - View count increased
   - Comment count: 1+ comments

---

## Advanced Setup

### OAuth Configuration

Enable social login for your users.

#### Google OAuth

1. **Create Google OAuth Credentials**

   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing
   - Navigate to **APIs & Services** → **Credentials**
   - Click **Create Credentials** → **OAuth client ID**
   - Application type: **Web application**
   - Authorized redirect URIs: `https://commento.example.com/api/oauth/google/callback`
   - Copy **Client ID** and **Client Secret**

2. **Add to Commento**

   **For Self-Hosted:**
   Edit `docker-compose.yml`:
   ```yaml
   environment:
     COMMENTO_GOOGLE_KEY: your-client-id.apps.googleusercontent.com
     COMMENTO_GOOGLE_SECRET: your-client-secret
   ```
   Restart: `docker-compose restart`

   **For Hosted:**
   Contact your provider to enable Google OAuth.

3. **Enable in Domain Settings**
   - Dashboard → Your Domain → General → Authentication Options
   - Check ☑ **Google login**
   - Click **Save Changes**

#### GitHub OAuth

1. **Create GitHub OAuth App**

   - Go to [GitHub Developer Settings](https://github.com/settings/developers)
   - Click **New OAuth App**
   - Application name: "Commento Comments"
   - Homepage URL: `https://yourdomain.com`
   - Authorization callback URL: `https://commento.example.com/api/oauth/github/callback`
   - Copy **Client ID** and **Client Secret**

2. **Add to Commento**

   ```yaml
   environment:
     COMMENTO_GITHUB_KEY: your-github-client-id
     COMMENTO_GITHUB_SECRET: your-github-client-secret
   ```

3. **Enable in Domain Settings**
   - Check ☑ **GitHub login**
   - Save

#### Twitter OAuth

1. **Create Twitter App**

   - Go to [Twitter Developer Portal](https://developer.twitter.com/en/portal/dashboard)
   - Create a new app
   - Enable OAuth 1.0a
   - Callback URL: `https://commento.example.com/api/oauth/twitter/callback`
   - Copy **API Key** and **API Secret**

2. **Add to Commento**

   ```yaml
   environment:
     COMMENTO_TWITTER_KEY: your-twitter-api-key
     COMMENTO_TWITTER_SECRET: your-twitter-api-secret
   ```

### Email Notifications

Configure SMTP to send email notifications.

#### Using Gmail

1. **Generate App Password**
   - Go to [Google Account Security](https://myaccount.google.com/security)
   - Enable 2-Step Verification
   - Go to **App passwords**
   - Generate a new app password for "Mail"

2. **Configure Commento**

   ```yaml
   environment:
     COMMENTO_SMTP_HOST: smtp.gmail.com
     COMMENTO_SMTP_PORT: 587
     COMMENTO_SMTP_USERNAME: your-email@gmail.com
     COMMENTO_SMTP_PASSWORD: your-app-password
     COMMENTO_SMTP_FROM_ADDRESS: noreply@yourdomain.com
   ```

#### Using SendGrid

1. **Create SendGrid Account**
   - Sign up at [SendGrid](https://sendgrid.com)
   - Generate an API key

2. **Configure Commento**

   ```yaml
   environment:
     COMMENTO_SMTP_HOST: smtp.sendgrid.net
     COMMENTO_SMTP_PORT: 587
     COMMENTO_SMTP_USERNAME: apikey
     COMMENTO_SMTP_PASSWORD: your-sendgrid-api-key
     COMMENTO_SMTP_FROM_ADDRESS: noreply@yourdomain.com
   ```

### Spam Protection

Enable Akismet spam detection.

1. **Get Akismet API Key**
   - Sign up at [Akismet](https://akismet.com/)
   - Get your API key

2. **Add to Commento**

   ```yaml
   environment:
     COMMENTO_AKISMET_KEY: your-akismet-key
   ```

3. **Enable in Dashboard**
   - Domain → Moderation → General
   - Check ☑ **Automatic spam filtering**
   - Save

### Custom Domain (Self-Hosted)

Use your own domain instead of a subdomain.

1. **Update DNS**
   ```
   Type: A
   Name: commento.yourdomain.com
   Value: [Your server IP]
   TTL: 3600
   ```

2. **Update Commento Configuration**
   ```yaml
   environment:
     COMMENTO_ORIGIN: https://commento.yourdomain.com
   ```

3. **Get SSL Certificate**
   ```bash
   sudo certbot --nginx -d commento.yourdomain.com
   ```

---

## Troubleshooting

### Comments Not Appearing

**Problem**: Widget doesn't load on your website.

**Solutions**:

1. **Check Domain Match**
   - Dashboard → Your Domain → Check domain URL
   - Must EXACTLY match your website (including `https://` and `www.`)

2. **Check Browser Console**
   - Press F12 → Console tab
   - Look for errors like:
     ```
     Access to script at 'https://commento.example.com/js/commento.js'
     from origin 'https://yourdomain.com' has been blocked by CORS policy
     ```
   - This means domain mismatch

3. **Check Script URL**
   - Verify `src="https://commento.example.com/js/commento.js"` is correct
   - Try visiting the URL directly in your browser
   - Should download a JavaScript file

4. **Check Page Source**
   - View page source (Ctrl+U)
   - Search for "commento"
   - Verify the script tag is present

### CORS Errors

**Problem**: Console shows CORS error.

**Solution**: Domain mismatch. Update your domain in the dashboard or your embed code.

### Comments Pending Moderation

**Problem**: Comments don't appear immediately.

**Solution**: This is by design if moderation is enabled.

- Dashboard → Domain → Moderation → General
- Uncheck **Require all comments to be approved manually**
- Save

### Email Notifications Not Working

**Problem**: Not receiving email notifications.

**Solutions**:

1. **Check SMTP Configuration**
   ```bash
   docker-compose logs server | grep SMTP
   ```

2. **Test SMTP Settings**
   - Use an online SMTP tester
   - Verify credentials are correct

3. **Check Email Settings**
   - Dashboard → Domain → Moderation → Email Schedule
   - Ensure it's not set to "Do not email moderators"

4. **Check Spam Folder**

### Widget Styling Issues

**Problem**: Widget looks broken or unstyled.

**Solutions**:

1. **Check Content Security Policy (CSP)**
   - Your site may have CSP headers blocking external scripts
   - Add to your CSP:
     ```
     script-src 'self' https://commento.example.com;
     style-src 'self' https://commento.example.com;
     ```

2. **Check Ad Blockers**
   - Disable ad blocker temporarily
   - If this fixes it, create an exception

3. **Clear Browser Cache**
   - Hard refresh: Ctrl+Shift+R (Windows/Linux) or Cmd+Shift+R (Mac)

### Self-Hosted: Container Won't Start

**Problem**: Docker container exits immediately.

**Solutions**:

1. **Check Logs**
   ```bash
   docker-compose logs server
   ```

2. **Common Issues**:
   - Database connection failed → Check PostgreSQL is running
   - Invalid configuration → Check environment variables
   - Port already in use → Change port in docker-compose.yml

3. **Database Connection**
   ```bash
   docker-compose exec db psql -U postgres commento
   ```

   If this fails, PostgreSQL isn't running properly.

### Performance Issues

**Problem**: Comments load slowly.

**Solutions**:

1. **Use CDN** (Self-Hosted)
   ```yaml
   environment:
     COMMENTO_CDN_PREFIX: https://cdn.yourdomain.com
   ```

2. **Enable Browser Caching**
   - Add cache headers to static assets
   - See ARCHITECTURE.md for Nginx configuration

3. **Optimize Database**
   ```bash
   docker-compose exec db psql -U postgres commento
   VACUUM ANALYZE;
   REINDEX DATABASE commento;
   ```

---

## Migration from Other Platforms

### Migrating from Disqus

1. **Export from Disqus**
   - Go to [Disqus Export](http://disqus.com/admin/discussions/export/)
   - Click **Export Comments**
   - Wait for email with download link

2. **Import to Commento**
   - Dashboard → Your Domain → Import/Export → Disqus tab
   - Paste the Disqus export URL
   - Click **Import**
   - Wait for completion (may take several minutes)

3. **Verify Import**
   - Check comment count in Statistics
   - View a few posts to verify comments appear

4. **Update Website Code**
   - Remove Disqus embed code
   - Add Commento embed code
   - Deploy changes

### Migrating from WordPress Comments

1. **Export WordPress Comments**
   - Use a WordPress plugin like "WP All Export"
   - Export comments to CSV

2. **Convert to Commento Format**
   - Currently no automated tool
   - Manual import via API or database
   - See ARCHITECTURE.md for database schema

3. **Alternative: Run Both Temporarily**
   - Keep WordPress comments visible
   - Add Commento for new comments
   - Gradually migrate old content

---

## Next Steps

### After Installation

1. **Customize Appearance**
   - Create custom CSS (see [Custom Styling](#custom-styling))
   - Match your website's design

2. **Configure Moderation**
   - Add moderators if needed
   - Set up spam filtering
   - Configure email notifications

3. **Enable OAuth**
   - Add Google, GitHub, or Twitter login
   - Makes it easier for users to comment

4. **Monitor Usage**
   - Check dashboard statistics regularly
   - Review flagged comments
   - Approve pending comments

5. **Optimize Performance**
   - Set up CDN for static assets
   - Enable browser caching
   - Monitor server resources

### Getting Help

- **Documentation**: [Official Docs](https://docs.commento.io)
- **GitHub Issues**: [Report bugs](https://github.com/adtac/commento/issues)
- **IRC**: #commento on Freenode
- **Email**: support@commento.io (hosted services)

---

## Frequently Asked Questions

### Q: Can I use Commento on multiple websites?

**A**: Yes! You can add unlimited domains (depending on your plan for hosted services). Each domain is configured separately in the dashboard.

### Q: Is Commento GDPR compliant?

**A**: Yes, Commento is designed with privacy in mind. It doesn't track users, doesn't sell data, and respects user privacy.

### Q: Can I migrate comments later?

**A**: Yes, you can export your comments at any time from Dashboard → Domain → General → Export Data.

### Q: Does Commento work with static site generators?

**A**: Absolutely! Commento works with Jekyll, Hugo, Gatsby, Next.js, and any other static site generator.

### Q: Can I customize the appearance?

**A**: Yes, use the `data-css-override` attribute to link to your custom CSS file.

### Q: How do I delete a comment?

**A**: As a moderator or owner, you'll see a **Delete** button on each comment when logged in on your website.

### Q: Can users edit their comments?

**A**: Yes, users can edit their own comments by clicking the **Edit** button.

### Q: Does it work on mobile?

**A**: Yes, Commento is fully responsive and works great on mobile devices.

### Q: Can I hide the "Powered by Commento" footer?

**A**: This depends on your hosting provider's terms. Self-hosted instances can modify the code freely.

### Q: How much does it cost?

**A**: Self-hosting is free (open source). Hosted services typically charge $5-50/month depending on features and number of domains.

---

## Summary

Congratulations! You've successfully installed Commento on your website. Here's what you accomplished:

✅ Created a Commento account (hosted or self-hosted)
✅ Added your website domain
✅ Installed the comment widget
✅ Configured settings and moderation
✅ Tested the installation

Your readers can now engage with your content through a privacy-respecting, fast, and beautiful commenting system.

**Enjoy your new commenting system! 🎉**
