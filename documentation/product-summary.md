# Product Summary

## Overview

`goog` is a command-line interface for Google Mail and Google Calendar, written in Go. It provides programmatic access to Gmail and Calendar APIs through a clean, scriptable interface designed for both human operators and AI agents.

## Goals

- **Unified CLI** - Single tool for Gmail and Calendar operations
- **Multi-account support** - Manage separate accounts for different contexts
- **Scriptability** - JSON output format for easy parsing and automation
- **Security** - OAuth2/PKCE authentication with tokens stored in system keyring
- **Clean Architecture** - Maintainable, testable codebase following Go best practices

## Key Capabilities

### Gmail Operations
- Message management: list, read, search, send, reply, forward
- Message actions: trash, archive, delete, modify labels, mark read/unread
- Draft management: create, update, send, delete drafts
- Label management: create, update, delete labels
- Thread operations: view conversations, bulk label changes

### Calendar Operations
- Event management: list, create, update, delete events
- Quick views: today, this week, upcoming
- Recurring events: list instances, modify individual occurrences
- Calendar management: create, update, delete calendars
- Sharing: ACL rules for calendar access control
- Availability: free/busy queries
- RSVP: accept, decline, tentative responses

### Multi-Account Support
- Add multiple Google accounts with aliases
- Switch default account easily
- Override per-command with `--account` flag
- Separate credential storage per account

## Target Users

1. **Developers** - Automate Gmail and Calendar workflows
2. **Power users** - Efficient email and calendar management from terminal
3. **AI agents** - Programmatic access with JSON output for easy parsing
4. **DevOps** - Script calendar and email notifications

## Status

**Phase**: Complete (v1.0)

All planned features implemented:
- Foundation layer (auth, config, accounts)
- Gmail core operations
- Calendar core operations
- Advanced features (threads, ACL, recurring events)
