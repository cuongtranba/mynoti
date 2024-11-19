![Go Test](https://github.com/cuongtranba/mynoti/actions/workflows/go-test.yml/badge.svg)
# Comic Tracking Bot

## Overview
This project is designed to track websites where I usually read comics. Unfortunately, these sites are often unofficial and lack notification features for new chapters.

## Features
1. **Simple Cron Job**: The project runs as a cron job, monitoring changes in the HTML structure of specified websites.
2. **Telegram Notifications**: Sends notifications to a Telegram chat whenever a new chapter is detected.

## How It Works
- Periodically checks for updates on the specified comic websites.
- Identifies changes or new chapters based on HTML structure.
- Notifies the user via a Telegram bot.

## Prerequisites
- Basic knowledge of cron jobs and scheduling.
- A Telegram bot set up to receive notifications.
- Access to the comic websites you want to track.

## Setup
1. Clone this repository.
2. Configure the websites to track and your Telegram bot token in the project settings.
3. Schedule the cron job using a task scheduler like `cron` (Linux/macOS) or Task Scheduler (Windows).

## Disclaimer
This project is for personal use only and is not intended to promote or endorse the use of unofficial or pirated websites.

## License
MIT License
