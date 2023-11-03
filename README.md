# HomeLinks Web Application

HomeLinks is a simple web application that allows you to manage and access your favorite links. It provides basic authentication for public IP addresses and serves as an example of a Go web application using the Gin framework, embedded files, and Docker containerization.

## Prerequisites

Before you begin, ensure you have the following dependencies installed:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

Follow these steps to build and run the HomeLinks web application:

### 1. Clone the Repository

```bash
git clone https://github.com/bariiss/homelinks.git
cd homelinks
```

### Create a Configuration File

You can customize the links displayed in the application by creating a homelinks.json file. Here's an example structure:

```json
[
  {
    "Name": "google",
    "Text": "Google",
    "URL": "https://www.google.com",
    "Color": "#4285F4",
    "TextColor": "white",
    "AltText": "Search with Google"
  },
  {
    "Name": "facebook",
    "Text": "Facebook",
    "URL": "https://www.facebook.com",
    "Color": "#F48042",
    "TextColor": "black",
    "AltText": "Connect with Facebook"
  },
  {
    "Name": "twitter",
    "Text": "Twitter",
    "URL": "https://www.twitter.com",
    "Color": "#F44242",
    "TextColor": "white",
    "AltText": "Tweet with Twitter"
  }
]
```

Place this file in the project root directory.

### 3. Build and Run with Docker

```bash
# Build the Docker image
make docker-build

# Run the application interactively (for debugging)
make docker-run-tty

# Run the application as a daemon
make docker-run-daemon
```

The application will be accessible at http://localhost:8080.

### 4. Access the HomeLinks Web Application

Open a web browser and navigate to http://localhost:8080 to access the HomeLinks web application.

### Configuration

* Basic authentication is required for public access. Use the following credentials:
> username: ```user```
> pass: ```pass```

### Customization

You can customize the links and the application's appearance by modifying the homelinks.json file and the template files in the templates directory.

### Troubleshooting

If you encounter any issues, you can stop the Docker container using the following command:

```make docker-stop```

To remove the Docker image, use the following command:

```make docker-image-rm```

This project is licensed under the MIT License - see the LICENSE file for details.