## Overview

- We discuss about the different features supported by Arka.

### Database

- Currently supported providers are:
    - Gorm

- To use multiple hosts (for replica setup):
    - Set the environment variable `DB_HOSTS` to an appropriate value.
    - See the `sample.env` file for an example.

### Cache

- Currently supported providers are:
    - Redis

### Email

- Currently supported providers are:
    - Mailgun
    - Amazon SES

### SMS

- Currently supported providers are:
    - SMS Broadcast
    - ClickSend
    - Amazon SNS
    - Termii

### Payment Gateway

- Currently supported providers are:
    - Stripe

### Scheduler

- Currently supported providers are:
    - Cron

### Queuing

- Currently supported providers are:
    - Amazon SQS
