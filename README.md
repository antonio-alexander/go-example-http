# go-example-http (github.com/antonio-alexander/go-example-http)

I put together this repository because I was _really_ curious about how I could serve static webpages with Go and by extension, how I could serve dynamic pages (using [html/templates](https://pkg.go.dev/html/template)). Although it wasn't very hard; I mostly had issues with SSL (https). Doing templates and serving static assets was relatively easy.

In _another_ project, I was figuring out how to do design work (html/java/css) and used the assets from that (see links in the [bibliography](#bibliography)). Although the [Dockerfile](./cmd/Dockerfile) and [docker-compose.yml](./docker-compose.yml) aren't super interesting and are run of the mill (please ignore them). The files used to generate a self-signed certificate are far more interesting, look at the [Makefile](./Makefile) and [openssl.conf](./config/openssl.conf).

## Getting Started

To get started, execute the following steps:

1. execute _make gen-certs_ to create the self-signed certificates for SSL
2. execute _make build_ to build the docker image
3. execute _make run_ to run the docker image in a container
4. ensure that your browser trusts the self-signed certificate
5. navigate to the webserver at [http://localhost](http://localhost) or [https://localhost](https://localhost)

> No, you don't have to run the docker container, you can run it using _go run_ or if you're a vs code kind of persion, the [launch.json](./.vscode/launch.json) has an entry. You can skip steps 2 and 3

Once everything is up and running, you can navigate to the website; there are two versions, one that's static and another that's served via template:

- [http://localhost](http://localhost)
- [http://localhost/index.html](http://localhost/index.html)

## Problems I Encountered...and Solved

While doing this, there were a handful of things I had to learn that weren't new, but I simply didn't know how. This is a list of all of the things that I encountered (or thought about):

Earlier, I mentioned that I had some significant trouble with https; one of the first things I had to figure out was how to generate a self-signed certificate

> The solve is in the [Makefile](./Makefile), but in general, I use [OpenSSL](https://www.openssl.org/) to generate the certs and output a key and cert file. There were a couple more caveats to this cert that I'll mention below. The actual command is _openssl req -x509 -newkey rsa:${rsa_bits} -sha256 -utf8 -days 1 -nodes -config ${ssl_config_file} -keyout ./certs/ssl.key -out ./certs/ssl.crt_.

When I ran the http server, with the generated certs, I found that chrome (and by extension Mac OS X) wouldn't load the page (because it was insecure). This was for two reasons: (1) the certificate wasn't signed by a valid [Certificate Authority](https://en.wikipedia.org/wiki/Certificate_authority) and (2) the certificate was lacking some of the (now) necessary requirements for a given certificate.

> To solve the first issue, I had to tell the operating system to trust the certificate using the [Keychain Access](https://support.apple.com/guide/keychain-access/what-is-keychain-access-kyca1083/mac); by doing this the OS trusts the certificate and by extension Chrome will trust the https site.
> to solve the second issue, I had to generate a certificate with additional attributes. I found that I couldn't supply these attributes via command line and found a nice stack overflow article that gave me a sample configuration. See the config in [openssl.conf](./config/openssl.conf)
> In addition, chrome has excellent tools to troubleshoot this specifically, if you right click the background of the website, click inspect and then go to the security tab, it'll tell you exactly why it's not "secure"

If you look at the Network tab of the Inspector, one of the errors I saw (and you'll probably see) is that [favicon.ico](./static/favicon.ico) is...missing. Functionally, it doesn't really matter (kinda); it's used to provide the icon in the tab/title bar of the browser.

> The solve is very easy for this, you simply need to be able to service favicon.ico and it should be at the root of your webserver; simple fix

I had a red herring early on, when I accessed the https website, I originally had an SSL issue, but the site wouldn't load correctly, like it couldn't load those files/assets as if they weren't available (i.e. 404). Again, I didn't have this problem with https, but only https.

> tldr; because the site was insecure (SSL) issue, it wouldn't load any additional assets nor would it render the page correctly (possibly related to the CSS not being loaded). Once I fixed the SSL/certificate issue, this issue also resolved itself, but I did learn a neat way to serve static files (look at the handler, it's super simple). In general, you don't "need" to _also_ serve static files that aren't being accessed directly; so if you have a template, I'm pretty sure you don't ALSO need a file server

I wondered if I'd have to worry about [CORS](https://en.wikipedia.org/wiki/Cross-origin_resource_sharing) issues and quickly realized I didn't need to worry about it because all of the resources...I was accessing were local to the website (no cross-site). But I could probably integrate it pretty easily using the solution I put together in [https://github.com/antonio-alexander/go-blog-cors](https://github.com/antonio-alexander/go-blog-cors)

Originally, I had a single web server where I would enable https or disable it, but I found that it was a little janky given that I'd hve to change the ports and it'd be a lot more interesting if I did both.

> This wasn't hard, you can see most if not all of the work I did in the [main.go](./internal/main.go). In general, you just need to manage the ports. I'm sure I could clean this up with nginx to be able to sense the protocol and forward it to the appropriate container port (maybe I'll look into it).

## Bibliography

Here are a list of links I used while trying to put this together (a lot are from stack overflow haha):

- [https://youtu.be/MBlkKE0GYGg?si=O-Hs2g9Js0edJwWt](https://youtu.be/MBlkKE0GYGg?si=O-Hs2g9Js0edJwWt)
- [https://gist.github.com/nguyendangminh/bd6e1e01df3c6cff139b1609fc1a646c](https://gist.github.com/nguyendangminh/bd6e1e01df3c6cff139b1609fc1a646c)
- [https://stackoverflow.com/questions/46100377/how-to-create-an-openssl-self-signed-certificate-using-san](https://stackoverflow.com/questions/46100377/how-to-create-an-openssl-self-signed-certificate-using-san)
- [https://serverfault.com/questions/47876/handling-http-and-https-requests-using-a-single-port-with-nginx](https://serverfault.com/questions/47876/handling-http-and-https-requests-using-a-single-port-with-nginx)
