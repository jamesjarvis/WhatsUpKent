# WhatsUpKent

## Find out everything going on at the University of Kent, without nagging all of your mates.

---

This repo is essentially the _backend_ of WhatsUpKent, and contains the code for managing the WhatsUpKent knowledge graph.

There are 2 main applications:

- `scraper`: The application which continually scrapes timetable information from the Kent timetabling service
- `api`: The caching READ-ONLY database interface which is used for external applications to access the knowledge graph.

The 'knowledge graph' is run on a dgraph cluster, and if you would like to run all of this locally for development, simply run:

```bash
docker-compose up --build
```

If you would like to wipe the local database and start from scratch (if for example you are updating the schema), then run:

```bash
docker-compose down
docker volume rm whatsupkent_dgraph
```

## 🚀 Deployment

This is currently hosted on a _tiny_ VM running lightweight kubernetes (k3s). As such, the goal is to keep resource usage to a minimum, while remaining performant.
CI/CD is set up, so that any commit to master builds a new image, and deploys to the cluster. Currently there is not a staging service (due to resource contraints), so be careful to make sure that your commits actually work!
