# WhisperHouse

The main goal is monitoring who is in your appartment when you aren't at home.
WhisperHouse accomplishes this assuming each person constantly has a smartphone with himself, that makes monitoring macaddresses of those devices an easy way to chart presences.

It's possible to make others interesting things with it. WH is modular and relies on [probes](#probes) which record different informations into the database.

WhisperHouse consists in *different pieces that could reside in different places*.
 1. **NoSQL database** act to logs each input sent by probes, in this case i used [InfluxDB](https://github.com/influxdata/influxdb).
 2. **Probes**, act to read values from sensors and constantly send it to the database. For instance:
  * thermometer connected to the RaspberryPI's GPIO ([temperature](/temperature))
  * wifi-card which scans for macaddresses ([presence](/presence))
 3. A tool to realize **charts** with data who comes from NoSQL database, in this case i used [Grafana](https://github.com/grafana/grafana).

## Tips for a good deployment

You should implement TLS certificates to secure connections from probes to InfluxDB's API, and for Grafana's web interface.
Authentication between InfluxDB and probes acts through [Basic Access Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication) which is sent in clear (just Base64 encoded) over HTTP headers.
If you don't configure HTTPS connections, everybody could read credential for write into your database.

I used [Certbot](https://github.com/certbot/certbot) to get certificate and [automatic regenerate it each 3 months](https://wiki.archlinux.org/index.php/Let%E2%80%99s_Encrypt#Automatic_renewal) from [Let’s Encrypt](https://letsencrypt.org) CA, and [Nginx](https://nginx.org) as reverse proxy.

# Probes

Probes are those *stand-alone components* who read data from a sensor and sends it via https to database.

For example I developed two probes for purposes that i needed.
#### Context
I live with some flatmates and i would like to know who is at home when i'm not, and cross those data with the temperature of the apartament. This helps me understand who uses  pump up radiators (up to absurd temperatures).

I deployed those two probes into a RaspberryPI that we use as media center ([Kodi](https://kodi.tv)), so nobody get suspect for an always connected device with router.

That's just in case someone often controls connected devices on the router's admin panel. And because *the better place to hide something is putting it under the eyes of all*. Plus the little thermomether sensor fits perfectly into RaspberryPI's case.

A *genuine use* could be equity divide bills consider who was at home when the radiators were on. This could be automaticaly calculated by surveyed data.

---

The [probes](/probes) directory contains all the probes and relative **docs**.

To add a new probe you just have to write a program who reads data from a sensor, manages and sends it to database.
