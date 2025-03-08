# mailReceipt
A simple (in terms of codebase and usability) mail receipt api.

## Usage
The api is easy to set up and use:
- Clone the repo with `git clone https://github.com/spacefall/mailReceipt.git`.
- Build it with `go build`.  
- Create a .env file with the following fields (replace the values with your own):
```env
# Replace the values with your own, the current values are just placeholders.
# URI to your postgres database, if you don't have one fill out the following fields.
DATABASE_URL=postgres://[username[:password]@][host[:port],]/[database][?parameter_list]
# Email address and password for the mail account that will send the receipts.
EMAIL_USERNAME=example@gmail.com
EMAIL_PASSWORD=examplePassword
# SMTP server host, the api will use SSL by default so check that your server supports it.
EMAIL_HOST=smtp.example.com
```
- Run the api with `./mailReceipt`.

And you're done! The api will be running on `localhost:3000` by default.

## Endpoints
The api has the following endpoints:
- '/track'
  - POST: Create a new tracking entry.
- '/track/{id}'
  - GET: Gets info about a tracking entry.
  - DELETE: Deletes a tracking entry.
- '/track/{id}/pixel'
  - GET: Returns a 1x1 pixel image that will be used to track when the email is opened.
- '/track/{id}/url/{url}'
  - GET: Redirects to the url provided and records that the link has been opened.

A full list and description of the endpoints can be found in the swagger file or [here](https://mailReceipt.5822.it).

## Why?
Mail trackers are kind of a pain (or at least the ones in the first page of Google):  
You have to register, install an extension, give them access to your account, say that you're a business, and you still get a service that sometimes doesn't work.

This, instead, is a simple api self-hostable api that works about the same, doesn't require any personal information, works about everywhere (even non emails if you wanted) and is open source.  
Also, it was a nice learning experience for me.
