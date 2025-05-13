# Go PR Reviewer

A GitHub pull request code review automation tool written in Go that uses Google's Gemini AI to provide automated code suggestions on pull requests.

This project automates code reviews for GitHub pull requests. When a pull request is created or updated, the application:

Retrieves the code changes from the pull request  
Sends the changes to Google's Gemini AI for analysis  
Posts AI-generated code suggestions as comments on the pull request

## Installation

### Prerequisites

Go 1.24 or higher  
GitHub personal access tokens(PAT) token with

- Read and Write access to pull requests
  Google Gemini API key

Create a .env file in the cmd/api-server directory with the following tokens

- GITHUB_KEY (your PAT token)
- GEMINI_KEY (Your Gemini API key)
- GEMINI_VERSION (The version of the gemini api you want to use)

## Usage

Start the server  
Configure a GitHub webhook to send pull request events to your server endpoint. This is done in the settings page of the repository.

The server will listen on port 3000 for webhook events.

To start receiving webhook request locally you can use [smee.io](https://smee.io/)
