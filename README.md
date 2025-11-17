# Blog Aggregator

A CLI application to follow RSS feeds and read blog posts directly from your terminal.

## Features
- 📰 Follow/unfollow RSS feeds
- 👤 User authentication and management  
- 🔄 Automatic feed aggregation
- 📖 Browse posts with pagination
- 💾 PostgreSQL data persistence

## Prerequisites
- **PostgreSQL** 18.1+
- **Go** 1.25.1+

## Installation

1. Clone the repository:
```bash
git clone git@github.com:luis-octavius/blog-aggregator.git
cd blog-aggregatorregator
```

2. Build and install:

```bash
go install .
```

3. (Optional) Create an alias for easier use:

```bash
alias gator="blog-aggregator"
```

## Usage

### User Management

```bash

# Register new user
gator register <username>

# Login as existing user  
gator login <username>

# List all users
gator users

# Reset all data
gator reset
```

### Feed Management

```bash

# Add a new RSS feed
gator addfeed <name> <url>

# List all feeds in system
gator feeds

# Follow/unfollow feeds
gator follow <url>
gator unfollow <url>

# List feeds followed by current user
gator following

Reading Content
bash

# Fetch latest posts from followed feeds
gator agg

# Browse posts (default: 2 most recent)
gator browse
gator browse 5    # Show 5 posts
gator browse 10   # Show 10 posts
```

### Example Workflow

```bash

# Reset and start fresh
gator reset

# Register and login
gator register luis
gator login luis

# Add some feeds
gator addfeed "Go Blog" "https://go.dev/blog/feed.atom"
gator addfeed "The New Stack" "https://thenewstack.io/feed"

# See what you're following
gator following

# Fetch and read posts
gator agg
gator browse 5
```

## Configuration

The application automatically creates a configuration file at ~/.gatorconfig.json storing your database URL and current user.
Database
Uses PostgreSQL with automatic schema migrations. Make sure your database is running and accessible via the connection string in the config file.


