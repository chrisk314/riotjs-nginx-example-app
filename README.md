# Example Riot.js app

## Run

### Development
Running in development provides hot reloading support via webpack-dev-server.

```
docker-compose up
```

### Production
In production a static version of the site is baked into an nginx docker image which is deployed
with auto-renewing SSL certs facilitated by certbot.

---
##### First time setup
The production nginx config contains a generic server name: `example.com`. In order to make use of
https it's necessary to obtain certificates for a domain under your control. A domain can be obtained
for free from [freenom](https://www.freenom.com/) or purchased from one of the many domain registrars.
DNS records must then be configured to point to the ip of the host where the app will be deployed.

With a domain in hand, i.e., `my-domain.com`, pointing to your host, the nginx config can be updated
and the first SSL certs obtained from letsencrypt with certbot by executing these commands providing
a valid email address

```
scripts/set-domains.sh my-domain.com www.my-domain.com
scripts/init-letsencrypt.sh -e my-email@gmail.com -d my-domain.com -d www.my-domain.com
```
---

With valid SSL certs and acme-challenge already stored in the `certbot-etc` and `certbot-www` docker
volumes created in the previous steps, the app can be run on port 80 with

```
docker-compose -f docker-compose.yml up
```

Go to [SSLLabs](https://www.ssllabs.com/ssltest/) to run their security test suite and you should
see that the site gets the top A+ grade. Now run a site audit with Lighthouse in Google
Chrome dev console and check out the all round top marks! Enjoy your minimal, secure, and performant app ;-)
