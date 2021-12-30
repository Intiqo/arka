## Code Structure

Here, we'll learn how `arka` is organised in terms of packages

## Packages

- [Cache Manager](./cache) - Caches data and retrieves them
- [Cloud Manager](./cloud) - Manages all connections to cloud services such as AWS
- [Config Manager](./config) - Gets configuration data
- [Constants](./constants) - Commonly used constants across the app
- [Database Manager](./database) - Manages the generic database operations
- [Dependency Manager](./dependency) - Registers and retrieves dependencies across the app
- [Email Manager](./email) - Responsible for sending emails
- [Exception Manager](./exception) - Creates and manages errors across the app
- [Event Manager](./event) - Creates and manages events published and consumed across the app
- [Excelize](./excelize) - Provides APIs to create and manage an Excel workbook
- [File Manager](./file) - Manages file operations such as uploading to a cloud bucket
- [Logger](./logger) - Provides logging capabilities
- [Payment Manager](./payment) - Provides APIs to handle payment related operations
- [Scheduler Manager](./scheduler) - Provides APIs to create and manage scheduled jobs
- [SMS](./sms) - Responsible for sending SMS
- [Template](./template) - Provides methods to create text and html templates
- [Util](./util) - Provides utility functionality such as generating a token, shortening a url, etc
- [Version](./version) - Provides version information of a service
