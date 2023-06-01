## Running the example

### First you need to install dependencies for the back-end and front-end server

#### For the Back-end

You need to run this inside the main folder if you didn't before.
```bash
go mod tidy
```

#### For the Fron-end

You need to run this inside the ./chat/public folder:

```bash
npm install
# or
yarn install
# or
pnpm install
```

### Once you install all dependencies

#### For the Back-end

On the chat folder run the Pub/Sub server with this command:

```bash
go run ./...
```

After that, you will see your server running on: `http://localhost:8080`

#### For the Fron-end

And for running the front end:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
```

After that, you will see the front end running on: `http://localhost:3000`

## Testing the app

Now you can open a couple of tabs in your preferred browser on this URL `http://localhost:3000`, write a name.
And send messages, you will see the messages in all the tabs you open.

And in your terminal, you can see the messages that the back-end server receives.
