## SnippetBox - A Text and Code Snippet Sharing Application
##### SnippetBox is a lightweight web application that allows users to paste and share snippets of text and code, similar to Pastebin or GitHub's Gists. It uses a managed MySQL database hosted on Aiven and is built with Go.

#### Features

- Create and store text/code snippets
- View individual snippets via unique URLs
- Simple and intuitive API for snippet management

#### Getting Started
##### Prerequisites
- Go (v1.16 or higher)
- MySQL (or access to a managed MySQL instance like Aiven)

##### Installation

1. Clone the repository:

        git clone https://www.github.com/MlondiMchunu/snippetbox.git
        cd snippetbox

2. Configure the database:

- Set up your MySQL connection details in cmd/web/main.go (or via environment variables).

3. Run the application:

        go run ./cmd/web  

- The app will start on http://localhost:4000


##### API Endpoints:

1. Create a Snippet:

- Method POST 

      http://localhost:4000/snippet/create

- Request Body:

      {  
        "title": "Example Snippet",  
        "content": "This is a test snippet.",  
        "expires": "24h"  // Optional: "1h", "24h", "7d", or "never"  
      }  

2. View A Snippet:

- Method GET

      http://localhost:4000/snippet/view/2

- Example: 

      GET http://localhost:4000/snippet/view/2  

3. View List of all Snippets:

- Method GET
       http://localhost:4000


##### Deployment
This project is deployed on Render for easy access.

##### License
This project is open-source and available under the MIT License.


<img width="857" height="400" alt="image" src="https://github.com/user-attachments/assets/2aa2ed6b-fdef-4afa-8dd9-e5ec3c54b1b2" />



🛠 🚧 under construction 😐
