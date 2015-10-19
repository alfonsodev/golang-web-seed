# Golang web seed
This project is a basic starting point for building a website in golang without frameworks,
but with small packages for each functionallity.  
## Features
- Middleware with [negroni](https://github.com/codegangsta/negroni)
- Template system with [pongo2](https://github.com/flosch/pongo2)
- Routing with [gorilla mux](https://github.com/gorilla/mux)
- Sessions [gorilla sessions](https://github.com/gorilla/sessions)
- Google auth []()

## Configuration 
The app is configured with the following environment variables.
In order to obtain google client: id, secret and redirect url, you'll need to viist your [developer console](https://console.developers.google.com)
### Mandatory
- GOOGLE_CLIENT_ID=000000000000-00000000000000000000000000000000.apps.googleusercontent.com
- GOOGLE_CLIENT_SECRET=000000000000000000000000
- GOOGLE_CLIENT_REDIRECT=http://localhost:3000/callback
### Optional 
- SGW_PORT
  Port where the application will listen for connection, defaults to `3000` 
- SGW_SESSION
 Session name, defaults to `sgw`
